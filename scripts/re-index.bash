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

# Index webpages
bleveIndexes=${CAIT_BLEVE/:/ }
for I in $bleveIndexes; do
    consolelog "Updating $I"
    pids=$(pgrep cait-servepages)
    if [ "$pids" != "" ]; then
        consolelog "Sending signal to swap out index $I"
        kill -s HUP $pids
    fi
    consolelog "Rebuilding index $I"
    bin/cait-indexpages >> logs/harvest.$WEEKDAY.log
done

