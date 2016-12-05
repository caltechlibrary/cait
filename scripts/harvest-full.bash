#!/bin/bash
#

# You can set the "HOME" directory for the deployment.
# e.g. /Sites/archives.example.edu
# export HOME=/Sites/archives.example.edu
# cd $HOME

# You can set the environment file for the deployment.
# e.g. /Sites/archives.example.edu/etc/cait.bash
# export CONFIG=/Sites/archives.example.edu

# This is an example cronjob to be run from the root account.
export WEEKDAY=$(date +%A)
if [ ! -d "logs" ]; then
    mkdir -p logs
fi

function consolelog {
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@"
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@" >> logs/harvest.$WEEKDAY.log
}

# Change directory to where cait is installed
consolelog  "Running as $USER"

if [ "$CONFIG" = "" ]; then
    export CONFIG=etc/cait.bash
fi

consolelog "Working path: $(pwd)"
# Load the cait configuration
if [ -f $CONFIG ]; then
    consolelog "Sourcing configuration $CONFIG"
    . $CONFIG
fi

# Export the current content from ArchivesSpace
bin/cait archivesspace export >> logs/harvest.$WEEKDAY.log

# Generate webpages
bin/genpages >> logs/harvest.$WEEKDAY.log

# Generate sitemap
bin/sitemapper htdocs htdocs/sitemap.xml $CAIT_SITE_URL >> logs/harvest.$WEEKDAY.log

# Index webpages
bleveIndexes=${CAIT_BELVE/:/ }
for I in $bleveIndexes; do
    consolelog "Updating $I"
    pids=$(pgrep servepages)
    if [ "$pids" != "" ]; then
        consolelog "Sending signal to swap out index $I"
        kill -s HUP $pids
    fi
    consolelog "Rebuilding index $I"
    bin/indexpages >> logs/harvest.$WEEKDAY.log
done

#
# Systemd setup notes:
# + copy etc/systemd/system/servepages.service to /etc/systemd/system/
# + Update /etc/systemd/system/servepages.service to match your system deployment
# + copy etc/cait.env-example to etc/cait.env
# + Update etc/cait.env to match your system deployment
# + start the web service with "systemctl daemon-reload && systemctl start servcepages"
#

# For development and non-systemd configuraiton try: 
# "cd $HOME && . etc/cait.bash && bin/servepages"


