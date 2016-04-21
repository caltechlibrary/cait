#!/bin/bash
#

# This is an example cronjob to be run from the root account.

# Change directory to where cait is installed
echo "$(date +'%Y-%M-%D %H:%I:%S') Running as $USER"
if [ "$USER" = "root" ]; then
    cd /archivesspace/cait
fi
echo "$(date +'%Y-%M-%D %H:%I:%S') Working path $(pwd)"
# Load the cait configuration
if [ -f etc/setup.conf ]; then
    echo "$(date +'%Y-%M-%D %H:%I:%S') Configuration $(pwd)/etc/setup.conf"
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
    etc/init.d/servepages restart
fi
