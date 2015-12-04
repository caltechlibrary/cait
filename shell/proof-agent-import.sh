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
echo "This is very dangerous. You are deleting data!!!!"
echo "Enter key to proceed: $KEY"
read INPUT_KEY
if [ "$KEY" != "$INPUT_KEY" ]; then
    echo "Aborting delete of /agents/people"
    exit 1
fi
for I in $(./aspace agent list '{"uri":"/agents/people"}' | tr "," " " | sed -E "s/\[|]/ /g"); do 
    if [ "$I" != "1" ]; then 
        ./aspace agent delete '{"id":'$I',"uri":"/agents/people/'$I'"}'
    fi
done


# Try to import all the people in data/agents/people/*.json
make build

echo "Loading data from data/agents/people (this will take a while)..."
find data/agents/people -type f | while read ITEM; do ./aspace -i $ITEM agent create; done > agent-people-import.log

echo "Extracting errors from agent-people-import.log..."
# Extract the errors from the agent-people-import.log
cat agent-people-import.log | grep -E '^(Could not decode |, error: )' > error-import.log

if [ $(wc -l error-import.log | cut -d e -f 1 | sed -e "s/ //g") != "0" ]; then
    echo "We have errors....."
    cat error-import.log
    exit 1
fi
echo "Done!"
