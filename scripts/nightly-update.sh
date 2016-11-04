#!/bin/bash
#

# You can set the "HOME" directory for the deployment.
# e.g. /Sites/archives.example.edu
# export HOME=/Sites/archives.example.edu

# This is an example cronjob to be run from the root account.
function consolelog {
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@"
}

# Change directory to where cait is installed
consolelog  "Running as $USER"
cd $HOME

consolelog "Working path $(pwd)"
# Load the cait configuration
if [ -f $HOME/etc/cait.bash ]; then
    consolelog "Sourcing configuration $HOME/etc/cait.bash"
    . $HOME/etc/cait.bash
fi

# Export the current content from ArchivesSpace
$HOME/bin/cait archivesspace export
# Generate webpages
$HOME/bin/genpages
if [ "$USER" = "root" ]; then
    /etc/init.d/servepages stop
fi
# Index webpages
$HOME/bin/indexpages

# Generate sitemap
$HOME/bin/sitemapper htdocs htdocs/sitemap.xml $CAIT_SITE_URL

# You should now be ready to reload the search engine/servepage service
if [ "$USER" = "root" ]; then
    /etc/init.d/servepages start
else
    $HOME/bin/servepages
fi
