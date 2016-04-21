#!/bin/bash
#

# This is an example cronjob to be run from the root account.

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

CWD=$(pwd)
mkdir -p testoutput
cd testoutput
ASPACE_API_TOKEN=$(curl -Fpassword=$CAIT_PASSWORD $CAIT_API_URL/users/$CAIT_USERNAME/login | jq -r '.session')
if [ "$ASPACE_API_TOKEN" = "" ]; then
    echo "login $CAIT_API_URL/users/$CAIT_USERNAME/login failed"
    exit 1
fi
ACCESSION_IDS=$(curl -H "X-ArchivesSpace-Session: $ASPACE_API_TOKEN" $CAIT_API_URL/repositories/2/accessions?all_ids=true)
echo "Accession IDS: $ACCESSION_IDS"
let L=$(echo $ACCESSION_IDS | jq length)-1
if [ "$L" = "-1" ]; then
    echo "ERROR: No accessions found $ACCESSION_IDS"
    exit 1
fi
## FIXME: Using $L to compute a range and pull a random element with my range command.
echo "Not finished adding tests yet!"
cd $CWD
exit 1
