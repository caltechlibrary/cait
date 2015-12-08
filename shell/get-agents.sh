#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

ASPACE_API_URL="$ASPACE_PROTOCOL://$ASPACE_HOST:$ASPACE_PORT"
echo "Accessing ArchivesSpace via $ASPACE_API_URL"

# Login
TOKEN=$(curl -Fpassword=$ASPACE_PASSWORD $ASPACE_API_URL/users/$ASPACE_USERNAME/login | jq '.session' | cut -d \" -f 1)


function getAgents {
    AGENT_TYPE=$1
    echo "Setting up data directory for $AGENT_TYPE"
    mkdir -p data/agents/$AGENT_TYPE
    # Get a list of all agents ids
    echo "Getting ids for /agents/$AGENT_TYPE"
    curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/agents/$AGENT_TYPE?all_ids=true | sed -E "s/\[//;s/,/ /g;s/]//" | tr " " "\n" > data/$AGENT_TYPE-ids.txt

    # Now for each agent id in data/agents-*-ids.txt get a full record.
    echo "Reading /agent/$AGENT_TYPE ids and fetch their JSON records "
    cat data/$AGENT_TYPE-ids.txt | while read AGENT_ID; do
        if [ "$AGENT_ID" != "" ]; then
            curl -H "X-ArchivesSpace-Session: $TOKEN" $ASPACE_API_URL/agents/$AGENT_TYPE/$AGENT_ID > data/agents/$AGENT_TYPE/$AGENT_ID.json
        fi
    done
    echo "Completed $AGENT_TYPE"
}

getAgents people
getAgents corporate_entities
getAgents families
getAgents software
echo ""
echo "Done."
