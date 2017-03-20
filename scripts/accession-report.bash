#!/bin/bash

# 
# Here is my first cut of the data needed for a basic AS data spreadsheet:
# 
# Accession Identifier
# Accession Title - LINK TO ACCESSION
# Accession Date
# Accession Publish (T/F)
# 
# Date Creation
# 
# Extent Type
# 
# Agent Role
# Agent Name - LINK TO AGENT
# 
# Subject - LINK TO SUBJECT
# 
# Let's start with this and work out from there.
# Links optional.
# Thanks!
# 
# Stephen
# 

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
    ID=$(jsoncols -i $FNAME .id )
    TITLE=$(jsoncols -i $FNAME .title)
    ID_0=$(jsoncols -i $FNAME  .id_0 | sed -E 's/"//g')
    ID_1=$(jsoncols -i $FNAME .id_1 | sed -E 's/"//g')
    ACCESSION_DATE=$(jsoncols -i $FNAME .accession_date | sed -E 's/"//g')
    PUBLISH=$(jsoncols -i $FNAME .publish | sed -E 's/"//g')
    IDENTIFIER="$ID_0 $ID_1"
    URL="$CAIT_ARCHIVESSPACE_URL/accessions/$ID"
    EXTENT_COUNT=$(jsonrange -i $FNAME -length .extents)
    ## FIXME: generate a semi-colon delimited list of people associated this Accession
    if [ "$EXTENT_COUNT" = "0" ]; then
        EXTENT_TYPE=""
        PHYSICAL_DETAILS=""
        csvcols -d "|" "$ID|$TITLE|$IDENTIFIER|$EXTENT_TYPE|$PHYSICAL_DETAILS|$ACCESSION_DATE|$PUBLISH|$URL"
     else
         for I in $(jsonrange -i $FNAME .extents); do
             # Fetch the extent for 
             EXTENT_TYPE=$(jsoncols -i $FNAME .extents[$I].extent_type)
             PHYSICAL_DETAILS=$(jsoncols -i $FNAME .extents[$I].physical_details)
             if [ "$EXTENT_TYPE" != "null" ]; then
                 # Output delimited record
                 csvcols -d "|" "$ID|$TITLE|$IDENTIFIER|$EXTENT_TYPE|$PHYSICAL_DETAILS|$ACCESSION_DATE|$PUBLISH|$URL"
             fi
         done
    fi 
}


#
# Main code
#
CheckEnv CAIT_DATASET CAIT_ARCHIVESSPACE_URL
CheckSoftware cut grep findfile csvcols jsoncols range
csvcols -d "|" "ArchivesSpace ID | Title | Identifier | Extent Type | Physical Details | Accession Date | Publish | URL"
findfile -f -s .json $CAIT_DATASET/repositories/2/accessions | while read ITEM; do
    RECORD=$(GetRecord $ITEM)
    if [ "$RECORD" != "" ]; then
        echo "$RECORD"
    fi
done

