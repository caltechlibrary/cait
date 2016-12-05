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

# Generate sitemap
$HOME/bin/sitemapper htdocs htdocs/sitemap.xml $CAIT_SITE_URL

# Index webpages
bleveIndexes=${CAIT_BELVE/:/ }
for I in $bleveIndexes; do
    console "Updating $I"
    pids=$(pgrep servepages)
    if [ "$pids" != "" ]; then
        console "Sending signal to swap out index $I"
        kill -s HUP $pids
    fi
    echo "Rebuilding index $I"
    $HOME/bin/indexpages
done

#if [ "$USER" = "root" ]; then
# Sys Init restart 
#    /etc/init.d/servepages restart
# Systemd reload environment and restart 
#    systemctl daemon-reload
#    systemctl restart servepages
#fi

