
# General deployment notes

For illustration purposes the deployment directory and site URL,
and release version are

+ /Sites/archives.example.edu
+ http://archives.example.edu
+ v0.0.9

Overview steps taken

1. Get the release zip file from http://github.com/caltechlibrary/cait/releases/latest
2. unzip the release file into a temporary directory
3. copy the binaries for the appropriate architecture (e.g. linux-amd64) to an appropraite bin directory (e.g. /Sites/archives.example.edu/bin)
4. copy, modify, and source the example configuration file (e.g. etc/cait.bash-example to /etc/cait.bash)
5. copy and modify scripts/update-website.bash (if using Github webhooks)
6. copy and modify scripts/nightly-update.bash for running under cron
7. Test everything works
9. If everything is OK then setup cronjob

Example shell commands run

```shell
    # Step 1
    curl -O https://github.com/caltechlibrary/cait/releases/download/v0.0.8/cait-binary-release.zip
    # Step 2
    mkdir -p tmp && cd tmp
    unzip cait-binary-release.zip
    # Step 3
    mkdir -p /Sites/archives.example.edu/bin
    cp -v dist/linux-amd64/* /Sites/archives.example.edu/bin/
    # Step 4
    cp -v etc/cait.bash-example /Sites/archives.example.edu/etc/cait.bash
    # e.g. setup the value of $HOME to /Sites/archives.example.edu
    # If needed include /Sites/archives.example.edu in PATH
    vi /Sites/archives.example.edu/etc/cait.bash
    . /Sites/archives.example.edu/etc/cait.bash
    # Step 5
    cp -v scripts/update-website.bash /Sites/archives.example.edu/bin/
    vi /Sites/archives.example.edu/bin/update-website.bash
    # Step 6
    cp -v scripts/nightly-update.bash /Sites/archives.example.edu/bin/
    # e.g. Set the value of HOME to /Sites/archives.example.edu
    vi /Sites/archives.example.edu/bin/nightly-update.bash
    # Step 7
    cait -v
    cait-genpages -v
    cait-indexpages -v
    cait-servepages -v
    scripts/update-website.bash
    scripts/nightly-update.bash
    # Step 8
    # Add the cronjob for /Sites/archives.example.edu/scripts/nightly-update.bash
    cronjob -e
    # List the cronjob and verify it is correct.
    cronjob -l
```

## Example cronjob

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
    # Run archives site update everyday at 6:30am.
    30 6 * * * /Sites/archives.example.edu/scripts/nightly-update.bash >> /Sites/archives.example.edu/logs/nightly-update.log 2>&1
```

## Reference Links

+ [Explanation of authentication](https://github.com/archivesspace/archivesspace/blob/4c26d82b1b0e343b7e1aea86a11913dcf6ff5b6f/docs/slate/source/index.md#authentication)
+ [File API](https://archivesspace.github.io/archivesspace/doc/file.API.html)
+ [General Docs](https://archivesspace.github.io/archivesspace/)
+ [API Docs](http://archivesspace.github.io/archivesspace/api/) (dynamically generated so doesn't show in Google search results)
+ [Wiki](https://archivesspace.atlassian.net/wiki/display/ADC/ArchivesSpace)
+ [Duke's Python scripts for ArchivesSpace](https://github.com/noahgh221/archivesspace-duke-scripts)
