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
    ID=$(jsoncols -i $FNAME '.id')
    URL="$CAIT_ARCHIVESSPACE_URL/agents/agent_person/$ID"
    for I in $(jsonrange -i $FNAME .names); do
        PRIMARY_NAME=$(jsoncols -i $FNAME ".names[$I].primary_name")
        REST_OF_NAME=$(jsoncols -i $FNAME ".names[$I].rest_of_name")
        SORT_NAME=$(jsoncols -i $FNAME ".names[$I].sort_name")
        IS_DISPLAY_NAME=$(jsoncols -i $FNAME ".names[$I].is_display_name")
        # Output delimited record
        csvcols -d "|" "$ID|$PRIMARY_NAME|$REST_OF_NAME|$SORT_NAME|$IS_DISPLAY_NAME|$URL"
    done
}


#
# Main code
#
CheckEnv CAIT_DATASET CAIT_ARCHIVESSPACE_URL
CheckSoftware cut grep findfile csvcols jsoncols jsonrange
csvcols -d "|" "Agent/People ID|Primary Name|Rest of Name|Sort Name|Is Display Name|URL"
findfile -s .json $CAIT_DATASET/agents/people | while read ITEM; do
    RECORD=$(GetRecord $CAIT_DATASET/agents/people/$ITEM)
    echo "$RECORD"
done

