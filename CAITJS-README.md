
# caitjs

_caitjs_ is a command line JavaScript runner for working with the ArchivesSpace API.
_caitjs_ exposes some of the exported API calls implemented int the cait Golang package 
to JavaScript. This makes it easy to implement ad-hoc scripts to process data, as well
as import/export content.

In addition to the ArchivesSpace API functions there are three additional functions
which maybe of some use

+ Getenv() which allows you to query the contentsof the shell environment where caitjs is running
+ HttpGet() which performs a HTTP GET request
+ HttpPost() which performs a HTTP POST request

## Examples

In the examples directory under the caitjs subdirect you'll find a number of scripts
demonstrating the usage of the API features.


