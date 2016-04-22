
# cait JavaScript implementation

_cait_ is a command line JavaScript runner and repl for working with the ArchivesSpace API.
It provides access via JavaScript to some of the exported API methods implemented int the _cait_ package.
This makes it easy to implement ad-hoc scripts to process data, as well as import/export content.

## Interactive

When run in interactive mode (e.g. `cait -i`) you can use autocomplete to explore the additional non-standard
JavaScript objects available (i.e. os, http, api). To active command line completion use Ctrl+I (or tab)
and to cancel command line completion use Ctrl+G at the command prompt '> '.

The interactive mode provides a simple help command as well as comment line completion of object and methods names using your *tab* key.

The JavaScript shell also provides support for multi-line statements.  Multi-line statements can be aborted with the *.break*
command.

The shell can be exited with the *quit()* function.

The shell displays the results of the evaluated expression in bold.


## Help

In interactive mode you can get help by using the help function

+ help() - list the objects available in the JavaScript environment
+ get object level help by providing the object name as a string
    + help("api") - would list help for the *api* object
+ get help on an object function by providing the object name, dot, and function name as a string
    + help("api.login") - would list help information about logging into the API

You can also press the *tab* key to get a list of possible completions. Those completions can be navigated around
with your arrow keys.

## ArchivesSpace API object and methods

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
+ api.createResource(),  creates a resource
+ api.getResource(), gets a resource by id
+ api.updateResource(),  updates a resource
+ api.deleteResource(),  deletes a resource
+ api.listResources(),  lists all the resource ids

## Object and methods for working with the file system and environment

In addition to the ArchivesSpace API functions there are three additional functions
which maybe of some use

+ os.args() return an array of command line arguments
+ os.getEnv() which allows you to query the contents of the shell environment where *cait* is running
+ os.exit() exit the program with a given error code (defaults to 0)
+ os.chmod() return true for success false otherwise
+ os.readFile() which reads in the contents of a file as text
+ os.writeFile() which writes the contents to a file as text
+ os.rename() rename old path to new path
+ os.remove() remove a file
+ os.find() return an array of a directory's content
+ os.mkdir() create a directory
+ os.mkdirAll() create a directory and sub-directories needed
+ os.rmdir() remove a directory
+ os.rmdirAll() remove a directory and contents

## Object and methods for working with HTTP resources

+ http.get() which performs a HTTP GET request
+ http.post() which performs a HTTP POST request

They're similar JavaScript functions available in _xlsximporter_.

## Examples

In the [examples directory](./examples) under the cait sub-direct you'll find a number of scripts
demonstrating the usage of the API features. These should be treated as proofs of concepts and exploratory
as they were developed primarily in response to figuring out how _cait_ should work.

+ examples/convert-linked-agents-creators-to-subject.js - shows how you could change an agents' type in an accession record from creator to subject
+ examples/fix-agents-people.js - demonstrates modifying agents people record
+ examples/fix-digital_object_id.js - demontstrates modifying digitalobject identifiers
+ examples/get-repository.js - demonstrates listing repositories available
+ examples/helloworld.js - your basic "Hello World" example
+ examples/link-digital-objects-to-accessions.js - demonstrates linking digital objects to accessions with matching titles
+ examples/list-date-types-in-accessions.js - list the types of dates associated with accessions (very slow running)
+ examples/list-linked-agents-roles.js - list the types of roles that agents have linked to accessions
+ examples/signon-signoff.js - demonstrate api.login() and api.logout() function
