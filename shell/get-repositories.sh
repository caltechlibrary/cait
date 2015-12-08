#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_HOST" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

ASPACE_API_URL="$ASPACE_PROTOCOL://$ASPACE_HOST:$ASPACE_PORT"
echo "Accessing ArchivesSpace via $ASPACE_API_URL"

# Login
TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_API_URL/users/$ASPACE_USERNAME/login | jq '.session' | cut -d \" -f 1)


function getRepository {
    echo "Setting up data directory for repositories"
    mkdir -p data/repositories
    # Get a list of all agents ids
    echo "Getting ids for /agents/$AGENT_TYPE"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories |  jq -r ".[].uri" | cut -d / -f 3 > data/repository-ids.txt

    # Now for each agent id in data/agents-*-ids.txt get a full record.
    echo "Reading repository-paths and fetch their JSON records "
    cat data/repository-ids.txt | while read REPO_ID; do
        if [ "$REPO_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID > data/repositories/$REPO_ID.json
        fi
    done
    echo "Completed repositories list"
}

getRepository
echo ""
echo "Done."
