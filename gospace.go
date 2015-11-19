/**
 * gospace - Golang ArchicesSpace package. A collection of structures and functions
 * for interacting with ArchivesSpace's REST API
 *
 * @author R. S. Doiel, <rsdoiel@caltech.edu>
 * copyright (c) 2015
 * Caltech Library
 */
package gospace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

// Repository represents an ArchivesSpace repository from the client point of view
type Repository struct {
	ID             int       `json:"id"`
	RepoCode       string    `json:"repo_code"`
	Name           string    `json:"name"`
	LockVersion    int       `json:"lock_version,omitempty"`
	CreatedBy      string    `json:"created_by,omitempty"`
	LastModifiedBy string    `json:"last_modified_by,omitempty"`
	CreateTime     time.Time `json:"create_time,omitempty"`
	SystemTime     time.Time `json:"system_time,omitempty"`
	UserMTime      time.Time `json:"user_mtime,omitempty"`
	URI            string    `json:"uri,omitempty"`
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(protocol, host, port, username, password string) *ArchivesSpaceAPI {
	aspace := new(ArchivesSpaceAPI)

	u, err := url.Parse(fmt.Sprintf("%s://%s:%s", protocol, host, port))
	if err != nil {
		log.Fatal(err)
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
	var data map[string]interface{}

	aspace.URL.Path = fmt.Sprintf("/users/%s/login", aspace.Username)
	form := url.Values{}
	form.Add("password", aspace.Password)

	res, err := http.PostForm(aspace.URL.String(), form)
	defer res.Body.Close()
	if err != nil {
		return err
	}
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
	aspace.URL.Path = fmt.Sprintf(`/users/%s/logout`, aspace.Username)
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
