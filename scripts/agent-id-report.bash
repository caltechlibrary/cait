#!/bin/bash

#
# Report all the Agent/Person ids found, their full name(s), if the record looks complete
# and link back to Archives's Person objects for details.
#

# check for software
function CheckSoftware () {
    for ITEM in $@; do
        PROG=$(which "$ITEM")
        if [ "$PROG" = "" ]; then
            echo "Missing $1";
            exit 1
        fi
    done
    #echo "Found required: $@" 
}

# check for VAR
function CheckEnv () {
    VAR=$(env | grep $1)
    if [ "$VAR" = "" ]; then
        echo "Missing environment varaible: \$$1"
        exit 1
    fi
}

# GetRecord from JSON blob
function GetRecord () {
    FNAME="$1"
    ID=$(jq '.id' $FNAME)
    PRIMARY_NAME=$(jq '.names[0].primary_name' $FNAME)
    REST_OF_NAME=$(jq '.names[0].rest_of_name' $FNAME)
    SORT_NAME=$(jq '.names[0].sort_name' $FNAME)
    IS_DISPLAY_NAME=$(jq '.names[0].is_display_name' $FNAME)

    # Output delimited record
    csvcols -d "|" "agent:person:$ID|$PRIMARY_NAME|$REST_OF_NAME|$SORT_NAME|$IS_DISPLAY_NAME"
}


#
# Main code
#
CheckEnv CAIT_DATASET
CheckSoftware cut grep findfile csvcols jq
csvcols -d "|" "ArchivesSpace ID|Primary Name|Rest of Name| Sort Name| Is Display Name"
findfile -s .json $CAIT_DATASET/agents/people | while read ITEM; do
    RECORD=$(GetRecord $CAIT_DATASET/agents/people/$ITEM)
    echo "$RECORD"
done

