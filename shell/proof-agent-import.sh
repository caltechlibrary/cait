#!/bin/bash
#

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
echo 'curl -H "X-ArchicesSpace-Session: $TOKEN" '$ASPACE_URL
echo ""

# Remove in stale /agents/people
KEY=$(date +%s)
echo "This is a destructive test. You are deleting data from $ASPACE_URL!!!!"
echo "Enter key to proceed: $KEY"
read INPUT_KEY
if [ "$KEY" != "$INPUT_KEY" ]; then
    echo "Aborting delete, aborting proof-agent-import.sh"
    exit 1
fi


function testAgents {
    aType=$1
    echo "Deleting /agents/$aType..."
    # ./aspace agent list '{"uri":"/agents/$aType"}' | sed -E "s/^\[|\]$//g;s/,/ /g"
    for I in $(./aspace agent list '{"uri":"/agents/'$aType'"}' | sed -E "s/^\[|\]$//g;s/,/ /g"); do
        if [ "$I" = "1" ]; then
            echo "Not deleting /agents/$aType/$I"
        else
            ./aspace agent delete '{"id":'$I',"uri":"/agents/'$aType'/'$I'"}'
        fi
    done

    # Try to import all the $aType in data/agents/$aType/*.json
    echo "Loading data from data/agents/$aType (this will take a while)..."
    find test/data/agents/$aType -type f | while read ITEM; do ./aspace -i $ITEM agent create; done > agent-$aType-import.log

    echo "Extracting errors from agent-$aType-import.log..."
    # Extract the errors from the agent-$aType-import.log
    cat agent-$aType-import.log | grep -E '^(Could not decode |, error: )' >> error-import.log

    if [ $(wc -l error-import.log | cut -d e -f 1 | sed -e "s/ //g") != "0" ]; then
        echo "We have errors....."
        cat error-import.log
        exit 1
    fi
}

# Build ./aspace binary to make sure we're testing the most current version
make build
START=$(date)
if [ -f error-import.log ];then
    echo "Clearing stale error-import.log"
    rm error-import.log
    touch error-import.log
fi
testAgents people
testAgents corporate_entities
testAgents families
testAgents software
echo "Done! Started: $START Completed: $(date)"
