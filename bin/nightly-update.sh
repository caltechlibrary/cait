#!/bin/bash
#

# This is an example cronjob to be run from the root account.


# Load the cait configuration
. /archivesspace/cait/etc/setup.conf

# Change directory to where cait is installed
cd /archivesspace/cait
# Export the current content from ArchivesSpace
./bin/cait archivesspace export
# Generate webpages
./bin/genpages
# Index webpages
./bin/indexpages

# You should now be ready to reload the search engine
/etc/init.d/servepages stop
/etc/init.d/servepages start

