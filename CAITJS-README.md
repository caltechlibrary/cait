
# caitjs

_caitjs_ is a command line JavaScript runner and repl for working with the ArchivesSpace API.
_caitjs_ exposes some of the exported API calls implemented int the cait Golang package
to JavaScript. This makes it easy to implement ad-hoc scripts to process data, as well
as import/export content.


+ api is the object that holds the *cait api* methods
+ api.login(), logs into the ArchivesSpace API based on the current environment variables
+ api.logout(), logs you out of the ArchivesSpace API
+ api.createRepository() creates a new repository in ArchivesSpace
+ api.getRepository(), gets a repository record
+ api.updateRepository(), updates a repository record
+ api.deleteRepository(), deletes a repository record
+ api.listRepositories(), lists all repositories ids
+ api.createAgent(), creates an agent by type
+ api.getAgent(),  gets an agent by id
+ api.updateAgent(), update an agent
+ api.deleteAgent(), deletes an agent
+ api.listAgents(), lists agent ids by type
+ api.createAccession(), creates an accession
+ api.getAccession(), gets an accession by id
+ api.updateAccession(), updates an accession
+ api.deleteAccession(), deletes an accession
+ api.listAccessions(), lists all accession ids
+ api.createSubject(), creates a subject
+ api.getSubject(), gets a subject by id
+ api.updateSubject(),  updates a subject
+ api.deleteSubject(), deletes a subject
+ api.listSubjects(), lists all subject ids
+ api.createDigitalObject(),  creates a digital object
+ api.getDigitalObject(), gets a digital object by id
+ api.updateDigitalObject(),  updates a digital object
+ api.deleteDigitalObject(),  deletes a digital object
+ api.listDigitalObjects(),  lists all the digital object ids

In addition to the ArchivesSpace API functions there are three additional functions
which maybe of some use

+ os.Getenv() which allows you to query the contentsof the shell environment where caitjs is running
+ http.Get() which performs a HTTP GET request
+ http.Post() which performs a HTTP POST request

These are the same JavaScript functions available in _xlsximporter_.

## Examples

In the examples directory under the caitjs sub-direct you'll find a number of scripts
demonstrating the usage of the API features.
