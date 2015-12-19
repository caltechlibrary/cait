
# aspace

## Golang ArchivesSpace package and utility

A proof of concept Golang package for working with ArchivesSpace REST API.

## Requires

+ A working ArchivesSpace instance reachable on the network
+ Golang 1.5 or better to compile

See NOTES.md for more details.

If you want to run the shell scripts see [github.com/caltechlibrary/aspace-shell-scripts](https://github.com/caltechlibrary/aspace-shell-scripts).

You can setup the environment to use the _aspace_ command by sourcing _shell/api-login.sh_

```
    . shell/api-login.sh # Answer the prompts to set things up
```


## _aspace_ command examples

Current _aspace_ supports operations on repositories. It supports

+ create
+ list (individually or all repositories)
+ update (uses a JSON blob generated from listing a specific repository)
+ delete

Here's an example of using the _aspace_ command line tool

```shell
    . shell/api-login.sh # Load the connection info into the environment
    aspace repository create "My Archive" "This is an example of my archive"
    aspace repository list all # show a list of archives, for example purposes we'll use archive ID of 11
    aspace repository list 11   # Show only the archive JSON for repository ID equal to 11,
    # Example output is {"id":11,"repo_code":"My Archive","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"    ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","s    ystem_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    # Change 'My Archive to Test Archives'
    aspace repository update {"id":11,"repo_code":"Test Archives","name":"This is an example of my archive","uri":"/repositories/11","agent_representation":{"ref":"/agents/corporate_entities/9"},"image_url":"","lock_version":1,"created_by":"admin","last_modified_by":"admin","create_time":"2015-12-01T00:52:55Z","system_time":"0001-01-01T00:00:00Z","user_mtime":"2015-12-01T01:00:29Z"}
    aspace repository list 11 # See the update output for repo ID 11
    aspace repository delete 11 # remove repository ID 2
```

The _aspace_ commands uses the following environment variables

+ ASPACE_API_URL
+ ASPACE_API_TOKEN
+ ASPACE_USERNAME
+ ASPACE_PASSWORD

You could put these in a shell script such as _setup.sh_ and source them `. setup.sh`

```
    #
    # this is an example setup configuration for running the API tests.
    #

    #
    # Local Development setup
    #
    export ASPACE_API_URL=http://localhost:8089
    export ASPACE_USERNAME=admin
    export ASPACE_PASSWORD=admin
```
