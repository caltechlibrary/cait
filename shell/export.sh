#!/bin/bash

# Sanity check
if [ "$ASPACE_PROTOCOL" = "" ] || [ "$ASPACE_HOST" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

ASPACE_URL="$ASPACE_PROTOCOL://$ASPACE_HOST:$ASPACE_PORT"
echo "Accessing ArchivesSpace via $ASPACE_URL"

# Login
TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_URL/users/$ASPACE_USERNAME/login | jq -r '.session')


function getAccessions {
    REPO_ID=$1
    echo "Setting up data directory for repositories/$REPO_ID/accessions"
    mkdir -p data-export/repositories/$REPO_ID/accessions
    # Get a list of all agents ids
    echo "Getting ids for /repositories/$REPO_ID/accessions"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/repositories/$REPO_ID/accessions?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data-export/$REPO_ID-accession-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading /repositories/$REPO_ID/accessions ids and fetch their JSON records "
    cat data-export/$REPO_ID-accession-ids.txt | while read ACCESSION_ID; do
        if [ "$ACCESSION_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/repositories/$REPO_ID/accessions/$ACCESSION_ID > data-export/repositories/$REPO_ID/accessions/$ACCESSION_ID.json
        fi
    done
    echo "Completed Accession export for repository $REPO_ID"
}


function getAgents {
    AGENT_TYPE=$1
    echo "Setting up data directory for $AGENT_TYPE"
    mkdir -p data-export/agents/$AGENT_TYPE
    # Get a list of all agents ids
    echo "Getting ids for /agents/$AGENT_TYPE"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/agents/$AGENT_TYPE?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data-export/$AGENT_TYPE-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading /agent/$AGENT_TYPE ids and fetch their JSON records "
    cat data-export/$AGENT_TYPE-ids.txt | while read AGENT_ID; do
        if [ "$AGENT_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/agents/$AGENT_TYPE/$AGENT_ID > data-export/agents/$AGENT_TYPE/$AGENT_ID.json
        fi
    done
    echo "Completed agent/$AGENT_TYPE export"
}


function getRepository {
    echo "Setting up data directory for repositories"
    mkdir -p data-export/repositories
    # Get a list of all agents ids
    echo "Getting ids for /agents/$AGENT_TYPE"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/repositories |  jq -r ".[].uri" | cut -d / -f 3 > data-export/repository-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading repository-paths and fetch their JSON records "
    cat data-export/repository-ids.txt | while read REPO_ID; do
        if [ "$REPO_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_URL/repositories/$REPO_ID > data-export/repositories/$REPO_ID.json
        fi
    done
    echo "Completed repositories export"
}


STARTED=$(date)
mkdir -p data-export
touch data-export/export.log
echo -e "$(date)\tExporting data from $ASPACE_URL\tstarted" >> data-export/export.log

echo -e "$(date)\tExporting /agents/people\tstarted" >> data-export/export.log
getAgents people
echo -e "$(date)\tExporting /agents/people\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /agents/corporate_entities\tstarted" >> data-export/export.log
getAgents corporate_entities
echo -e "$(date)\tExporting /agents/corporate_entities\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /agents/families\tstarted" >> data-export/export.log
getAgents families
echo -e "$(date)\tExporting /agents/families\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /agents/software\tstarted" >> data-export/export.log
getAgents software
echo -e "$(date)\tExporting /agents/software\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /repositories\tstarted" >> data-export/export.log
getRepository
echo -e "$(date)\tExporting /repositories\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting accessions by repositories\tstarted" >> data-export/export.log
cat data-export/repository-ids.txt | while read ID; do
    echo -e "$(date)\tExporting accessions for repository $ID\tstarted" >> data-export/export.log
    getAccessions $ID
    echo -e "$(date)\tExporting accessions for repository $ID\tfinished" >> data-export/export.log
done
echo -e "$(date)\tExporting accessions by repositories\tfinished" >> data-export/export.log
echo -e "$(date)\tExporting data from $ASPACE_URL\tfinished" >> data-export/export.log
echo ""
echo "Done. Started: $STARTED, Completed: $(date)"
