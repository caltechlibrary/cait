#!/bin/bash

function getEADList {
    SAVE_PATH="$1"
    curl -o $SAVE_PATH/ead-index.html http://voro.cdlib.org/oac-ead/prime2002/caltech/
}

function processEADList {
    LIST_HTML="$1"
    SAVE_PATH="$2"
    echo "Processing list in $LIST_HTML"
    for ITEM in $(grep '<li><a href=' "$LIST_HTML" | cut -d\>  -f 3-1000 | cut -d\<  -f 1 | sed -e "s/^ //g"); do
        echo "Fetching $ITEM"
        curl -o $SAVE_PATH/$ITEM http://voro.cdlib.org/oac-ead/prime2002/caltech/$ITEM
    done
    exit 0
}

mkdir -p ead/oac-download
getEADList "ead/oac-download"
processEADList "ead/oac-download/ead-index.html" "ead/oac-download"
