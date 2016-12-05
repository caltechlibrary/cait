#!/bin/bash

#
# Set the value of HOME to project directory
# export HOME=/Sites/archives.example.edu
#

# Run the web service with logging.
cd $HOME
export WEEKDAY=$(date +%A)
if [ -f etc/cait.bash ]; then
    . etc/cait.bash
fi
bin/servepages >> logs/servepages.$WEEKDAY.log

