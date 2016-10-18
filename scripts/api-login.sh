#!/bin/bash

# Sanity check
if [ "$CAIT_API_URL" = "" ]; then
    echo "Enter the URL to the AchivesSpace REST API (e.g. http://localhost:8089) "
    read -p "URL: " CAIT_API_URL
    export CAIT_API_URL
fi
if [ "$CAIT_API_TOKEN" = "" ] && [ "$CAIT_PASSWORD" = "" ]; then
    echo "Enter the ArchivesSpace username and password to authenticate and get token:"
    read -p "Username: " CAIT_USERNAME
    read -s -p "Password: " CAIT_PASSWORD
    export CAIT_USERNAME
    export CAIT_PASSWORD
fi

# Login
echo "Accessing ArchivesSpace via $CAIT_API_URL getting CAIT_API_TOKEN"
CAIT_API_TOKEN=$(curl -Fpassword=$CAIT_PASSWORD $CAIT_API_URL/users/$CAIT_USERNAME/login | jq -r '.session')
if [ "$CAIT_API_TOKEN" = "" ] | [ "$CAIT_API_TOKEN" = "null" ]; then
    echo "Login failed."
    echo 
    echo 'Try the collowing to view the error returned.'
    echo 
    echo '    curl -Fpassword=$CAIT_PASSWORD $CAIT_API_URL/users/$CAIT_USERNAME/login | jq -r ".error"'
    echo 
    exit 1
fi
echo 'Running the following export command--'
echo ''
echo 'export CAIT_API_TOKEN='$CAIT_API_TOKEN
echo ''
echo 'Example curl command usage'
echo '   curl -H "X-ArchivesSpace-Session: $CAIT_API_TOKEN" '$CAIT_API_URL
echo ""
export CAIT_API_TOKEN
