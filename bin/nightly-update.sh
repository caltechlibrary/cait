#!/bin/bash
#

# This is an example cronjob to be run from the root account.
function consolelog {
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@"
}

# Change directory to where cait is installed
consolelog  "Running as $USER"
if [ "$USER" = "root" ]; then
    cd /archivesspace/cait
fi
consolelog "Working path $(pwd)"
# Load the cait configuration
if [ -f etc/setup.conf ]; then
    consolelog "Configuration $(pwd)/etc/setup.conf"
    . etc/setup.conf
fi

# Export the current content from ArchivesSpace
./bin/cait archivesspace export
# Generate webpages
./bin/genpages
# Index webpages
./bin/indexpages

# You should now be ready to reload the search engine/servepage service
if [ "$USER" = "root" ]; then
    /etc/init.d/servepages restart
else
    ./bin/servepages
fi
