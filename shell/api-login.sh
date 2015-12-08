#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

echo "Accessing ArchivesSpace via $ASPACE_API_URL"

# Login
export TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_API_URL/users/$ASPACE_USERNAME/login | jq -r '.session')
echo 'export TOKEN='$TOKEN
echo 'curl -H "X-ArchivesSpace-Session: $TOKEN" '$ASPACE_API_URL
echo ""
