#!/bin/bash
#

# You can set the "HOME" directory for the deployment.
# e.g. /Sites/archives.example.edu
# export HOME=/Sites/archives.example.edu

# You can set the environment file for the deployment.
# e.g. /Sites/archives.example.edu/etc/cait.bash
# export CONFIG=/Sites/archives.example.edu

# This is an example cronjob to be run from the root account.
function consolelog {
    echo $(date +"%Y/%m/%d %H:%M:%S")" $@"
}

# Change directory to where cait is installed
consolelog  "Running as $USER"
cd $HOME

if [ "$CONFIG" = "" ]; then
    export CONFIG=$HOME/etc/cait.bash
fi

consolelog "Working path $(pwd)"
# Load the cait configuration
if [ -f $CONFIG ]; then
    consolelog "Sourcing configuration $CONFIG"
    . $CONFIG
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


