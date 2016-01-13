
# aspace

## Golang ArchivesSpace REST API package and utilities

_aspace_ is a proof of concept Golang package for working with the ArchicesSpace
REST API. In addition it provides several (hopefully) useful tools

+ aspace - a command line utility for ArchivesSpace interaction (CRUD and Search)
+ aspaceindexer - for indexing exported JSON structures with [Bleve](https://github.com/blevesearch/bleve)
+ aspacesearch - a web service providing public search services and content browsing
+ aspacedashboard - a web service providing administrative search services

Most the ArchivesSpace REST API for Agents and Accessions have be implemented in the
aspace package (aspace.go, models.go). It should be possible to support the full API
give time to create all the JSON models and appropriate test code.


## Requires

+ A working ArchivesSpace instance reachable on the network
+ Golang 1.5 or better to compile
+ [Bleve](http://blevesearch.com) Golang based Search/Indexing library (think Lucene lite implemented in Golang)

See NOTES.md for more details.

If you want to run the shell scripts see [github.com/caltechlibrary/aspace-shell-scripts](https://github.com/caltechlibrary/aspace-shell-scripts).

You can setup the environment to use the _aspace_ command by sourcing _shell/api-login.sh_

```
    . shell/api-login.sh # Answer the prompts to set things up
```

## configured by environment

The command line tools and services are configured via Unix environment. This is
trivial to set in a shell script and to source into your current environment.
Here's is an example Bash script that would set the environment variables.

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
    export ASPACE_BLEVE_INDEX=index.bleve
    export ASPACE_BLEVE_MAPPING=index.map
    export ASPACE_SEARCH_TEMPLATES=templates/search
    export ASPACE_SEARCH_SITE=public
    export ASPACE_SEARCH_URL=http://localhost:8580
    export ASPACE_DASHBOARD_TEMPLATES=templates/dashboard
    export ASPACE_DASHBOARD_SITE=dashboard
    export ASPACE_DASHBOARD_URL=http://localhost:8581
```

Assuming Bash and that you've named the file _setup.sh_ you could
source the file from your shell prompt by typing `. setup.sh`.

## Utilities

### _aspace_

This command is a general purpose tool for fetch ArchivesSpace data from the
ArchivesSpace REST API, saving or modifying that data as well as querying the
locally capture output of the API.

Current _aspace_ supports operations on repositories, agents, and accessions.
It supports

+ create
+ list (individually or all repositories)
+ update (uses a JSON blob generated from listing a specific repository)
+ delete
+ export (for use with _aspaceindexer_)

Here's an example session of using the _aspace_ command line tool.

```shell
    . setup.sh # Source my setup file so I can get access to the API
    aspace repository create "My Archive" "This is an example of my archive"
    aspace repository list all # show a list of archives, for example purposes we'll use archive ID of 11
    aspace repository list 11   # Show only the archive JSON for repository ID equal to 11,
    # Example output is {"id":11,"repo_code":"My Archive","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"    ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","s    ystem_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    # Change 'My Archive to Test Archives'
    aspace repository update {"id":11,"repo_code":"Test Archives","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","system_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    aspace repository list 11 # See the update output for repo ID 11
    aspace repository delete 11 # remove repository ID 2
```

_aspace_ also supports searching JSON content exported if you've index the content
with _aspaceindexer_.

The _aspace_ command uses the following environment variables

+ ASPACE_API_URL, the URL to the ArchivesSpace API (e.g. http://localhost:8089 in v1.4.2)
+ ASPACE_USERNAME, username to access the ArchivesSpace API
+ ASPACE_PASSWORD, to access the ArchivesSpace API
+ ASPACE_BLEVE_INDEX, (optional) the Bleve index file created with _aspaceindexer_
    + you only need for the _aspace search_ type commands


### _aspaceindexer_

This command creates [bleve](http://blevesearch.com) indexes for using by aspace, aspacesearch and aspacedashboard.

Current _aspaceindexer_ operates on JSON content exported with _aspace_. It expects
a specific directory structure with each individual JSON blob names after its
numeric ID and the extension .json. E.g. Agents/People would be found in
_data/agents/people/_ with filenames like _1.json_ and _2800.json_ depending on
the numeric ID of the person's record.

_aspaceindexer_ depends on four environment variables

+ ASPACE_DATASET, the root directory where the JSON blobs are saved
+ ASPACE_BLEVE_INDEX, the name of the Bleve index (created or maintained)
+ ASPACE_BLEVE_MAPPING, the name of the Bleve map file (assuming you're not using the default mapping)

### _aspacesearch_

This command runs a web service for publicly accessible ASPACE content (e.g. agents and
accessions that are published and not suppressed or restricted).

Current _aspacesearch_ uses the Bleve indexes created with _aspaceindexer_. It also
uses the templates defined in ASPACE_SEARCH_TEMPLATES for rendering the search pages
results and browsable records.

+ ASPACE_BLEVE_INDEX, the Bleve index to use to drive the search service
+ ASPACE_SEARCH_TEMPLATES, templates for search service as well as browsable static pages
+ ASPACE_SEARCH_SITE, static content for the _aspacesearch_ web service
+ ASPACE_SEARCH_PREFIX, the website prefix (e.g. http://locahost:8580/seach)

## _aspacedashboard_ command

Current _aspacedashboard_ uses the Bleve indexes created with _aspaceindexer_. It also
uses the templates defined in ASPACE_DASHBOARD_TEMPLATES for rendering the search pages
results and browsable records.

+ ASPACE_BLEVE_INDEX, the Bleve index to use to drive the dashboard service
+ ASPACE_DASHBOARD_SITESEARCH_TEMPLATES, templates for search service as well as browsable static pages
+ ASPACE_DASHBOARD_SITE, static content for the _aspacedashboard_ web service
