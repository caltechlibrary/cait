//
// Package gospace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2015
// Caltech Library
//
package gospace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ArchivesSpaceAPI is a struct holding the essentials for communicating
// with the ArchicesSpace REST API
type ArchivesSpaceAPI struct {
	URL       *url.URL
	Username  string
	Password  string
	AuthToken string
}

/*
	// Blog post on handling [JSON DATA](http://blog.golang.org/json-and-go)
	Example Repo JSON object

	{
	    "agent_representation": {
	        "ref": "/agents/corporate_entities/8"
	    },
	    "country": "UM",
	    "create_time": "2015-11-20T18:39:49Z",
	    "created_by": "admin",
	    "image_url": "http://identity.example.org",
	    "jsonmodel_type": "repository",
	    "last_modified_by": "admin",
	    "lock_version": 1,
	    "name": "This is a test generated from go_test",
	    "org_code": "Orangization or Agency Code",
	    "parent_institution_name": "Parent Institution Name",
	    "repo_code": "1448044788",
	    "system_mtime": "2015-11-20T18:56:57Z",
	    "uri": "/repositories/9",
	    "url": "http://example.org",
	    "user_mtime": "2015-11-20T18:56:57Z"
	}
*/

// Repository represents an ArchivesSpace repository from the client point of view
type Repository struct {
	ID                    int                    `json:"id"`
	RepoCode              string                 `json:"repo_code"`
	Name                  string                 `json:"name"`
	URI                   string                 `json:"uri,omitempty"`
	URL                   string                 `json:"url,omitempty"`
	AgentRepresentation   map[string]interface{} `json:"agent_representation,omitempty"`
	Country               string                 `json:"country,omitempty"`
	ImageURL              string                 `json:"image_url"`
	OrgCode               string                 `json:"org_code,omitempty"`
	ParentInstitutionName string                 `json:"parent_institution_name,omitempty"`
	LockVersion           int                    `json:"lock_version"`
	CreatedBy             string                 `json:"created_by,omitempty"`
	LastModifiedBy        string                 `json:"last_modified_by,omitempty"`
	CreateTime            time.Time              `json:"create_time,omitempty"`
	SystemTime            time.Time              `json:"system_time,omitempty"`
	UserMTime             time.Time              `json:"user_mtime,omitempty"`
}

// Agent represents an ArchivesSpace agent from the client point of view
type Agent struct {
	ID int `json:"id"`
}

func checkEnv(protocol, host, username, password string) bool {
	if strings.TrimSpace(protocol) == "" {
		return false
	}
	if strings.TrimSpace(host) == "" {
		return false
	}
	if strings.TrimSpace(username) == "" {
		return false
	}
	if strings.TrimSpace(password) == "" {
		return false
	}
	return true
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(protocol, host, port, username, password string) *ArchivesSpaceAPI {
	var connectString string

	aspace := new(ArchivesSpaceAPI)
	if checkEnv(protocol, host, username, password) == false {
		log.Fatal("Cannot create a new ArchivesSpace API connection")
	}

	if port == "" {
		connectString = fmt.Sprintf("%s://%s", protocol, host)
	} else {
		connectString = fmt.Sprintf("%s://%s:%s", protocol, host, port)
	}

	u, err := url.Parse(connectString)
	if err != nil {
		log.Fatalf("ArchivesSpace connection not configured, %s", err)
	}
	aspace.URL = u
	aspace.Username = username
	aspace.Password = password
	aspace.AuthToken = ""
	return aspace
}

// IsAuth returns true if the auth token has been set, false otherwise
func (aspace *ArchivesSpaceAPI) IsAuth() bool {
	if aspace.AuthToken == "" {
		return false
	}
	return true
}

// Login authenticates against the ArchivesSpace REST API setting the AuthToken
// value in the ArchivesSpaceAPI struct.
func (aspace *ArchivesSpaceAPI) Login() error {
	// See https://golang.org/pkg/net/url/#pkg-examples for example building a URL from parts.
	// Command line example: curl -F "password=admin" "http://localhost:8089/users/admin/login"
	var data map[string]interface{}

	aspace.URL.Path = fmt.Sprintf("/users/%s/login", aspace.Username)
	form := url.Values{}
	form.Add("password", aspace.Password)

	res, err := http.PostForm(aspace.URL.String(), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Status != "200 OK" {
		return fmt.Errorf("ArchivesSpace returned HTTP status %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("ArchivesSpace return unreadable body: %s", err)
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("Can't process JSON response %s\n\t%s", content, err)
	}
	aspace.AuthToken = data["session"].(string)
	return nil
}

// Logout clear the authentication token for the session with the API
func (aspace *ArchivesSpaceAPI) Logout() error {
	aspace.URL.Path = `/logout`
	client := &http.Client{}
	req, err := http.NewRequest("GET", aspace.URL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// CreateRepository will create a respository via the REST API for the
// ArchivesSpace instance defined in the ArchivesSpaceAPI struct.
func (aspace *ArchivesSpaceAPI) CreateRepository(repoCode, name string) (*Repository, error) {
	aspace.URL.Path = "/repositories"
	payload := strings.NewReader(fmt.Sprintf(`{"repo_code":%q,"name":%q}`, repoCode, name))

	client := &http.Client{}
	req, err := http.NewRequest("POST", aspace.URL.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("Can't create request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	// content should look something like
	// {"status":"Created","id":3,"lock_version":0,"stale":null,"uri":"/repositories/3","warnings":[]}
	repo := new(Repository)
	err = json.Unmarshal(bytes.TrimSpace([]byte(content)), repo)
	if err != nil {
		return repo, err
	}
	repo.RepoCode = repoCode
	repo.Name = name
	return repo, nil
}

// GetRepository returns the repository details based on Id
func (aspace *ArchivesSpaceAPI) GetRepository(id int) (*Repository, error) {
	aspace.URL.Path = fmt.Sprintf(`/repositories/%d`, id)
	//payload := strings.NewReader(fmt.Sprintf(`{"repo_code":%q,"name":%q}`, repoCode, name))

	client := &http.Client{}
	req, err := http.NewRequest("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can't get repository: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	if res.Status != "200 OK" {
		return nil, fmt.Errorf("ArchiveSpace API error %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	// content should look something like
	// {"lock_version":0,"repo_code":"1447893780","name":"This is a test generated from go_test","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-19T00:43:00Z","system_mtime":"2015-11-19T00:43:00Z","user_mtime":"2015-11-19T00:43:00Z","jsonmodel_type":"repository","uri":"/repositories/16","agent_representation":{"ref":"/agents/corporate_entities/15"}}
	repo := new(Repository)
	repo.ID = id
	err = json.Unmarshal(bytes.TrimSpace([]byte(content)), repo)
	if err != nil {
		return repo, err
	}
	return repo, nil
}

// UpdateRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) UpdateRepository(repo *Repository) error {
	/*
		Example get a repo with curl:
			curl -H "X-ArchivesSpace-Session: $TOKEN" --request GET "http://localhost:8089/repositories/9" | python -m json.tool

		Example output

		{
		    "agent_representation": {
		        "ref": "/agents/corporate_entities/8"
		    },
		    "country": "UM",
		    "create_time": "2015-11-20T18:39:49Z",
		    "created_by": "admin",
		    "image_url": "http://identity.example.org",
		    "jsonmodel_type": "repository",
		    "last_modified_by": "admin",
		    "lock_version": 1,
		    "name": "This is a test generated from go_test",
		    "org_code": "Orangization or Agency Code",
		    "parent_institution_name": "Parent Institution Name",
		    "repo_code": "1448044788",
		    "system_mtime": "2015-11-20T18:56:57Z",
		    "uri": "/repositories/9",
		    "url": "http://example.org",
		    "user_mtime": "2015-11-20T18:56:57Z"
		}

		Example Update the repo with curl (updating coutry and url):
			export RECORD='{"lock_version":1,"repo_code":"1448044788","name":"This is a test generated from go_test","org_code":"Orangization or Agency Code","parent_institution_name":"Parent Institution Name","url":"http://www.example.org","image_url":"http://identity.example.org","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-20T18:39:49Z","system_mtime":"2015-11-20T18:56:57Z","user_mtime":"2015-11-20T18:56:57Z","country":"US","jsonmodel_type":"repository","uri":"/repositories/9","agent_representation":{"ref":"/agents/corporate_entities/8"}}'
			curl -H "X-ArchivesSpace-Session: $TOKEN" -d $RECORD --request PUT "http://localhost:8089/repositories/9"
	*/

	aspace.URL.Path = repo.URI
	jsonSrc, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("Can't JSON encode update %v %s", repo, err)
	}
	payload := strings.NewReader(fmt.Sprintf("%s", jsonSrc))

	client := &http.Client{}
	req, err := http.NewRequest("POST", aspace.URL.String(), payload)
	if err != nil {
		return fmt.Errorf("Can't POST update request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Request error: %s", err)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Read body error: %s", err)
	}
	// content should look something like
	// {"status":"Created","id":3,"lock_version":0,"stale":null,"uri":"/repositories/3","warnings":[]}
	type msgStatus struct {
		Status   string   `json:"status"`
		ID       int      `json:"id"`
		URI      string   `json:"uri,omitempty"`
		Warnings []string `json:"warnings,omitempty"`
	}
	data := new(msgStatus)
	err = json.Unmarshal(bytes.TrimSpace([]byte(content)), data)
	if err != nil {
		return fmt.Errorf("Could not unpack UpdateRepository() response [%s] %s", content, err)
	}
	return nil
}

// DeleteRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) DeleteRepository(repo *Repository) error {
	/*
		Example Listing the repo with curl:
			curl -H "X-ArchivesSpace-Session: $TOKEN" --request GET "http://localhost:8089/repositories" | python -m json.tool

		Example Delete the repo with curl:
			curl -H "X-ArchivesSpace-Session: $TOKEN" -d '{"repo_code": "1448043078"}' --request DELETE "http://localhost:8089/repositories/8" | python -m json.tool
	*/
	aspace.URL.Path = fmt.Sprintf("/repositories/%d", repo.ID)
	payload := strings.NewReader(fmt.Sprintf(`{"repo_code":%q}`, repo.RepoCode))

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", aspace.URL.String(), payload)
	if err != nil {
		return fmt.Errorf("Can't create request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("Request error: %s", err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Read body error: %s", err)
	}
	// content should look something like
	// {"status":"Deleted","id":8}
	type statusMsg struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}
	data := new(statusMsg)
	err = json.Unmarshal(bytes.TrimSpace([]byte(content)), data)
	if err != nil {
		return fmt.Errorf("Cannnot decode DeleteRepository() response %s", err)
	}
	if strings.Compare(data.Status, "Deleted") != 0 {
		return fmt.Errorf("DeleteRepository() unexpected status [%s]", data.Status)
	}
	if data.ID != repo.ID {
		return fmt.Errorf("DeleteRepository() unexpected repository id %d", data.ID)
	}
	return nil
}

// ListRepositories returns a list of repositories available via the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) ListRepositories() ([]Repository, error) {
	aspace.URL.Path = `/repositories`

	client := &http.Client{}
	req, err := http.NewRequest("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can't get repository: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	if res.Status != "200 OK" {
		return nil, fmt.Errorf("ArchiveSpace API error %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	// content should look something like
	// {"lock_version":0,"repo_code":"1447893780","name":"This is a test generated from go_test","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-19T00:43:00Z","system_mtime":"2015-11-19T00:43:00Z","user_mtime":"2015-11-19T00:43:00Z","jsonmodel_type":"repository","uri":"/repositories/16","agent_representation":{"ref":"/agents/corporate_entities/15"}}
	var repos []Repository
	err = json.Unmarshal(bytes.TrimSpace([]byte(content)), &repos)
	if err != nil {
		return nil, err
	}
	// Now I need to populate the repos[?].ID fields
	for i := range repos {
		if id, err := strconv.Atoi(strings.TrimPrefix(repos[i].URI, "/repositories/")); err == nil {
			repos[i].ID = id
		}
	}
	return repos, nil
}

// CreateAgent creates a Agent recod via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) CreateAgent(repo *Repository) (*Agent, error) {
	return nil, fmt.Errorf("CreateAgent() not implemented")
}

// ListAgents return a list of agents available in the respository via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) ListAgents() ([]Agent, error) {
	return nil, fmt.Errorf("ListAgents() not implemented")
}

//FIXME Need to implemenent similar methods on agents and accessions
