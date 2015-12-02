#!/bin/bash

#!/bin/bash

# Sanity check
if [ "$ASPACE_PROTOCOL" = "" ] || [ "$ASPACE_HOST" = "" ] || [ "$ASPACE_USERNAME" = "" ]; then
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
./aspace repository create '{"repo_code": "'$REPO_CODE'", "name": "'$REPO_NAME'"}'
echo "Done."
