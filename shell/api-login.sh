#!/bin/bash

# Sanity check
if [ "$ASPACE_PROTOCOL" = "" ] || [ "$ASPACE_HOST" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

ASPACE_URL="$ASPACE_PROTOCOL://$ASPACE_HOST:$ASPACE_PORT"
echo "Accessing ArchivesSpace via $ASPACE_URL"

# Login
export TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_URL/users/$ASPACE_USERNAME/login | jq -r '.session')
echo 'export TOKEN='$TOKEN
echo 'curl -H "X-ArchivesSpace-Session: $TOKEN" '$ASPACE_URL
echo ""
