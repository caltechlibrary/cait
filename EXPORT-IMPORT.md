
# The Utilities

This document covers some examples of using the _cait_ command line utilities to export
content from a production ArchivesSpace deployment to a local development ArchivesSpace deployment.
The most recent version of this document can be found at https://github.com/caltechlibrary/cait.


## Exporting from a production deployment

The easiest way to export content from a production ArchivesSpace deployment is using the _cait_ utility.

1. Set you environment variables
2. Use the *archivesspace export* option to create a local dump

### Example Assumptions

+ CAIT_USERNAME admin
+ CAIT_PASSWORD admin
+ CAIT_API_URL (for your production system) http://archives.example.edu:8089
+ CAIT_DATASET dataset

The following environment variables not note used in the export process

+ CAIT_SITE_URL
+ CAIT_HTDOCS
+ CAIT_HTDOCS_INDEX
+ CAIT_TEMPLATES

I am also assuming you have installed the _cait_ utility in *./bin/cait*

```
    export CAIT_API_URL=http://archives.example.edu:8089
    export CAIT_USERNAME=admin
    export CAIT_PASSWORD=admin
    export CAIT_DATASET=dataset

    ./bin/cait archivesspace export
    unset CAIT_USERNAME
    unset CAIT_PASSWORD
    unset CAIT_API_URL
```

This will take a while but it will create a local dump of the content in a directory called *data*. Each file is a JSON blob.
Since you don't want to accidentally disturb your production system it is a good idea that you unset the environment variables
when the export is complete.

## Importing into a development deployment

In this example we're assuming your *data* directory is already populated, you are using the Bash shell,
and the _cait_ utilities are installed in *./bin/*.

The basic setups are

1. Bring up an empty ArchivesSpace (follow the instructions at http://archiesspace.org)
2. Create a repository (this usually gets created as Repo ID 2)
3. Create any custom controlled vocabularies you need (e.g. extent types)
4. Load the Agents (I am assuming you only are interested in the people in this example)
5. Load the Subjects
6. Load the Accessions
7. Load the Digital Objects

### Example assumptions

+ CAIT_API_URL http://localhost:8089
+ CAIT_USERNAME admin
+ CAIT_PASSWORD admin
+ CAIT_DATASET dataset

The following environment variables not note used in the import process

+ CAIT_SITE_URL
+ CAIT_HTDOCS
+ CAIT_HTDOCS_INDEX
+ CAIT_TEMPLATES

Here's the stops to populate your local development ArchivesSpace. In this example I am assuming you're importing
into repository id of 2.

If you have any non-default extent_extent_type create them before proceeding

1. login to AS as admin
2. click on System
3. click Manage Controlled Value Lists
4. Select (from List Name) Extent Extent Type (extent_extent_type)
5. Add you additional values.


```
    export CAIT_API_URL=http://localhost:8089
    export CAIT_USERNAME=admin
    export CAIT_PASSWORD=admin
    export CAIT_DATASET=dataset

    # If you have non-default extent extent types, create them before proceeding
    # e.g. Multimedia, ProRes Master file, DVD
    ./bin/cait repository create -input dataset/repositories/2.json
    find dataset/subjects -type f | while read ITEM; do ./bin/cait subject create -input $ITEM; done
    find dataset/agents/people -type f | while read ITEM; do ./bin/cait agent create -input $ITEM; done
    find dataset/agents/corporate_entities -type f | while read ITEM; do ./bin/cait agent create -input $ITEM; done
    find dataset/repositories/2/digital_objects -type f | while read ITEM; do ./bin/cait digital_object create -input $ITEM; done
    find dataset/repositories/2/accessions -type f | while read ITEM; do ./bin/cait accession create -input $ITEM; done
```



You can import content from one ArchivesSpace deployment to the next using a combination of the _cait_ utility and basic shell scripting.

## Building a local dev site

The basic steps I take after having setup ArchivesSpace for development and loaded it with data is as follows.
The instructions assume you're in your *cait* repository directory and that all the *cait* tools are compiled and
installed in *./bin*.

### Environment required

+ CAIT_DATASET
+ CAIT_HTDOCS
+ CAIT_HTDOCS_INDEX
+ CAIT_SITE_URL
+ CAIT_TEMPLATES

### The workflow to generate website

1. Make sure the *CAIT_* environment variables are set.
2. Build the website with `./bin/genpages`
3. Create/update the sitemap with `./bin/sitemapper $CAIT_HTDOCS $CAIT_HTDOCS/sitemap.xml $CAIT_SITE_URL`
4. Index the site (this takes a while on my machine) `./bin/indexpages`
5. Launch `./bin/servepages` and test with your web browser
