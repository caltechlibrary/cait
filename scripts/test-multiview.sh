#!/bin/bash
#

# This is test script for multiviews handling by servepages

# Change directory to where cait is installed
echo "$(date +'%Y-%M-%D %H:%I:%S') Running as $USER"
if [ "$USER" = "root" ]; then
    echo "ERROR: Should not run this test as root"
    exit 1
fi
echo "$(date +'%Y-%M-%D %H:%I:%S') Working path $(pwd)"
# Load the cait configuration
if [ -f etc/setup.conf ]; then
    echo "$(date +'%Y-%M-%D %H:%I:%S') Configuration $(pwd)/etc/setup.conf"
    . etc/setup.conf
fi

CAIT_API_TOKEN=$(curl -Fpassword=$CAIT_PASSWORD $CAIT_API_URL/users/$CAIT_USERNAME/login | jq -r '.session')
if [ "$CAIT_API_TOKEN" = "" ]; then
    echo "login $CAIT_API_URL/users/$CAIT_USERNAME/login failed"
    exit 1
fi
ACCESSION_ID=990
echo "Starting servepages"
./bin/servepages &
PID=$!
echo "Waiting for 5 seconds"
sleep 5
echo "Checking for $CAIT_SITE_URL/repositories/2/accession/$ACCESSION_ID"
DATA=$(curl -H "X-ArchivesSpace-Session: $CAIT_API_TOKEN" $CAIT_SITE_URL/repositories/2/accessions/$ACCESSION_ID)
echo "Data: $DATA"
echo "Stopping servepage"
kill $PID
