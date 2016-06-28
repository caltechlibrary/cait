
# Date types

The templates for rendering website pages (via genpages and servepages) need to
accommodate various date types in "dates" fields.


## single

```
    "dates": [
        {
          "date_type": "single",
          "label": "creation",
          "expression": "1757-01-01",
          "lock_version": 0,
          "jsonmodel_type": "date",
          "created_by": "admin",
          "last_modified_by": "admin",
          "user_mtime": "2015-10-19T22:45:20Z",
          "system_mtime": "2015-10-19T22:45:20Z",
          "create_time": "2015-10-19T22:45:20Z"
        }
    ]
```

## inclusive

```
    "dates": [
        {
          "date_type": "inclusive",
          "label": "creation",
          "expression": "Date Created",
          "begin": "1804-01-01",
          "end": "1804-12-31",
          "lock_version": 0,
          "jsonmodel_type": "date",
          "created_by": "admin",
          "last_modified_by": "admin",
          "user_mtime": "2015-10-19T22:45:20Z",
          "system_mtime": "2015-10-19T22:45:20Z",
          "create_time": "2015-10-19T22:45:20Z"
        }
    ]
```

## bulk

```
    "dates": [
        {
          "date_type": "bulk",
          "label": "creation",
          "certainty": "approximate",
          "begin": "1940",
          "end": "1941",
          "lock_version": 0,
          "jsonmodel_type": "date",
          "created_by": "mariella",
          "last_modified_by": "mariella",
          "user_mtime": "2016-03-02T23:57:58Z",
          "system_mtime": "2016-03-02T23:57:58Z",
          "create_time": "2016-03-02T23:57:58Z"
        }
    ]
```

## Finding records with date types

If you have [jq](https://stedolan.github.io/jq) and installed then it is easy
to scan the dataset directory for examples in your accessions.

### Finding all single dates entries

```shell
    find dataset/repositories/2/accessions -type f |\
    while read ITEM; do
        jq '{"accession_id": .id, "dates":.dates}| select(.dates[].date_type == "single")' "$ITEM";
    done
```

### Finding all inclusive dates entries

```shell
    find dataset/repositories/2/accessions -type f |\
    while read ITEM; do
        jq '{"accession_id": .id, "dates":.dates}| select(.dates[].date_type == "inclusive")' "$ITEM";
    done
```

### Finding all bulk dates entries

```shell
    find dataset/repositories/2/accessions -type f |\
    while read ITEM; do
        jq '{"accession_id": .id, "dates":.dates}| select(.dates[].date_type == "bulk")' "$ITEM";
    done
```


## Interesting fields

+ label
+ expression (free text)
+ date_type
    + inclusive
    + single
    + bulk
+ begin
    + formats: YYYY, YYYY-MM, YYYY-MM-DD
+ end
    + formats: YYYY, YYYY-MM, YYYY-MM-DD
+ certainty
    + ""
    + approximate
    + inferred
    + questionable
+ era
    + ""
    + ce
+ calendar
    + ""
    + Georgian

