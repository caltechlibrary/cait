
# The Utilities

This document covers some examples of using the _aspace_ command line utlities to export
content from a production ArchivesSpace deployment to a local development ArchivesSpace deployment.
The most recent version of this document can be found at https://github.com/rsdoiel/aspace.


## Exporting from a production deployment

The easist way to export content from a production ArchivesSpace deployment is using the _aspace_ utility.

1. Set you environment variables
2. Use the instance export option to create a local dump

### Example Assumptions

+ ASPACE_USERNAME admin
+ ASPACE_PASSWORD admin
+ ASPACE_API_URL (for your production system) http://archives.example.edu:8089
+ ASPACE_DATASET data

The following environment variables not note used in the export process

+ ASPACE_SEARCH_URL
+ ASPACE_HTDOCS
+ ASPACE_TEMPLATES
+ ASPACE_BLEVE_INDEX

I am also assuming you have installed the _aspace_ utility in *./bin/aspace*

```
    export ASPACE_API_URL=http://archives.example.edu:8089
    export ASPACE_USERNAME=admin
    export ASPACE_PASSWORD=admin
    export ASPACE_DATASETS=data

    ./bin/aspace instance export
    unset ASPACE_USERNAME
    unset ASPACE_PASSWORD
    unset ASPACE_API_URL
```

This will take a while but it will create a local dump of the content in a directory called *data*. Each file is a JSON blob.
Since you don't want to accidentally disturb your production system it is a good idea that you unset the environent variables
when the export is complete.

## Importing into a development deployment

In this example we're assuming your *data* directory is already populated, you are using the Bash shell,
and the _aspace_ utilities are installed in *./bin/*.

The basic setups are

1. Bring up an empty ArchivesSpace instance (follow the instructions at http://archiesspace.org)
2. Create a repository (this usually gets created as Repo ID 2)
3. Create any custom controlled vocabularies you need (e.g. extent types)
4. Load the Agents (I am assuming you only are interested in the people in this example)
5. Load the Subjects
6. Load the Accessions
7. Load the Digital Objects

### Example assumptions

+ ASPACE_API_URL http://localhost:8089
+ ASPACE_USERNAME admin
+ ASPACE_PASSWORD admin
+ ASPACE_DATASETS data

The following environment variables not note used in the import process

+ ASPACE_SEARCH_URL
+ ASPACE_HTDOCS
+ ASPACE_TEMPLATES
+ ASPACE_BLEVE_INDEX

Here's the stops to populate your local development ArchivesSpace. In this example I am assuming you're importing
into repository id of 2.


```
    export ASPACE_API_URL=http://localhost:8089
    export ASPACE_USERNAME=admin
    export ASPACE_PASSWORD=admin
    export ASPACE_DATASETS=data

    ./bin/aspace repository create -i data/repositories/2.json
    find data/agents/people -type f | while read ITEM; do ./bin/aspace agent create -i $ITEM; done
    find data/subjects -type f | while read ITEM; do ./bin/aspace subject create -i $ITEM; done
    find data/repositories/2/accessions -type f | while read ITEM; do ./bin/aspace accession create -i $ITEM; done
    find data/repositories/2/digital_objects -type f | while read ITEM; do ./bin/aspace digital_object create -i $ITEM; done
```



You can import content from one ArchivesSpace deployment to the next using a combination of the _aspace_ utility and basic shell scripting.
