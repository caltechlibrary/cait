//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package cait

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Version of library
var (
	Version = "v0.0.15"

	LicenseText = `
%s %s

Copyright (c) 2017, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

func getenv(envvar, defaultValue string) string {
	tmp := os.Getenv(envvar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(apiURL, username, password string) *ArchivesSpaceAPI {
	api := new(ArchivesSpaceAPI)
	api.BaseURL, _ = url.Parse(getenv("CAIT_API_URL", apiURL))
	api.CallURL, _ = url.Parse(getenv("CAIT_API_URL", apiURL))
	api.AuthToken = getenv("CAIT_API_TOKEN", "")
	api.Username = getenv("CAIT_USERNAME", username)
	api.Password = getenv("CAIT_PASSWORD", password)
	api.Dataset = getenv("CAIT_DATASET", "dataset")
	api.DatasetIndex = getenv("CAIT_DATASET_INDEX", "dataset.bleve")
	api.Htdocs = getenv("CAIT_HTDOCS", "htdocs")
	api.HtdocsIndex = getenv("CAIT_HTDOCS_INDEX", "htdocs.bleve")
	api.Templates = getenv("CAIT_TEMPLATES", "templates/default")
	return api
}

// UpdateCallPath takes the BaseURL Path attribute, copies it into CallURL, applies appends a path for next API call
func (api *ArchivesSpaceAPI) UpdateCallPath(p string) string {
	api.CallURL.Path = api.BaseURL.Path + p
	return api.CallURL.Path
}

// IsAuth returns true if the auth token has been set, false otherwise
func (api *ArchivesSpaceAPI) IsAuth() bool {
	if api.AuthToken == "" {
		return false
	}
	return true
}

// Login authenticates against the ArchivesSpace REST API setting the AuthToken
// value in the ArchivesSpaceAPI struct.
func (api *ArchivesSpaceAPI) Login() error {
	// See https://golang.org/pkg/net/url/#pkg-examples for example building a URL from parts.
	// Command line example: curl -F "password=admin" "http://localhost:8089/users/admin/login"
	var data map[string]interface{}

	// If we already have a token set then logout and get a new one
	if api.IsAuth() == true {
		api.Logout()
	}

	api.UpdateCallPath(fmt.Sprintf("/users/%s/login", api.Username))
	form := url.Values{}
	form.Add("password", api.Password)

	res, err := http.PostForm(api.CallURL.String(), form)
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
	api.AuthToken = data["session"].(string)
	return nil
}

// Logout clear the authentication token for the session with the API
func (api *ArchivesSpaceAPI) Logout() error {
	// Save the token and invalidate the one in our cait struct.
	token := api.AuthToken
	api.AuthToken = ""
	// Using the copied token try to logout from the service.
	api.UpdateCallPath(`/logout`)
	client := &http.Client{}
	req, err := http.NewRequest("GET", api.CallURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-ArchivesSpace-Session", token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// API the common HTTP request processing for interacting with ArchivesSpaceAPI
func (api *ArchivesSpaceAPI) API(method string, url string, data interface{}) ([]byte, error) {
	var (
		payload []byte
		err     error
	)
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("API(%q, %q, data), %s", method, url, err)
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("Can't create request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", api.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	if method == "POST" {
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Request error: %s", err)
		}
		defer res.Body.Close()
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Read body error: %s", err)
		}
		return content, nil
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	defer res.Body.Close()
	if res.Status != "200 OK" {
		return nil, fmt.Errorf("ArchiveSpace API error %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	return content, nil
}

// CreateAPI is a generalized call to create an object form an interface.
func (api *ArchivesSpaceAPI) CreateAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := api.API("POST", url, obj)
	if err != nil {
		return nil, fmt.Errorf("Create API, %s, %s", content, err)
	}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Create API, unmarshal response msg, %s", err)
	}
	return data, nil
}

// GetAPI is a generalized call to get a specific object from an interface
// obj is modified as a side effect
func (api *ArchivesSpaceAPI) GetAPI(url string, obj interface{}) error {
	content, err := api.API("GET", url, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, obj)
	if err != nil {
		return fmt.Errorf("unmarshal error %s, %s\n", content, err)
	}
	return nil
}

// UpdateAPI is a generalized call to update an object from an interface.
func (api *ArchivesSpaceAPI) UpdateAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := api.API("POST", url, obj)
	if err != nil {
		return nil, fmt.Errorf("UpdateAPI(%q, obj) %s", url, err)
	}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Could not unpack UpdateAPI() response [%s] %s", content, err)
	}
	return data, nil
}

// DeleteAPI is a generalized call to update an object form an interface
func (api *ArchivesSpaceAPI) DeleteAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := api.API("DELETE", url, obj)
	if err != nil {
		return nil, fmt.Errorf("DeleteAPI(%q, obj) %s", url, err)
	}

	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Cannnot decode DeleteAPI() response %s", err)
	}
	return data, nil
}

// ListAPI return a list of IDs from ArchivesSpace for given URL
func (api *ArchivesSpaceAPI) ListAPI(url string) ([]int, error) {
	content, err := api.API("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ListAPI(%q) %s", url, err)
	}

	// content should look something like
	// [1,2,3,4]
	var ids []int
	err = json.Unmarshal(content, &ids)
	if err != nil {
		return nil, fmt.Errorf("ListAPI(%q) %s", url, err)
	}
	return ids, nil
}

// CreateRepository will create a repository via the REST API for
// ArchivesSpace defined in the ArchivesSpaceAPI struct.
// It will return the created record.
func (api *ArchivesSpaceAPI) CreateRepository(repo *Repository) (*ResponseMsg, error) {
	api.UpdateCallPath("/repositories")
	return api.CreateAPI(api.CallURL.String(), repo)
}

// GetRepository returns the repository details based on Id
func (api *ArchivesSpaceAPI) GetRepository(id int) (*Repository, error) {
	api.UpdateCallPath(fmt.Sprintf(`/repositories/%d`, id))
	repo := new(Repository)
	err := api.GetAPI(api.CallURL.String(), repo)
	if err != nil {
		return nil, fmt.Errorf("GetRepostiory(%d) %s", id, err)
	}
	repo.ID = URIToID(repo.URI)
	return repo, nil
}

// UpdateRepository takes a repository structure and sends it to the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) UpdateRepository(repo *Repository) (*ResponseMsg, error) {
	api.UpdateCallPath(repo.URI)
	return api.UpdateAPI(api.CallURL.String(), repo)
}

// DeleteRepository takes a repository structure and sends it to the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) DeleteRepository(repo *Repository) (*ResponseMsg, error) {
	api.UpdateCallPath(fmt.Sprintf("/repositories/%d", repo.ID))
	return api.DeleteAPI(api.CallURL.String(), repo)
}

// ListRepositoryIDs returns the numeric ids for all respoistories via the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) ListRepositoryIDs() ([]int, error) {
	var ids []int
	var repos []Repository

	api.UpdateCallPath(`/repositories`)
	content, err := api.API("GET", api.CallURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ListRepositoryIDs() %s", err)
	}
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, fmt.Errorf("ListRepositoryIDs() %s", err)
	}
	// Now I need to populate out id list
	for i := range repos {
		if id, err := strconv.Atoi(strings.TrimPrefix(repos[i].URI, "/repositories/")); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// ListRepositories returns a list of repositories available via the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) ListRepositories() ([]Repository, error) {
	api.UpdateCallPath(`/repositories`)

	content, err := api.API("GET", api.CallURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ListRepositories() %s", err)
	}

	var repos []Repository
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, fmt.Errorf("ListRepositories() %s", err)
	}
	// Now I need to populate the repos[?].ID fields
	for i := range repos {
		repos[i].ID = URIToID(repos[i].URI)
	}
	return repos, nil
}

// CreateAgent creates a Agent recod via the ArchivesSpace API
func (api *ArchivesSpaceAPI) CreateAgent(aType string, agent *Agent) (*ResponseMsg, error) {
	api.UpdateCallPath(fmt.Sprintf("/agents/%s", aType))
	agent.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), agent)
}

// GetAgent return an Agent via the ArchivesSpace API
func (api *ArchivesSpaceAPI) GetAgent(agentType string, agentID int) (*Agent, error) {
	api.UpdateCallPath(fmt.Sprintf(`/agents/%s/%d`, agentType, agentID))

	agent := new(Agent)
	err := api.GetAPI(api.CallURL.String(), agent)
	if err != nil {
		return nil, fmt.Errorf("GetAgent(%s, %d) %s", agentType, agentID, err)
	}
	agent.ID = URIToID(agent.URI)
	return agent, nil
}

// UpdateAgent creates a Agent recod via the ArchivesSpace API
func (api *ArchivesSpaceAPI) UpdateAgent(agent *Agent) (*ResponseMsg, error) {
	api.UpdateCallPath(agent.URI)
	return api.UpdateAPI(api.CallURL.String(), agent)
}

// DeleteAgent creates a Agent record via the ArchivesSpace API
func (api *ArchivesSpaceAPI) DeleteAgent(agent *Agent) (*ResponseMsg, error) {
	api.UpdateCallPath(agent.URI)
	return api.DeleteAPI(api.CallURL.String(), agent)
}

// ListAgents return an array of Agents via the ArchivesSpace API
func (api *ArchivesSpaceAPI) ListAgents(agentType string) ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/agents/%s`, agentType))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}

// CreateAccession creates a new Accession record in a Repository
func (api *ArchivesSpaceAPI) CreateAccession(repoID int, accession *Accession) (*ResponseMsg, error) {
	api.UpdateCallPath(fmt.Sprintf("/repositories/%d/accessions", repoID))
	accession.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), accession)
}

// GetAccession retrieves an Accession record from a Repository
func (api *ArchivesSpaceAPI) GetAccession(repoID, accessionID int) (*Accession, error) {
	api.UpdateCallPath(fmt.Sprintf("/repositories/%d/accessions/%d", repoID, accessionID))

	accession := new(Accession)
	err := api.GetAPI(api.CallURL.String(), accession)
	if err != nil {
		return nil, fmt.Errorf("GetAccession(%d, %d) %s", repoID, accessionID, err)
	}
	p := strings.Split(accession.URI, "/")
	accession.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return accession, fmt.Errorf("Accession ID parse error %d %s", accession.ID, err)
	}
	return accession, nil
}

// UpdateAccession updates an existing Accession record in a Repository
func (api *ArchivesSpaceAPI) UpdateAccession(accession *Accession) (*ResponseMsg, error) {
	api.UpdateCallPath(accession.URI)
	return api.UpdateAPI(api.CallURL.String(), accession)
}

// DeleteAccession deleted an Accession record from a Repository
func (api *ArchivesSpaceAPI) DeleteAccession(accession *Accession) (*ResponseMsg, error) {
	api.UpdateCallPath(accession.URI)
	return api.DeleteAPI(api.CallURL.String(), accession)
}

// ListAccessions return a list of Accession IDs from a Repository
func (api *ArchivesSpaceAPI) ListAccessions(repositoryID int) ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/repositories/%d/accessions`, repositoryID))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}

// CreateSubject creates a new Subject in ArchivesSpace
func (api *ArchivesSpaceAPI) CreateSubject(subject *Subject) (*ResponseMsg, error) {
	api.UpdateCallPath("/subjects")
	subject.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), subject)
}

// GetSubject retrieves a subject record from ArchivesSpace
func (api *ArchivesSpaceAPI) GetSubject(subjectID int) (*Subject, error) {
	api.UpdateCallPath(fmt.Sprintf("/subjects/%d", subjectID))

	subject := new(Subject)
	err := api.GetAPI(api.CallURL.String(), subject)
	p := strings.Split(subject.URI, "/")
	subject.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return subject, fmt.Errorf("Accession ID parse error %d %s", subject.ID, err)
	}
	return subject, nil
}

// UpdateSubject updates an existing subject record in ArchivesSpace
func (api *ArchivesSpaceAPI) UpdateSubject(subject *Subject) (*ResponseMsg, error) {
	api.UpdateCallPath(subject.URI)
	return api.UpdateAPI(api.CallURL.String(), subject)
}

// DeleteSubject deletes a subject from ArchivesSpace
func (api *ArchivesSpaceAPI) DeleteSubject(subject *Subject) (*ResponseMsg, error) {
	api.UpdateCallPath(subject.URI)
	return api.DeleteAPI(api.CallURL.String(), subject)
}

// ListSubjects return a list of Subject IDs from ArchivesSpace
func (api *ArchivesSpaceAPI) ListSubjects() ([]int, error) {
	api.UpdateCallPath(`/subjects`)
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}

// CreateVocabulary creates a new Vocabulary in ArchivesSpace
func (api *ArchivesSpaceAPI) CreateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	api.UpdateCallPath("/vocabularies")
	vocabulary.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), vocabulary)
}

// GetVocabulary retrieves a vocabulary record from ArchivesSpace
func (api *ArchivesSpaceAPI) GetVocabulary(vocabularyID int) (*Vocabulary, error) {
	api.UpdateCallPath(fmt.Sprintf("/vocabularies/%d", vocabularyID))

	vocabulary := new(Vocabulary)
	err := api.GetAPI(api.CallURL.String(), vocabulary)
	p := strings.Split(vocabulary.URI, "/")
	vocabulary.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return vocabulary, fmt.Errorf("Accession ID parse error %d %s", vocabulary.ID, err)
	}
	return vocabulary, nil
}

// UpdateVocabulary updates an existing vocabulary record in ArchivesSpace
func (api *ArchivesSpaceAPI) UpdateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	api.UpdateCallPath(vocabulary.URI)
	return api.UpdateAPI(api.CallURL.String(), vocabulary)
}

// DeleteVocabulary deletes a vocabulary from ArchivesSpace
func (api *ArchivesSpaceAPI) DeleteVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	api.UpdateCallPath(vocabulary.URI)
	return api.DeleteAPI(api.CallURL.String(), vocabulary)
}

// ListVocabularies return a list of Vocabulary IDs from ArchivesSpace
func (api *ArchivesSpaceAPI) ListVocabularies() ([]int, error) {
	api.UpdateCallPath(`/vocabularies`)
	content, err := api.API("GET", api.CallURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ListVocabularies() %s", err)
	}
	var (
		ids          []int
		vocabularies []Vocabulary
	)
	err = json.Unmarshal([]byte(content), &vocabularies)
	if err != nil {
		return nil, fmt.Errorf("ListVocabularies() %s", err)
	}
	for _, val := range vocabularies {
		p := strings.Split(val.URI, "/")
		id, err := strconv.Atoi(p[len(p)-1])
		if err != nil {
			return nil, fmt.Errorf("ListVocabularies() %s", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateTerm creates a new Term in ArchivesSpace
func (api *ArchivesSpaceAPI) CreateTerm(vocabularyID int, term *Term) (*ResponseMsg, error) {
	api.UpdateCallPath(fmt.Sprintf("/vocabularies/%d/terms", vocabularyID))
	term.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), term)
}

// GetTerm retrieves a term record from ArchivesSpace
func (api *ArchivesSpaceAPI) GetTerm(vocabularyID, termID int) (*Term, error) {
	api.UpdateCallPath(fmt.Sprintf("/vocabularies/%d/terms", vocabularyID))

	terms, err := api.ListTerms(vocabularyID)
	if err != nil {
		return nil, fmt.Errorf("GetTerm(%d, %d) %s", vocabularyID, termID, err)
	}
	for _, term := range terms {
		term.ID = URIToID(term.URI)
		if term.ID == termID {
			return term, nil
		}
	}
	return nil, nil
}

// UpdateTerm updates an existing term record in ArchivesSpace
func (api *ArchivesSpaceAPI) UpdateTerm(term *Term) (*ResponseMsg, error) {
	api.UpdateCallPath(term.URI)
	return api.UpdateAPI(api.CallURL.String(), term)
}

// DeleteTerm deletes a term from ArchivesSpace
func (api *ArchivesSpaceAPI) DeleteTerm(term *Term) (*ResponseMsg, error) {
	api.UpdateCallPath(term.URI)
	return api.DeleteAPI(api.CallURL.String(), term)
}

// ListTermIDs return a list of Term IDs from ArchivesSpace
func (api *ArchivesSpaceAPI) ListTermIDs(vocabularyID int) ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/vocabularies/%d/terms`, vocabularyID))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	data, err := api.API("GET", api.CallURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can't get Terms for vocabulary %d, %s", vocabularyID, err)
	}
	// Now Unpack list of terms into a []Term
	var terms []*Term
	err = json.Unmarshal(data, &terms)
	if err != nil {
		return nil, fmt.Errorf("Can't decode terms for vocabularly %d, %s", vocabularyID, err)
	}
	var ids []int
	for _, term := range terms {
		//FIXME: Get the Term id and set terms[i].ID to that value.
		p := strings.Split(term.URI, "/")
		id, err := strconv.Atoi(p[len(p)-1])
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// ListTerms return a list of Term IDs from ArchivesSpace
func (api *ArchivesSpaceAPI) ListTerms(vocabularyID int) ([]*Term, error) {
	api.UpdateCallPath(fmt.Sprintf(`/vocabularies/%d/terms`, vocabularyID))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	data, err := api.API("GET", api.CallURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can't get Terms for vocabulary %d, %s", vocabularyID, err)
	}
	// Now Unpack list of terms into a []Term
	var terms []*Term
	if err := json.Unmarshal(data, &terms); err != nil {
		return nil, fmt.Errorf("Can't decode terms for vocabularly %d, %s", vocabularyID, err)
	}
	for _, term := range terms {
		//FIXME: Get the Term id and set terms[i].ID to that value.
		p := strings.Split(term.URI, "/")
		term.ID, _ = strconv.Atoi(p[len(p)-1])
	}
	return terms, nil
}

// CreateLocation creates a new Location in ArchivesSpace
func (api *ArchivesSpaceAPI) CreateLocation(location *Location) (*ResponseMsg, error) {
	api.UpdateCallPath(fmt.Sprintf("/locations"))
	location.LockVersion = "0"
	return api.CreateAPI(api.CallURL.String(), location)
}

// GetLocation retrieves a location record from ArchivesSpace
func (api *ArchivesSpaceAPI) GetLocation(ID int) (*Location, error) {
	api.UpdateCallPath(fmt.Sprintf("/locations/%d", ID))

	location := new(Location)
	err := api.GetAPI(api.CallURL.String(), location)
	if err != nil {
		return nil, fmt.Errorf("GetLocation(%d) %s", ID, err)
	}
	p := strings.Split(location.URI, "/")
	location.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return location, fmt.Errorf("Accession ID parse error %d %s", location.ID, err)
	}
	return location, nil
}

// UpdateLocation updates an existing location record in ArchivesSpace
func (api *ArchivesSpaceAPI) UpdateLocation(location *Location) (*ResponseMsg, error) {
	api.UpdateCallPath(location.URI)
	return api.UpdateAPI(api.CallURL.String(), location)
}

// DeleteLocation deletes a location from ArchivesSpace
func (api *ArchivesSpaceAPI) DeleteLocation(location *Location) (*ResponseMsg, error) {
	api.UpdateCallPath(location.URI)
	return api.DeleteAPI(api.CallURL.String(), location)
}

// ListLocations return a list of Location IDs from ArchivesSpace
func (api *ArchivesSpaceAPI) ListLocations() ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/locations`))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}

// CreateDigitalObject - return a new digital object
func (api *ArchivesSpaceAPI) CreateDigitalObject(repoID int, obj *DigitalObject) (*ResponseMsg, error) {
	// NOTE: attempt extract accession ID for the edge of importing a digital object as opposed to a clean create
	uriPrefix := fmt.Sprintf("/repositories/%d/digital_objects", repoID)
	obj.JSONModelType = "digital_object"
	obj.LockVersion = "0"
	api.UpdateCallPath(uriPrefix)
	// We need to create the object
	responseMsg, responseErr := api.CreateAPI(api.CallURL.String(), obj)
	if responseErr != nil || responseMsg.Status != "created" {
		return responseMsg, responseErr
	}
	// NOTE: In the case we're importing a digital_object from another ArchivesSpace instance.
	// We need to correct the URI assignment, lock version info and attach to the accession of necessary
	obj.URI = responseMsg.URI
	obj.LockVersion = responseMsg.LockVersion
	return responseMsg, responseErr
}

// GetDigitalObject - return a given digital object
func (api *ArchivesSpaceAPI) GetDigitalObject(repoID, objID int) (*DigitalObject, error) {
	api.UpdateCallPath(fmt.Sprintf("/repositories/%d/digital_objects/%d", repoID, objID))

	obj := new(DigitalObject)
	err := api.GetAPI(api.CallURL.String(), obj)
	if err != nil {
		return nil, fmt.Errorf("GetDigitalObject() %s, error, %s", api.CallURL.String(), err)
	}
	obj.ID = URIToID(obj.URI)
	return obj, nil
}

// UpdateDigitalObject - returns an updated digital
func (api *ArchivesSpaceAPI) UpdateDigitalObject(obj *DigitalObject) (*ResponseMsg, error) {
	api.UpdateCallPath(obj.URI)
	return api.UpdateAPI(api.CallURL.String(), obj)
}

// DeleteDigitalObject - return the results of deleting a digital object
func (api *ArchivesSpaceAPI) DeleteDigitalObject(obj *DigitalObject) (*ResponseMsg, error) {
	api.UpdateCallPath(obj.URI)
	//FIXME: If we're Updating we may need to unlink existing accessions
	return api.DeleteAPI(api.CallURL.String(), obj)
}

// ListDigitalObjects - return a list of digital object ids
func (api *ArchivesSpaceAPI) ListDigitalObjects(repoID int) ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/repositories/%d/digital_objects`, repoID))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}

// CreateResource - return a new resource
func (api *ArchivesSpaceAPI) CreateResource(repoID int, obj *Resource) (*ResponseMsg, error) {
	// NOTE: attempt extract accession ID for the edge of importing a digital object as opposed to a clean create
	uriPrefix := fmt.Sprintf("/repositories/%d/digital_objects", repoID)
	obj.JSONModelType = "digital_object"
	obj.LockVersion = "0"
	api.UpdateCallPath(uriPrefix)
	// We need to create the object
	responseMsg, responseErr := api.CreateAPI(api.CallURL.String(), obj)
	if responseErr != nil || responseMsg.Status != "created" {
		return responseMsg, responseErr
	}
	// NOTE: In the case we're importing a digital_object from another ArchivesSpace instance.
	// We need to correct the URI assignment, lock version info and attach to the accession of necessary
	obj.URI = responseMsg.URI
	obj.LockVersion = responseMsg.LockVersion
	return responseMsg, responseErr
}

// GetResource - return a given resource
func (api *ArchivesSpaceAPI) GetResource(repoID, objID int) (*Resource, error) {
	api.UpdateCallPath(fmt.Sprintf("/repositories/%d/resources/%d", repoID, objID))

	obj := new(Resource)
	err := api.GetAPI(api.CallURL.String(), obj)
	if err != nil {
		return nil, fmt.Errorf("GetResource() %s, error, %s", api.CallURL.String(), err)
	}
	//obj.ID = URIToID(obj.URI)
	return obj, nil
}

// UpdateResource - returns an updated resource
func (api *ArchivesSpaceAPI) UpdateResource(obj *Resource) (*ResponseMsg, error) {
	api.UpdateCallPath(obj.URI)
	return api.UpdateAPI(api.CallURL.String(), obj)
}

// DeleteResource - return the results of deleting a resource
func (api *ArchivesSpaceAPI) DeleteResource(obj *Resource) (*ResponseMsg, error) {
	api.UpdateCallPath(obj.URI)
	return api.DeleteAPI(api.CallURL.String(), obj)
}

// ListResources - return a list of resource ids
func (api *ArchivesSpaceAPI) ListResources(repoID int) ([]int, error) {
	api.UpdateCallPath(fmt.Sprintf(`/repositories/%d/resources`, repoID))
	q := api.CallURL.Query()
	q.Set("all_ids", "true")
	api.CallURL.RawQuery = q.Encode()
	return api.ListAPI(api.CallURL.String())
}
