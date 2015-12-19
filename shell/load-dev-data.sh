#!/bin/bash
export ASPACE_API_URL=http://localhost:8089
export ASPACE_API_TOKEN=
export ASPACE_USERNAME=admin
export ASPACE_PASSWORD=admin
echo "This script is intended to work with an empty developer instance of ArchiveSpace."

# Load the repository
./aspace repository create -i data-import/repositories/2.json
# Load Subjects
find data-import/subjects -type f | while read ITEM; do
    ./aspace subject create -i "$ITEM"
done
# Load Vocabulary & Terms
find data-import/vocabularies -type f -depth 1 | while read ITEM; do
    ./aspace vocabulary create -i "$ITEM"
    TERM_FILE=$(echo $ITEM | cut -d \. -f 1)/terms.json
    ./aspace term create -i "$TERM_FILE"
done
# Load Agents
find data-import/agents/people -type f | while read ITEM; do
    if [ "$ITEM" != "data-import/agents/people/1.json" ]; then
        ./aspace agent create -i "$ITEM"
    fi
done
find data-import/agents/corporate_entities -type f | while read ITEM; do
    if [ "$ITEM" != "data-import/agents/corporate_entities/1.json" ]; then
        ./aspace agent create -i "$ITEM"
    fi
done
find data-import/agents/families -type f | while read ITEM; do
    ./aspace agent create -i "$ITEM"
done
# NOTE: agents/software is defined by installing ArchivesSpace the first time.

# FIXME: Before Accessions can be successfully loaded we need to define the extent types
# Load Accessions
find data-import/repositories/2/accessions -type f | while read ITEM; do
    ./aspace accession create -i "$ITEM"
done
