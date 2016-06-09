
# Search for public website

_cait_ provides a public search interface for searching accessions data.It used the [Bleve](http://blevesearh.com) library for
implementing search. By default bleve implements a [tf–idf](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) algorithm for search results.
This works well for specific terms and term combination but can yield unexected orderings when look for single names (e.g. a family name).
