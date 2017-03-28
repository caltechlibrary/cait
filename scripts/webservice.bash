#!/bin/bash

#
# Set the value of HOME to project directory
# export HOME=/Sites/archives.example.edu
# cd $HOME
#

# Run the web service with logging.
WEEKDAY="$(date +%A)"
export WEEKDAY
if [ -f etc/cait.bash ]; then
	. etc/cait.bash
fi
if [ ! -d logs ]; then
	mkdir -p logs
fi
bin/servepages >>"logs/servepages.${WEEKDAY}.log"
