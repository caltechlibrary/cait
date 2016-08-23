#!/bin/bash
#

# This is an example cronjob to be run from the root account.
function consolelog {
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@"
}

# Change directory to where cait is installed
consolelog  "Running as $USER"
cd $HOME

consolelog "Working path $(pwd)"
# Load the cait configuration
if [ -f /etc/cait.bash ]; then
    consolelog "Configuration /etc/cait.bash"
    . /etc/cait.bash
fi

# Export the current content from ArchivesSpace
cait archivesspace export
# Generate webpages
genpages
if [ "$USER" = "root" ]; then
    /etc/init.d/servepages stop
fi
# Index webpages
indexpages

# You should now be ready to reload the search engine/servepage service
if [ "$USER" = "root" ]; then
    /etc/init.d/servepages start
else
    servepages
fi
