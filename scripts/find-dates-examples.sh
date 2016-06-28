#!/bin/bash

JQ=$(which jq)
if [ "$JQ" = "" ]; then
    echo "This script requires jq. See https://stedolan.github.io/jq for how to obtain it"
    exit 1
fi


if [ "$1" = "" ];then
    echo "USAGE: bash scripts/find-dates.sh DATE_TYPE"
    echo "Date types: single, inclusive, bulk or all"
    exit 1
fi

find dataset/repositories/2/accessions -type f |\
while read ITEM; do
    if [ "$1" = "all" ]; then
        jq "{\"accession_id\": .id, \"dates\":.dates}" "$ITEM"
    else
        jq "{\"accession_id\": .id, \"dates\":.dates}| select(.dates[].date_type == \"$1\")" "$ITEM"
    fi
done
