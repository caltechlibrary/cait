
# Curl Examples for working with ArchivesSpace

The ruby centric docs are in the ArchicesSpace [File API](https://archivesspace.github.io/archivesspace/doc/file.API.html)

## Useful tools

Found these tools in _Data Science at the Command Line_. For a complete list check "Appendix A. A Complete List of Tools"

+ Python based
    + [csvkit](http://csvkit.readthedocs.org/en/0.9.1/) various tools for working with CSV files (e.g. csvcut, csvgrep, csvlook, csvjson)
    + Pretty Print JSON output by piping through _python -m json.tool_
+ Golang based
    + [json2csv](https://github.com/jehiah/json2csv) convert JSON streams into CSV format (Golang)
+ NodeJS based
    + [jq](https://github.com/stedolan/jq/) A command line JSON selector analogous to XPath for XML
    + [xml2json](https://github.com/parmentf/xml2json) convert XML to JSON

### ArchivesSpace tools

+ [ashtin](https://www.npmjs.com/package/ashtin) a NPM module for working with Solr
    + [Project at Github](https://github.com/quoideneuf/ashtin)
    + Might be interesting to fold some of this into aspace

Install ashtin

```
    npm install ashtin -g
    ashtin setup
    ashtin prune-index --solr_url http://localhost:8090
```

## Authenticating

Assuming user is _admin_ and password is _admin_ for this example.

```
    curl -Fpassword=admin http://localhost:8089/users/admin/login
```

Save the authorize token (i.e. session value) in an environment variable for use in other requests.

See [YouTube](https://www.youtube.com/watch?v=iKd4ZME1uIE) at about 1:30.

Here's a one liner to login and save the session value as an environment variable _TOKEN_

```
    export TOKEN=$(curl -Fpassword=admin http://localhost:8089/users/admin/login | jq -r '.session')
    echo "Auth token is $TOKEN"
```


## Repositories

### List All Repositories

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/repositories
```

+ Adding  _| python -m json.tool_ to the above command will pretty print the results to the console.


## Agents

### List All Agents/People ids

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/agents/people?all_ids=true
```

### List specific Agents/People by Id

List the agent with id of 2

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/agents/people/2
```

### List All Agents/Corporate Entities

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/agents/corporate_entities?all_ids=true
```


### List a specific Agent in relation to repository 3

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/repositories/3/agents
```


## Accessions

Accessions are accessed in the context of a specific repository only.

### List all accession ids in repository 3

```
    curl -H "X-ArchivesSpace-Session: $TOKEN" http://localhost:8089/repositories/3/accessions?all_ids=true
```

# Refactoring aspace.go

Things to consider

+ With the exception of Get*() other methods on API should "return *ResponseMsg, string, err" where string is the Unmarshal value of the returned requested content
+ Should args to GetAgents() and CreateAgents() pass in the Agent.URI value or a string representing agent type?
