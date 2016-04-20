
# General development notes

## Running things in a production setting

### Example nightly update

```shell
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
    # Run archives site update everyday at 6:30am.
    30 6 * * * /archivesspace/cait/bin/nightly-update.sh >> /archivesspace/logs/nightly-update.log 2>&1
```


## Reference Links

+ [Explanation of authentication](https://github.com/archivesspace/archivesspace/blob/4c26d82b1b0e343b7e1aea86a11913dcf6ff5b6f/docs/slate/source/index.md#authentication)
+ [File API](https://archivesspace.github.io/archivesspace/doc/file.API.html)
+ [General Docs](https://archivesspace.github.io/archivesspace/)
+ [API Docs](http://archivesspace.github.io/archivesspace/api/) (dynamically generated so doesn't show in Google search results)
+ [Wiki](https://archivesspace.atlassian.net/wiki/display/ADC/ArchivesSpace) 

