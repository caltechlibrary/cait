#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

if [ "$1" = "" ]; then
    echo "Missing repository id (e.g. 2)."
    exit 1
fi

ASPACE_API_URL="$ASPACE_PROTOCOL://$ASPACE_HOST:$ASPACE_PORT"
echo "Accessing ArchivesSpace via $ASPACE_API_URL"

# Login
TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_API_URL/users/$ASPACE_USERNAME/login | jq -r '.session')


function getAccessions {
    REPO_ID=$1
    echo "Setting up data directory for repositories/$REPO_ID/accessions"
    mkdir -p data/repositories/$REPO_ID/accessions
    # Get a list of all agents ids
    echo "Getting ids for /repositories/$REPO_ID/accessions"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID/accessions?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data/$REPO_ID-accession-ids.txt

    # Now for each agent id in data/agents-*-ids.txt get a full record.
    echo "Reading /repositories/$REPO_ID/accessions ids and fetch their JSON records "
    cat data/$REPO_ID-accession-ids.txt | while read ACCESSION_ID; do
        if [ "$ACCESSION_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID/accessions/$ACCESSION_ID > data/repositories/$REPO_ID/accessions/$ACCESSION_ID.json
        fi
    done
    echo "Completed Accession dump for repository $REPO_ID"
}

getAccessions $1
echo ""
echo "Done."
