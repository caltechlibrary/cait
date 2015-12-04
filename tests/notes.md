
# Independant language testing for gospace

Since my plan of a golang based wrapper for ArchivesSpace is non trivial and the documentation about the shapes of the AS models if sparse I am writting some comparison analysis of inputs and output in R so I can determine if I am actually interacting with the AS File API correctly.

## Packages

+ [httr](https://github.com/hadley/httr) - A http client library wrapping RCurl for working with web API
+ [jsonlite](https://github.com/jeroenooms/jsonlite) - A JSON parsing library

