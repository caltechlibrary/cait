#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

echo "Accessing ArchivesSpace via $ASPACE_API_URL"

# Login
TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_API_URL/users/$ASPACE_USERNAME/login | jq -r '.session')


function getAccessions {
    REPO_ID=$1
    echo "Setting up data directory for repositories/$REPO_ID/accessions"
    mkdir -p data-export/repositories/$REPO_ID/accessions
    # Get a list of all agents ids
    echo "Getting ids for /repositories/$REPO_ID/accessions"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID/accessions?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data-export/$REPO_ID-accession-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading /repositories/$REPO_ID/accessions ids and fetch their JSON records "
    cat data-export/$REPO_ID-accession-ids.txt | while read ACCESSION_ID; do
        if [ "$ACCESSION_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID/accessions/$ACCESSION_ID > data-export/repositories/$REPO_ID/accessions/$ACCESSION_ID.json
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
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/agents/$AGENT_TYPE?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data-export/$AGENT_TYPE-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading /agent/$AGENT_TYPE ids and fetch their JSON records "
    cat data-export/$AGENT_TYPE-ids.txt | while read AGENT_ID; do
        if [ "$AGENT_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/agents/$AGENT_TYPE/$AGENT_ID > data-export/agents/$AGENT_TYPE/$AGENT_ID.json
        fi
    done
    echo "Completed agent/$AGENT_TYPE export"
}


function getRepository {
    echo "Setting up data directory for repositories"
    mkdir -p data-export/repositories
    # Get a list of all agents ids
    echo "Getting ids for /agents/$AGENT_TYPE"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories |  jq -r ".[].uri" | cut -d / -f 3 > data-export/repository-ids.txt

    # Now for each agent id in data-export/agents-*-ids.txt get a full record.
    echo "Reading repository-paths and fetch their JSON records "
    cat data-export/repository-ids.txt | while read REPO_ID; do
        if [ "$REPO_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/repositories/$REPO_ID > data-export/repositories/$REPO_ID.json
        fi
    done
    echo "Completed repositories export"
}

function getSubjects {
    echo "Setting up data directory for subjects"
    mkdir -p data-export/subjects
    # Get a list of all subject ids
    echo "Getting ids for /subjects"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/subjects?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" >  data-export/subject-ids.txt

    # Now for each id in data-export/subject-ids.txt get a full record.
    echo "Reading subjects and fetch their JSON records "
    cat data-export/subject-ids.txt | while read SUBJECT_ID; do
        if [ "$SUBJECT_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/subjects/$SUBJECT_ID > data-export/subjects/$SUBJECT_ID.json
        fi
    done
    echo "Completed subjects export"

}

function getLocations {
    echo "Setting up data directory for locations"
    mkdir -p data-export/locations
    # Get a list of all location ids
    echo "Getting ids for /locations"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/locations?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" >  data-export/location-ids.txt

    # Now for each id in data-export/location-ids.txt get a full record.
    echo "Reading locations and fetch their JSON records "
    cat data-export/location-ids.txt | while read ID; do
        if [ "$ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/locations/$ID > data-export/locations/$ID.json
        fi
    done
    echo "Completed locations export"

}

function getVocabularies {
    echo "Setting up data directory for vocabularies"
    mkdir -p data-export/vocabularies
    # Get a list of all ids
    echo "Getting ids for /vocabularies"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/vocabularies?all_ids=true | jq -r '.[].uri' | cut -d / -f 3 >  data-export/vocabulary-ids.txt

    # Now for each id in data-export/vocabulary-ids.txt get a full record.
    echo "Read and fetch the JSON records "
    cat data-export/vocabulary-ids.txt | while read ID; do
        if [ "$ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/vocabularies/$ID > data-export/vocabularies/$ID.json
        fi
    done
    echo "Completed locations export"

}

function getTerms {
    echo "Getting terms for each vocabulary"
    cat data-export/vocabulary-ids.txt | while read vocID; do
        echo "Setting up data directory for /vocabularies/$vocID"
        mkdir -p data-export/vocabularies/$vocID
        # Get a list of all ids
        echo "Getting /vocabularies/$vocID/terms.json"
        curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/vocabularies/$vocID/terms >  data-export/vocabularies/$vocID/terms.json
    done
    echo "Completed terms export"
}

STARTED=$(date)
mkdir -p data-export
touch data-export/export.log
echo -e "$(date)\tExporting data from $ASPACE_API_URL\tstarted" >> data-export/export.log

echo -e "$(date)\tExporting /vocabularies\tstarted" >> data-export/export.log
getVocabularies
echo -e "$(date)\tExporting /vocabularies\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /vocabularies/.../terms\tstarted" >> data-export/export.log
getTerms
echo -e "$(date)\tExporting /vocabularies/.../terms\tfinished" >> data-export/export.log

echo -e "$(date)\tExporting /locations\tstarted" >> data-export/export.log
getLocations
echo -e "$(date)\tExporting /locations\tfinished" >> data-export/export.log

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

echo -e "$(date)\tExporting /subjects\tstarted" >> data-export/export.log
getSubjects
echo -e "$(date)\tExporting /subjects\tfinished" >> data-export/export.log

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
echo -e "$(date)\tExporting data from $ASPACE_API_URL\tfinished" >> data-export/export.log
echo ""
echo "Done. Started: $STARTED, Completed: $(date)"
