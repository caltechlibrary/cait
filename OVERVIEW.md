
# overview

_cait_ is a go library wrapping the [ArchivesSpace](http://archivesspace.org) REST API.  It includes support for content export, static site generation, indexing and independent search engine and web service.  This means you can manage your content in ArchivesSpace but read and search the public content independent of the status of ArchivesSpace itself.  This gives you more options for deployment as well as providing a clean separation of concerns for public/admin uses.

All tools can be configured through environment variables. Some have additional command line options that can be invoked.  Generally launching the tool with a "-h" or "--help" will get you a list of basic features and options.

## tools

### cait

_cait_ command line utility is the workhorse for getting content out of ArchivesSpace and onto your local file system in a useful static form (JSON blobs).  _cait_ can be used to put some content back into ArchivesSpace. This gives you options for batch editing content with more general tools like R, Open Refine, etc.

### caitjs

_caitjs_ a shell (repl) and JavaScript file runner for working with ArchicesSpace REST API. It provides access to the same API provided internal to _cait_. Helpful for data migrations and fixes.

### genpages

_genpages_ renders the content dumped by _cait_ into a website structure suitable for hosting with _servepages_ search engine and web server.  It does NOT talk directly to ArchivesSpace and as a result does not increase the load on your ArchivesSpace server.

### indexpages

_indexpages_ is a utility for creating and updating a Bleve index used by _servepages_ web server.  It crawls the website tree an ingesting JSON files found in the accessions directories. It can be run manually but is also suited to run periodically via a cronjob (say once every day as needed).   It takes about 45 mimutes to run through my 10k or so of accessions. Your milleage may very. It runs a little faster creating a new index structure than updating.  The current implementation is overly simplistic and certainly can be improved (e.g. rather than indexing files individually it could batch and index)

### servepages

_servepages_ is a web server and search engine. It is intended to run behind a more traditional web server like NginX or Apache.  Output of the search results are controlled by the Golang HTML templates.  This is an early implementation so this will see change as the project gets deployed into a production setting.

_servepages_ can be started manually but more typically would be brought up by your init process (e.g. /etc/init.d/servepages start). An example init file is provided

### xlsximporter

_xlsximporter_ is a utility that reads an Excel file in xlsx format and turns each row of a spreadsheet into a JSON object. By default the properties correspond to the column names but you can also provide a JavaScript file and callback that let's you customize the objects produced. This latter capability is useful when migrating content into ArchivesSpace.

### indexdataset

_indexdataset_ is a utility that crawls the dataset directory (the raw content exported from ArchivesSpace) for use by _searchdataset_.

### searchdataset

_searchdataset_ is a utility for searching the content in the dataset directory (the raw content exported from ArchivesSpace).  Helpful in debugging problems with specific records.


## Workflow for website

1. Export ArchivesSpace content with _cait_
2. Generate pages with _genpages_
3. Index pages with _indexpages_
4. Serve content and search service with _servepages_

## Workflow for data fixes

+ Export ArchivesSpace content as needed with _cait_
+ Index the exported data with _indexdataset_
+ Explore dataset with _searchdataset_
+ Write a JavaScript program to make corrections or use _caitjs_ interactively to fix/update data

## Workflow for importing data from Excel

+ Make your Excel file
+ Map the column and rows to ArchivesSpace objects in JavaScript
+ Run _xlsximporter_ invoking the JavaScript mappings along with spreadsheet name and sheet number
