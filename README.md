
# cait

[cait](https://github.com/caltechlibrary/cait) is a set of utilities written in the [Go](http://golang.org) language that work with and augment the [ArchivesSpace](http://archivesspace.org) API.

+ cait - a command line utility for ArchivesSpace interaction (basic CRUD operations and export)
+ caitpage - a simple static page generator based on exported ArchivesSpace content
+ caitindexer - for indexing exported JSON structures with [Bleve](https://github.com/blevesearch/bleve)
+ caitserver - a web service providing public search services and content browsing
+ xlsximporter - a tool for turning Excel spreadsheets in .xlsx format into JSON files suitable for importing into ArchivesSpace
+ sitemapper - a simple tool to generate a sitemap.xml file from pages rendered with caitpage

## Requirements

+ A working deployment of ArchivesSpace
+ Golang 1.5.3 or better to compile
+ Three 3rd party Go packages
    + [Bleve](https://github.com/blevesearch/bleve) by [Blevesearch](http://blevesearch.com), Apache License, Version 2.0
    + [Otto](https://github.com/robertkrimen/otto) by Robert Krimen, MIT license
    + [xlsx](https://github.com/tealeg/xlsx) by Tealeg, BSD license

## Compiling

If you already have [Go](https://golang.org) setup and installed compiling the utilities are pretty straight forward.

1. Clone the git repository for the project
2. "Go get" the 3rd party libraries
3. Compile
4. Setup the necessary environment variables for using the utilities

Here's a typical example of setting things up.

```
    git clone git@github.com:caltechlibrary/cait.git
    cd cait
    go get -u github.com/blevesearch/bleve/...
    go get -u github.com/robertkrimen/otto
    go get -u github.com/tealeg/xlsx
    mkdir bin
    go build -o bin/cait cmds/cait/cait.go
    go build -o bin/caitpage  cmds/caitpage/caitpage.go
    go build -o bin/caitindexer cmds/caitindexer/caitindexer.go
    go build -o bin/caitserver cmds/caitserver/caitserver.go
    go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go
    go build -o bin/sitemapper cmds/sitemapper/sitemapper.go
```

At this point you should have your command line utilities ready to go in the *bin* directory. You are now ready to setup your environment variables.


## Setting up your environment

The command line tools and services are configured via environment variables. Below is an example of setting things up under Bash running on your favorite Unix-like system.


```
    #
    # setup.sh - this script sets the environment variables for cait project.
    # You would source file before using cait, caitindexer, or caitserver.
    #

    #
    # Local Development setup
    #
    export CAIT_API_URL=http://localhost:8089
    export CAIT_USERNAME=admin
    export CAIT_PASSWORD=admin
    export CAIT_DATASETS=data
    export CAIT_SEARCH_URL=http://localhost:8501
    export CAIT_HTDOCS=htdocs
    export CAIT_TEMPLATES=templates/default
    export CAIT_BLEVE_INDEX=index.bleve

    #
    # Create the necessary directory structure
    #
    mkdir -p $CAIT_DATASETS
    mkdir -p $CAIT_HTDOCS
    mkdir -p $CAIT_TEMPLATES

```

Assuming Bash and that you've named the file _setup.sh_ you could
source the file from your shell prompt by typing

```
    . setup.sh
```

### Setting up a dev box

I run ArchivesSpace in a vagrant box for development use. You can find details to set that up at [github.com/caltechlibrary/archivesspace_vagrant](https://github.com/caltechlibrary/archivesspace_vagrant).
I usually run the [cait](https://github.com/caltechlibrary/cait) tools locally. You can see
and example workflow in the document [EXPORT-IMPORT.md](EXPORT-IMPORT.md).

## Utilities

### _cait_

This command is a general purpose tool for fetch ArchivesSpace data from the
ArchivesSpace REST API, saving or modifying that data as well as querying the
locally capture output of the API.

Current _cait_ supports operations on repositories, subjects, agents, accessions and digital_objects.

These are the common actions that can be performed

+ create
+ list (individually or all ids)
+ update (can use a file instead of the command line, see -i option)
+ delete
+ export (useful with integrating into static websites or batch processing via scripts)

Here's an example session of using the _cait_ command line tool on the repository object.

```shell
    . setup.sh # Source my setup file so I can get access to the API
    cait repository create '{"uri":"/repositories/3","repo_code":"My Archive","name":"My Archive"}' # Create an archive called My Archive
    cait repository list # show a list of archives, for example purposes we'll use archive ID of 3
    cait repository list '{"uri":"/repositories/3"}' # Show only the archive JSON for repository ID equal to 3
    cait repository list '{"uri":"/repositories/3"}' > repo2.json # Save the output to the file repo3.json
    cait repository update -i repo3.json # Save your changes back to ArchivesSpace
    cait repository export '{"uri":"/repositories/3"}' # export the repository metadata to data/repositories/3.json
    cait repository delete '{"uri":"/repositories/3"}' # remove repository ID 3
```

This is the general pattern also used with subject, agent, accession, digital_object.


The _cait_ command uses the following environment variables

+ CAIT_API_URL, the URL to the ArchivesSpace API (e.g. http://localhost:8089 in v1.4.2)
+ CAIT_USERNAME, username to access the ArchivesSpace API
+ CAIT_PASSWORD, to access the ArchivesSpace API
+ CAIT_DATASET, the directory for exported content

### _caitpage_

This command generates static webpages from exported ArchivesSpace content.

It relies on the following environment variables

+ CAIT_DATASET, where you've exported your ArchivesSpace content
+ CAIT_HTDOCS, where you want to write your static pages
+ CAIT_TEMPLATES, the templates to use (this defaults to template/defaults but you probably want custom templates for your site)

The typical process would use _cait_ to export all your content and then run _caitpage_ to generate your website content.

```
    ./bin/cait archivesspace export # this takes a while
    ./bin/caitpage # this is faster
```

Assuming the default settings you'll see new webpages in your local *htdocs* directory.


### _caitindexer_

This command creates [bleve](http://blevesearch.com) indexes for use by _caitserver_.

Current _caitindexer_ operates on JSON content exported with _cait_. It expects
a specific directory structure with each individual JSON blob named after its
numeric ID and the extension .json. E.g. data/repositories/2/accession/1.json would
correspond to accession id 1 for repository 2.

_caitindexer_ depends on four environment variables

+ CAIT_DATASET, the root directory where the JSON blobs are saved
+ CAIT_BLEVE_INDEX, the name of the Bleve index (created or maintained)
+ CAIT_BLEVE_MAPPING, the name of the Bleve map file (assuming you're not using the default mapping)

### _caitserver_

_caitserver_ provides both a static web server as well as web search service.

Current _caitserver_ uses the Bleve indexes created with _caitindexer_. It also
uses the search page and results templates defined in CAIT_TEMPLATES.

It uses the following environment variables

+ CAIT_BLEVE_INDEX, the Bleve index to use to drive the search service
+ CAIT_TEMPLATES, templates for search service as well as browsable static pages
+ CAIT_SEARCH_URL, the url you want to run the search service on (e.g. http://localhost:8501)

Assuming the default setup, you could start the like

```
    ./bin/caitserver
```

Or you could add a startup script to /etc/init.d/ as appropriate.

### _xlsximporter_

_xlsximporter_ is a utilty to transform sheets from an Excel file in xlsx format to JSON blobs suitable for importation into ArchivesSpace (e.g. Digital Objects).  By default it transforms each row in the spreadsheet into an object where the property names correspond to the column headers (in the initial row of the spreadsheet).  You can perform more elaborate mappings using a javascript callback function.  You can see an example of that in the *xlsximporter-javascript-example* directory.

The general workflow would be to transform you rows to JSON objects on your local disc with _xlsximporter_ then use the _cait_ utility to push the JSON blobs into ArchivesSpace itself.

### _sitemapper_

_sitemapper_ generates a sitemap.xml file based on the arguments you envoke. It you're site's URL is http://archives.example.edu, your htdocs directory is htdocs and you want to save you sitemap.xml file as htdocs/sitemap.xml you could run the command with

```
    sitemapper htdocs htdocs/sitemap.xml http://archives.example.edu
```

This will generate a site map of the HTML files found in *htdocs* with the results saved in *htdocs/sitemap.xml*. For more informaiton about sitemaps see the [sitemaps.org](http://sitemaps.org) website.
