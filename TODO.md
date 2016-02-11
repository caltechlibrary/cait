
# To do and action items

## Action Item

### Bugs

+ Need to determine exact Bleve var that needs to change to page through results
+ Need to strip unneeded characters from input to prevent XSS attacks (or figure out how to use html/template and still get highlight)
+ Need a better title sort (e.g. remove stop words)
+ Empty Search results cause broken page
+ Simple Search Page
    + Missing Link to Advanced Search
    + Missing Link to browse "Manuscript Collections" and "Oral History"
+ Result Pages missing elements
    + Missing link back to search/advanced search page
    + "Online File" link in results (Digital Object File Version link)
    + Extent info
    + Created date
    + Created by
    + "collection" subject (e.g. Manuscript Collection, Oral History)
+ Detail Pages missing element
    + Title
    + "Online File" link in results (Digital Object File Version link)
    + Extent info
    + Created date
    + Created by
    + "collection" subject (e.g. Manuscript Collection, Oral History)
    + Previous, Next for browsing
    + Subjects (e.g. Manuscript Collection, Oral History)
    + Access and use
    + Description of Contents
+ Search results need paging values (modify view or create a nav object?)
+ Simple search and Advanced submit empty search string causes a Go panic
    + side effect is broken page when clicking on detail view page
    + Default behavior should be to go in browsing mode by title

# some day maybe list

+ create a JSON to xlsx utility
+ create an digital_object or resources to EAD utility
