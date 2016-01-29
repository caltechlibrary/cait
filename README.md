
# aspace

_aspace_ is a set of utilities written in the [Go](http://golang.org) language that
work with and augment the [ArchivesSpace](http://archivesspace.org) API.

+ aspace - a command line utility for ArchivesSpace interaction (basic CRUD operations and export)
+ aspacepage - a simple static page generator based on export ArchivesSpace content
+ aspaceindexer - for indexing exported JSON structures with [Bleve](https://github.com/blevesearch/bleve)
+ aspacesearch - a web service providing public search services and content browsing
+ xlsximporter - a tool for turning Excel spreadsheets in .xlsx format into JSON files suitable for importing into ArchivesSpace

## Requirements

+ A working deployment of ArchivesSpace
+ Golang 1.5.3 or better to compile
+ [Bleve](http://blevesearch.com) Golang based Search/Indexing library (think Lucene lite implemented in Golang)
+ Two 3rd part Go packages
    + [Otto](https://github.com/robertkrimen/otto) by Robert Krimen, MIT license
    + [xlsx](https://github.com/tealeg/xlsx) by Tealeg, BSD license

## Compiling

If you already have [Go](https://golang.org) setup and installed compiling the utilties are pretty
straight forward.

1. Clone the git repository for the project
2. "Go get" the 3rd party libraries
3. Compile
4. Setup the necessary environment variables for using the utilities

Here's a typical example of setting things up.

```
    git clone git@github.com:rsdoiel/aspace.git
    cd aspace
    go get github.com/robertkrimen/otto
    go get github.com/tealeg/xlsx
    mkdir bin
    go build -o bin/aspace cmds/aspace/aspace.go
    go build -o bin/aspacepage  cmds/aspacepage/aspacepage.go
    go build -o bin/aspaceindexer cmds/aspaceindexer/aspaceindexer.go
    go build -o bin/aspacesearch cmds/aspacesearch/aspacesearch.go
    go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go
```

At this point you should have your command line utilities ready to go in the *bin* directory. You
are now ready to setup your environment variables.


## Setting up your environment

The command line tools and services are configured via environment variables. Below is an example
of setting things up under Bash running on your favorite Unix-like system.


```
    #
    # setup.sh - this script sets the environment variables for aspace project.
    # You would source file before using aspace, aspaceindexer, aspacesearch
    # or aspacedashbaord.
    #

    #
    # Local Development setup
    #
    export ASPACE_API_URL=http://localhost:8089
    export ASPACE_USERNAME=admin
    export ASPACE_PASSWORD=admin
    export ASPACE_DATASETS=data
    export ASPACE_SEARCH_URL=http://localhost:8501
    export ASPACE_HTDOCS=htdocs
    export ASPACE_TEMPLATES=templates/default
    export ASPACE_BLEVE_INDEX=index.bleve

    #
    # Create the necessary directory structure
    #
    mkdir -p $ASPACE_DATASETS
    mkdir -p $ASPACE_HTDOCS
    mkdir -p $ASPACE_TEMPLATES

```

Assuming Bash and that you've named the file _setup.sh_ you could
source the file from your shell prompt by typing `. setup.sh`.

## Utilities

### _aspace_

This command is a general purpose tool for fetch ArchivesSpace data from the
ArchivesSpace REST API, saving or modifying that data as well as querying the
locally capture output of the API.

Current _aspace_ supports operations on repositories, subjects, agents, accessions and digital_objects.

These are the common actions that can be performed

+ create
+ list (individually or all ids)
+ update (can use a file instead of the command line, see -i option)
+ delete
+ export (useful with integrating into static websites or batch processing via scripts)

Here's an example session of using the _aspace_ command line tool on the repository object.

```shell
    . setup.sh # Source my setup file so I can get access to the API
    aspace repository create '{"uri":"/repositories/3","repo_code":"My Archive","name":"My Archive"}' # Create an archive called My Archive
    aspace repository list # show a list of archives, for example purposes we'll use archive ID of 3
    aspace repository list '{"uri":"/repositories/3"}' # Show only the archive JSON for repository ID equal to 3
    aspace repository list '{"uri":"/repositories/3"}' > repo2.json # Save the output to the file repo3.json
    aspace repository update -i repo3.json # Save your changes back to ArchivesSpace
    aspace repository export '{"uri":"/repositories/3"}' # export the repository metadata to data/repositories/3.json
    aspace repository delete '{"uri":"/repositories/3"}' # remove repository ID 3
```

This is the general pattern also used with subject, agent, accession, digital_object.


The _aspace_ command uses the following environment variables

+ ASPACE_API_URL, the URL to the ArchivesSpace API (e.g. http://localhost:8089 in v1.4.2)
+ ASPACE_USERNAME, username to access the ArchivesSpace API
+ ASPACE_PASSWORD, to access the ArchivesSpace API
+ ASPACE_DATASET, the directory for exported content

### _aspacepage_

This command generates static webpages from exported ArchivesSpace content.

It relies on the following environment variables

+ ASPACE_DATASET, where you've exported your ArchivesSpace content
+ ASPACE_HTDOCS, where you want to write your static pages
+ ASPACE_TEMPLATES, the templates to use (this defaults to template/defaults but you probably want custom templates for your site)

The typical process would use _aspace_ to export all your content and then run _aspacepage_ to generate your website content.

```
    ./bin/aspace instance export # this takes a while
    ./bin/aspacepage # this is faster
```

Assuming the default settings you'll see new webpages in your local *htdocs* directory.


### _aspaceindexer_

This command creates [bleve](http://blevesearch.com) indexes for use by _aspacesearch_.

Current _aspaceindexer_ operates on JSON content exported with _aspace_. It expects
a specific directory structure with each individual JSON blob named after its
numeric ID and the extension .json. E.g. data/repositories/2/accession/1.json would
correspond to accession id 1 for repository 2.

_aspaceindexer_ depends on four environment variables

+ ASPACE_DATASET, the root directory where the JSON blobs are saved
+ ASPACE_BLEVE_INDEX, the name of the Bleve index (created or maintained)
+ ASPACE_BLEVE_MAPPING, the name of the Bleve map file (assuming you're not using the default mapping)

### _aspacesearch_

_aspacesearch_ provides both a static webserver as well as web search service.

Current _aspacesearch_ uses the Bleve indexes created with _aspaceindexer_. It also
uses the search page and results templates defined in ASPACE_TEMPLATES.

It uses the following environment variables

+ ASPACE_BLEVE_INDEX, the Bleve index to use to drive the search service
+ ASPACE_TEMPLATES, templates for search service as well as browsable static pages
+ ASPACE_SEARCH_URL, the url you want to run the search service on (e.g. http://localhost:8501)

Assuming the default setup, you could start the like

```
    ./bin/aspacesearch
```

Or you could add a startup script to /etc/init.d/ as appropriate.
