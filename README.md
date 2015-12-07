
# gospace

## Golang ArchivesSpace package and utility

A proof of concept Golang package for working with ArchivesSpace REST API.

## Requires

+ A working ArchivesSpace instance reachable on the network
+ Golang 1.5 or better to compile

If you want to run the shell scripts used to prototype the API interactions then you'll also need Bash, curl, cut, tr, sed, and jq.
See NOTES.md for more details.


## Initial prototypes implemented via bash, curl, cut, sed, tr, and jq

The scripts in the shell directory were used to dump examples from a populated ArchivesSpace instance for determining the appropriate structures needed in Go. They are retained in this repository because they are useful and illustrate basic ArchivesSpace File API concepts.

+ shell/get-agents.sh - will dump all the agent records for people, corporate_entities, families and software
+ shell/get-repositories.sh - will dump the currently defined repositories in the target ArchivesSpace instance
+ shell/get-accessions.sh - will dump all the accessions for a specific repository

Data collect by these scripts are placed in the data/* directory hierarchy created when the scripts are run.



## aspace command examples

Current aspace supports operations on repositories. It supports

+ create
+ list (individually or all repositories)
+ update (uses a JSON blob generated from listing a specific repository)
+ delete

Here's an example of using the _aspace_ command line tool

```shell
    . setup.conf # Load the connection info into the environment
    aspace repository create "My Archive" "This is an example of my archive"
    aspace repository list all # show a list of archives, for example purposes we'll use archive ID of 11
    aspace repository list 11   # Show only the archive JSON for repository ID equal to 11,
    # Example output is {"id":11,"repo_code":"My Archive","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"    ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","s    ystem_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    # Change 'My Archive to Test Archives'
    aspace repository update {"id":11,"repo_code":"Test Archives","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","system_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    aspace repository list 11 # See the update output for repo ID 11
    aspace repository delete 11 # remove repository ID 2
```

The file _setup.conf_ is a shell script which exports the needed environment variables -

+ ASPACE_PROTOCOL
+ ASPACE_HOST
+ ASPACE_PORT
+ ASPACE_USERNAME
+ ASPACE_PASSWORD

If you're shell is back then @. setup.conf@ setups the variables for a file that looks like

```
    #
    # this is an example setup configuration for running the API tests.
    #

    #
    # Local Development setup
    #
    export ASPACE_PROTOCOL=http
    export ASPACE_HOST=localhost
    export ASPACE_PORT=8089
    export ASPACE_USERNAME=admin
    export ASPACE_PASSWORD=admin
```
