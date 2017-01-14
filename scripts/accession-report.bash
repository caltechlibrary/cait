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
    ID=$(jsoncols .id $FNAME)
    TITLE=$(jsoncols .title $FNAME)
    ID_0=$(jsoncols .id_0 $FNAME | sed -E 's/"//g')
    ID_1=$(jsoncols .id_1 $FNAME | sed -E 's/"//g')
    IDENTIFIER="$ID_0 $ID_1"
    EXTENT_TYPE=$(jsoncols .extents[0].extent_type $FNAME)
    PHYSICAL_DETAILS=$(jsoncols .extents[0].physical_details $FNAME)
    URL="https://caltecharchives.lyrasistechnology.org/repositories/2/accession/$ID"

    # Output delimited record
    csvcols -d "|" "$URL|repositories:2:accession:$ID|$TITLE|$IDENTIFIER|$EXTENT_TYPE|$PHYSICAL_DETAILS"
}


#
# Main code
#
CheckEnv CAIT_DATASET
CheckSoftware cut grep findfile csvcols jsoncols
csvcols -d "|" "url | ArchivesSpace ID | Title | Identifier | Extent Type| Physical Details"
findfile -s .json $CAIT_DATASET/repositories/2/accessions | while read ITEM; do
    RECORD=$(GetRecord $CAIT_DATASET/repositories/2/accessions/$ITEM)
    echo "$RECORD"
done

