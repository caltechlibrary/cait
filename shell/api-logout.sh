#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" != "" ] && [ "$ASPACE_API_TOKEN" != "" ]; then
    echo "Sending logging out with $ASPACE_API_URL/logout"
    curl -H "X-ArchivesSpace-Session: $ASPACE_API_TOKEN" $ASPACE_API_URL/logout > /dev/null
fi
export ASPACE_API_URL=""
export ASPACE_API_TOKEN=""
export ASPACE_USERNAME=""
export ASPACE_PASSWORD=""
echo 'ASPACE_API_URL, ASPACE_API_TOKEN, ASPACE_USERNAME and ASPACE_PASSWORD set to empty strings'
