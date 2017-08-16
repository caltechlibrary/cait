
# Action Items

## Bug

+ [ ] Indexing dates is failing to be handled properly (need to use a some logic and derive a useful value like with templated defs in dsindexer)
    + Two problems, dates need to be normalized for all varieties (e.g. single, single partial, inclusive, inclusive partial dates)
    + Document type isn't identitied so the Mappings aren't overriding the default indexing, this is why date sorts fail

## Next (Sprint)

+ [ ] Develop an Agent/Person Template, include name, bio and links to accessions if available
+ [ ] Create an accessions report for Loma to answer Patrons requests (needs to include both Accession, list of persons and link to the access record)0:w

## Bugs (Sprint)

+ [ ] Get to the bottom of why keys.json have zero items after export
    + look at Create() versus Open() and come up with a better ApiCollection shim.
+ [ ] Aprox., Circa, C.E. dates need to be formated correctly in archives.caltech.edu website
+ [ ] Single Date rendering bug


## Some day, Maybe list

+ [ ] Add sortable results
+ [ ] Add support to core cait for resource objects, archival_objects, etc.
+ [ ] Implement incremental update support for AS export (see Humdol plugin at Github)
+ [ ] Add harvesting of agents/corporate entity
+ [ ] Migrate from cait-indexer to mkpage's general purpose indexer
+ [ ] Migrate from cait-servepages to mkpage's ws embedded search enabled
+ [ ] Add support to auto-render useful shared Google spreadsheet(s) based on AS harvest
+ [ ] Modernize website design (e.g. use CSS for layout instead of tables)
 
## Completed

+ [x] Add harvesting of agents/person
+ [x] Search results lists have character [encoding issues](http://archives.caltech.edu/search/basic/?q=Marble&-search.x=7&-search.y=0)
+ [x] Update harvest representation to use [dataset](https://caltechlibrary.github.io/dataset) package
    + io.go has a WriteJSON that can be used to create/write JSON to collections (e.g. repositories/2/accessions is a collection)
    + view.go has various ioutil.ReadDir() and ioutil.ReadFile() that could be replaced with either dataset functions for getting keys or get data
    + [x] Store dataset in S3 rather than on local dist, update processing to support S3
+ [x] Some templates have meta tag setting charset explicitly to iso-8859-1, remove or switch to utf-8
