#!/bin/bash

#!/bin/bash

# Sanity check
if [ "$ASPACE_API_URL" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
    echo "You need to setup your environment variables for accessing your ArchivesSpace deployment"
    exit 1
fi

echo "Beginning test of aspace tool"
if [ ! -f ./aspace ]; then
    make test
    make build
fi

REPO_CODE="TEST "$(date "+%Y-%m-%d %H:%M:%S")
REPO_NAME="This is a test of aspace tool"
echo "Creating a repository called $REPO_CODE"
echo ./aspace repository create '{"repo_code": "'$REPO_CODE'", "name": "'$REPO_NAME'"}'
RESPONSE=$(./aspace repository create '{"repo_code": "'$REPO_CODE'", "name": "'$REPO_NAME'"}')
REPO_ID=$(echo $RESPONSE | jq ".id")
if [ "$REPO_ID" = "" ]; then
    echo $RESPONSE
    exit 1
fi
echo ./aspace repository list '{"id": '$REPO_ID'}'
RESPONSE=$(./aspace repository list '{"id": '$REPO_ID'}')
T=$(echo $RESPONSE | jq -r ".repo_code")
if [ "$REPO_CODE" != "$T" ]; then
    echo "Can't find .repo_code: $REPO_CODE != $T"
    echo $RESPONSE
    exit 1
fi
PAYLOAD=$(echo $RESPONSE | sed -e "s/TEST/testme/")
echo ./aspace repository update $PAYLOAD
RESPONSE=$(./aspace repository update $PAYLOAD)
T=$(echo $RESPONSE | jq -r ".status")
if [ "$T" = "" ];then
    echo $RESPONSE
    exit 1
fi

find test/data/agents/people -type f | while read ITEM; do
    CHECK_NAME=$(echo $ITEM | grep -vE "/1.json$")
    if [ "$CHECK_NAME" != "" ];then
        echo "Importing $ITEM"
        RESPONSE=$(./aspace agent create -i $ITEM)
        ERROR=$(echo $RESPONSE | jq -r ".error")
        STATUS=$(echo $RESPONSE | jq -r ".status")
        if [ "$ERROR" != "null" ] ||  [ "$STATUS" != "Created" ]; then
            echo $RESPONSE
            exit 1
        fi
    fi
done

EXIT_CODE=$?
if [ "$EXIT_CODE" != "0" ]; then
    echo "Error importing agents/people"
    exit 1
fi

find test/data/agents/corporate_entities -type f | while read ITEM; do
    CHECK_NAME=$(echo $ITEM | grep -vE "/1.json$")
    if [ "$CHECK_NAME" != "" ];then
        echo "Importing $ITEM"
        RESPONSE=$(./aspace agent create -i $ITEM)
        ERROR=$(echo $RESPONSE | jq -r ".error")
        STATUS=$(echo $RESPONSE | jq -r ".status")
        if [ "$ERROR" != "null" ] || [ "$STATUS" != "Created" ]; then
            echo $RESPONSE
            exit 1
        fi
    fi
done

EXIT_CODE=$?
if [ "$EXIT_CODE" != "0" ]; then
    echo "Error importing agents/corporate_entities"
    exit 1
fi

find test/data/agents/families -type f | while read ITEM; do
    CHECK_NAME=$(echo $ITEM | grep -vE "/1.json$")
    if [ "$CHECK_NAME" != "" ];then
        echo "Importing $ITEM"
        RESPONSE=$(./aspace agent create -i $ITEM)
        ERROR=$(echo $RESPONSE | jq -r ".error")
        STATUS=$(echo $RESPONSE | jq -r ".status")
        if [ "$ERROR" != "null" ] || [ "$STATUS" != "Created" ]; then
            echo $RESPONSE
            exit 1
        fi
    fi
done

EXIT_CODE=$?
if [ "$EXIT_CODE" != "0" ]; then
    echo "Error importing agents/families"
    exit 1
fi

#find test/data/repositories/2/accessions -type f | while read ITEM; do
#    echo "Importing $ITEM"
#done

echo ./aspace repository delete '{"id": '$REPO_ID'}'
REPONSE=$(./aspace repository delete '{"id": '$REPO_ID'}')
T=$(echo $RESPONSE | jq -r ".status")
if [ "$T" = "" ];then
    echo $RESPONSE
    exit 1
fi

echo -e "PASS\nok $(date)"
