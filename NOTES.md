
# General development notes

## Installation

1. Build a new release (or pick one available at https://github.com/caltechlibrary/cait/latest). 
2. Copy the zip archive to the production machine
3. Extract the appropriate binaries (e.g. linux-amd42/*) to an appropriate bin directory (e.g. /Sites/archives.example.edu/bin)
4. Create a configuration and source it to your environment (e.g. `. /etc/cait.bash`)
5. Test the commands (e.g. cait, genpages, indexpages, sitemapper, servepages)
6. Setup cronjob to run and harvest ArchivesSpace content (see below)
7. Manually run your nightly script and watch the logs for errors.
8. Fix deployment bugs as needed.

Here's an example of the commands for install v0.0.8 on a Linux machine
where the deployment directory is /Sites/archives.example.edu.

```shell
    curl -O https://github.com/caltechlibrary/cait/releases/download/v0.0.8/cait-binary-release.zip
    unzip cait-binary-release.zip
    cp -v dist/linix-amd64/* /Sites/archives.example.edu/bin/
    sudo cp -v dist/etc/setup.conf-example /etc/cait.bash
    # Edit /etc/cait.bash to make sense for your envinronment
    # E.g. add "export PATH=/Sites/archives.example.edu/bin:$PATH"
    # for this example and set the other CAIT_* variables.
    sudo vi /etc/cait.bash
    # Source the configuration into your environment
    . /etc/cait.bash
    cait -v
    genpages -v
    indexpages -v
    sitemapepr -v
    serverpages -v
    # Now test your nightly script, watch the output for problems
    /Sites/archives.example.edu/scripts/nightly-update.bash
```

## Running things in a production setting

### Example nightly update

This is an example script that could be run as a nightly cronjob. Output from the cait tools is suitable to sending to a log file (e.g. /Sites/archives.example.edu/logs/nightly-update.log)

```shell
    #!/bin/bash
    #

    # This is an example cronjob to be run from the root account.


    # Load the cait configuration
    . /etc/cait.conf

    # Change directory to where cait is installed
    cd /Sites/archives.example.edu/
    # Export the current content from ArchivesSpace
    ./bin/cait archivesspace export
    # Generate webpages
    ./bin/genpages
    # Index webpages
    ./bin/indexpages

    # You should now be ready to reload the search engine
    /etc/init.d/servepages stop
    /etc/init.d/servepages start
```

### Example cronjob

```shell
    #!/bin/bash
    #
    #  field         allowed values
    #  -----         --------------
    #  minute        0-59
    #  hour          0-23
    #  day of month  1-31
    #  month         1-12 (or names, see below)
    #  day of week   0-7 (0 or 7 is Sun, or use names)
    #
    # Run archives site update everyday at 3:00am.
    0 3 * * * /Sites/archives.example.edu/scripts/nightly-update.sh >> /Sites/archives.example.edu/logs/nightly-update.log 2>&1
```


## Reference Links

+ [Explanation of authentication](https://github.com/archivesspace/archivesspace/blob/4c26d82b1b0e343b7e1aea86a11913dcf6ff5b6f/docs/slate/source/index.md#authentication)
+ [File API](https://archivesspace.github.io/archivesspace/doc/file.API.html)
+ [General Docs](https://archivesspace.github.io/archivesspace/)
+ [API Docs](http://archivesspace.github.io/archivesspace/api/) (dynamically generated so doesn't show in Google search results)
+ [Wiki](https://archivesspace.atlassian.net/wiki/display/ADC/ArchivesSpace)
+ [Duke's Python scripts for ArchivesSpace](https://github.com/noahgh221/archivesspace-duke-scripts)
