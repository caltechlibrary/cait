#!/bin/bash

# Sanity check
if [ "$CAIT_API_URL" != "" ] && [ "$CAIT_API_TOKEN" != "" ]; then
    echo "Sending logging out with $CAIT_API_URL/logout"
    curl -H "X-ArchivesSpace-Session: $CAIT_API_TOKEN" $CAIT_API_URL/logout > /dev/null
fi
export CAIT_API_URL=""
export CAIT_API_TOKEN=""
export CAIT_USERNAME=""
export CAIT_PASSWORD=""
echo 'CAIT_API_URL, CAIT_API_TOKEN, CAIT_USERNAME and CAIT_PASSWORD set to empty strings'
