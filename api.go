//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2015
// Caltech Library
//
package aspace

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

// MergeEnv - given an environment variable get the value or set a default.
func MergeEnv(environmentVar, defaultValue string) string {
	if s := os.Getenv(environmentVar); s != "" {
		return s
	}
	return defaultValue
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(apiURL, username, password string) *ArchivesSpaceAPI {
	api := new(ArchivesSpaceAPI)
	api.URL, _ = url.Parse(MergeEnv("ASPACE_API_URL", apiURL))
	api.AuthToken = MergeEnv("ASPACE_API_TOKEN", "")
	api.Username = MergeEnv("ASPACE_USERNAME", username)
	api.Password = MergeEnv("ASPACE_PASSWORD", password)
	api.DataSet = MergeEnv("ASPACE_DATASET", "data")
	api.Htdocs = MergeEnv("ASPACE_HTDOCS", "htdocs")
	api.Templates = MergeEnv("ASPACE_TEMPLATES", "templates")
	api.BleveIndex = MergeEnv("ASPACE_BLEVE_INDEX", "index.bleve")
	return api
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

	u := api.URL
	u.Path = fmt.Sprintf("/users/%s/login", api.Username)
	form := url.Values{}
	form.Add("password", api.Password)

	res, err := http.PostForm(u.String(), form)
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
	// Save the token and invalidate the one in our aspace struct.
	token := api.AuthToken
	api.AuthToken = ""
	// Using the copied token try to logout from the service.
	u := api.URL
	u.Path = `/logout`
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
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
			return nil, err
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
		defer res.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("Request error: %s", err)
		}
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Read body error: %s", err)
		}
		return content, nil
	}
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
	return content, nil
}

// CreateAPI is a generalized call to create an object form an interface.
func (api *ArchivesSpaceAPI) CreateAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := api.API("POST", url, obj)
	if err != nil {
		return nil, err
	}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Cannnot decode DeleteAPI() response %s", err)
	}
	return data, nil
}

// ListAPI return a list of IDs from an ArchivesSpace instance for given URL
func (api *ArchivesSpaceAPI) ListAPI(url string) ([]int, error) {
	content, err := api.API("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// [1,2,3,4]
	var ids []int
	err = json.Unmarshal(content, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// CreateRepository will create a respository via the REST API for the
// ArchivesSpace instance defined in the ArchivesSpaceAPI struct.
// It will return the created record.
func (api *ArchivesSpaceAPI) CreateRepository(repo *Repository) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = "/repositories"
	return api.CreateAPI(u.String(), repo)
}

// GetRepository returns the repository details based on Id
func (api *ArchivesSpaceAPI) GetRepository(id int) (*Repository, error) {
	u := *api.URL
	u.Path = fmt.Sprintf(`/repositories/%d`, id)
	repo := new(Repository)
	err := api.GetAPI(u.String(), repo)
	if err != nil {
		return nil, err
	}
	repo.ID = URIToID(repo.URI)
	return repo, nil
}

// UpdateRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) UpdateRepository(repo *Repository) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = repo.URI
	return api.UpdateAPI(u.String(), repo)
}

// DeleteRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) DeleteRepository(repo *Repository) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/repositories/%d", repo.ID)
	return api.DeleteAPI(u.String(), repo)
}

// ListRepositoryIDs returns the numeric ids for all respoistories via the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) ListRepositoryIDs() ([]int, error) {
	var ids []int
	var repos []Repository

	u := *api.URL
	u.Path = `/repositories`
	content, err := api.API("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, err
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
	u := *api.URL
	u.Path = `/repositories`

	content, err := api.API("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"lock_version":0,"repo_code":"1447893780","name":"This is a test generated from go_test","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-19T00:43:00Z","system_mtime":"2015-11-19T00:43:00Z","user_mtime":"2015-11-19T00:43:00Z","jsonmodel_type":"repository","uri":"/repositories/16","agent_representation":{"ref":"/agents/corporate_entities/15"}}
	var repos []Repository
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, err
	}
	// Now I need to populate the repos[?].ID fields
	for i := range repos {
		repos[i].ID = URIToID(repos[i].URI)
	}
	return repos, nil
}

// CreateAgent creates a Agent recod via the ArchivesSpace API
func (api *ArchivesSpaceAPI) CreateAgent(aType string, agent *Agent) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/agents/%s", aType)
	return api.CreateAPI(u.String(), agent)
}

// GetAgent return an Agent via the ArchivesSpace API
func (api *ArchivesSpaceAPI) GetAgent(agentType string, agentID int) (*Agent, error) {
	u := *api.URL
	u.Path = fmt.Sprintf(`/agents/%s/%d`, agentType, agentID)

	agent := new(Agent)
	err := api.GetAPI(u.String(), agent)
	if err != nil {
		return nil, err
	}
	agent.ID = URIToID(agent.URI)
	return agent, nil
}

// UpdateAgent creates a Agent recod via the ArchivesSpace API
func (api *ArchivesSpaceAPI) UpdateAgent(agent *Agent) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = agent.URI
	return api.UpdateAPI(u.String(), agent)
}

// DeleteAgent creates a Agent record via the ArchivesSpace API
func (api *ArchivesSpaceAPI) DeleteAgent(agent *Agent) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = agent.URI
	return api.DeleteAPI(u.String(), agent)
}

// ListAgents return an array of Agents via the ArchivesSpace API
func (api *ArchivesSpaceAPI) ListAgents(agentType string) ([]int, error) {
	u := *api.URL
	u.Path = fmt.Sprintf(`/agents/%s`, agentType)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return api.ListAPI(u.String())
}

// CreateAccession creates a new Accession record in a Repository
func (api *ArchivesSpaceAPI) CreateAccession(repoID int, accession *Accession) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/repositories/%d/accessions", repoID)
	return api.CreateAPI(u.String(), accession)
}

// GetAccession retrieves an Accession record from a Repository
func (api *ArchivesSpaceAPI) GetAccession(repoID, accessionID int) (*Accession, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/repositories/%d/accessions/%d", repoID, accessionID)

	accession := new(Accession)
	err := api.GetAPI(u.String(), accession)
	if err != nil {
		return nil, err
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
	u := *api.URL
	u.Path = accession.URI
	return api.UpdateAPI(u.String(), accession)
}

// DeleteAccession deleted an Accession record from a Repository
func (api *ArchivesSpaceAPI) DeleteAccession(accession *Accession) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = accession.URI
	return api.DeleteAPI(u.String(), accession)
}

// ListAccessions return a list of Accession IDs from a Repository
func (api *ArchivesSpaceAPI) ListAccessions(repositoryID int) ([]int, error) {
	u := *api.URL
	u.Path = fmt.Sprintf(`/repositories/%d/accessions`, repositoryID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return api.ListAPI(u.String())
}

// CreateSubject creates a new Subject in ArchivesSpace instance
func (api *ArchivesSpaceAPI) CreateSubject(subject *Subject) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = "/subjects"
	return api.CreateAPI(u.String(), subject)
}

// GetSubject retrieves a subject record from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) GetSubject(subjectID int) (*Subject, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/subjects/%d", subjectID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"subject","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","vocabulary":"/vocabularies/1"}],"external_documents":[],"uri":"/subjects/1","is_linked_to_published_record":true,"vocabulary":"/vocabularies/1"}
	subject := new(Subject)
	err := api.GetAPI(u.String(), subject)
	p := strings.Split(subject.URI, "/")
	subject.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return subject, fmt.Errorf("Accession ID parse error %d %s", subject.ID, err)
	}
	return subject, nil
}

// UpdateSubject updates an existing subject record in an ArchivesSpace instance
func (api *ArchivesSpaceAPI) UpdateSubject(subject *Subject) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = subject.URI
	return api.UpdateAPI(u.String(), subject)
}

// DeleteSubject deletes a subject from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) DeleteSubject(subject *Subject) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = subject.URI
	return api.DeleteAPI(u.String(), subject)
}

// ListSubjects return a list of Subject IDs from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) ListSubjects() ([]int, error) {
	u := *api.URL
	u.Path = `/subjects`
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return api.ListAPI(u.String())
}

// CreateVocabulary creates a new Vocabulary in ArchivesSpace instance
func (api *ArchivesSpaceAPI) CreateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = "/vocabularies"
	return api.CreateAPI(u.String(), vocabulary)
}

// GetVocabulary retrieves a vocabulary record from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) GetVocabulary(vocabularyID int) (*Vocabulary, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/vocabularies/%d", vocabularyID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"vocabulary","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","vocabulary":"/vocabularies/1"}],"external_documents":[],"uri":"/vocabularys/1","is_linked_to_published_record":true,"vocabulary":"/vocabularies/1"}
	vocabulary := new(Vocabulary)
	err := api.GetAPI(u.String(), vocabulary)
	p := strings.Split(vocabulary.URI, "/")
	vocabulary.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return vocabulary, fmt.Errorf("Accession ID parse error %d %s", vocabulary.ID, err)
	}
	return vocabulary, nil
}

// UpdateVocabulary updates an existing vocabulary record in an ArchivesSpace instance
func (api *ArchivesSpaceAPI) UpdateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = vocabulary.URI
	return api.UpdateAPI(u.String(), vocabulary)
}

// DeleteVocabulary deletes a vocabulary from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) DeleteVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = vocabulary.URI
	return api.DeleteAPI(u.String(), vocabulary)
}

// ListVocabularies return a list of Vocabulary IDs from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) ListVocabularies() ([]int, error) {
	u := *api.URL
	u.Path = `/vocabularies`
	/*
		q := u.Query()
		q.Set("all_ids", "true")
		u.RawQuery = q.Encode()
	*/
	content, err := api.API("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	var (
		ids          []int
		vocabularies []Vocabulary
	)
	err = json.Unmarshal([]byte(content), &vocabularies)
	for _, val := range vocabularies {
		p := strings.Split(val.URI, "/")
		id, err := strconv.Atoi(p[len(p)-1])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateTerm creates a new Term in ArchivesSpace instance
func (api *ArchivesSpaceAPI) CreateTerm(vocabularyID int, term *Term) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/vocabularies/%d/terms", vocabularyID)
	return api.CreateAPI(u.String(), term)
}

// GetTerm retrieves a term record from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) GetTerm(vocabularyID, termID int) (*Term, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/vocabularies/%d/terms", vocabularyID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"term","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","term":"/terms/1"}],"external_documents":[],"uri":"/terms/1","is_linked_to_published_record":true,"term":"/terms/1"}
	terms, err := api.ListTerms(vocabularyID)
	if err != nil {
		return nil, err
	}
	for _, term := range terms {
		term.ID = URIToID(term.URI)
		if term.ID == termID {
			return term, nil
		}
	}
	return nil, nil
}

// UpdateTerm updates an existing term record in an ArchivesSpace instance
func (api *ArchivesSpaceAPI) UpdateTerm(term *Term) (*ResponseMsg, error) {
	u := api.URL
	u.Path = term.URI
	return api.UpdateAPI(u.String(), term)
}

// DeleteTerm deletes a term from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) DeleteTerm(term *Term) (*ResponseMsg, error) {
	u := api.URL
	u.Path = term.URI
	return api.DeleteAPI(u.String(), term)
}

// ListTermIDs return a list of Term IDs from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) ListTermIDs(vocabularyID int) ([]int, error) {
	u := api.URL
	u.Path = fmt.Sprintf(`/vocabularies/%d/terms`, vocabularyID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	data, err := api.API("GET", u.String(), nil)
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

// ListTerms return a list of Term IDs from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) ListTerms(vocabularyID int) ([]*Term, error) {
	u := api.URL
	u.Path = fmt.Sprintf(`/vocabularies/%d/terms`, vocabularyID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	data, err := api.API("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can't get Terms for vocabulary %d, %s", vocabularyID, err)
	}
	// Now Unpack list of terms into a []Term
	var terms []*Term
	err = json.Unmarshal(data, &terms)
	if err != nil {
		return nil, fmt.Errorf("Can't decode terms for vocabularly %d, %s", vocabularyID, err)
	}
	for _, term := range terms {
		//FIXME: Get the Term id and set terms[i].ID to that value.
		p := strings.Split(term.URI, "/")
		term.ID, err = strconv.Atoi(p[len(p)-1])
	}
	return terms, nil
}

// CreateLocation creates a new Location in ArchivesSpace instance
func (api *ArchivesSpaceAPI) CreateLocation(location *Location) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/locations")
	return api.CreateAPI(u.String(), location)
}

// GetLocation retrieves a location record from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) GetLocation(ID int) (*Location, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/locations/%d", ID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"location","external_ids":[],"publish":true,"locations":[{"lock_version":0,"location":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","location_type":"function","jsonmodel_type":"location","uri":"/locations/1","location":"/locations/1"}],"external_documents":[],"uri":"/locations/1","is_linked_to_published_record":true,"location":"/locations/1"}
	location := new(Location)
	err := api.GetAPI(u.String(), location)
	if err != nil {
		return nil, err
	}
	p := strings.Split(location.URI, "/")
	location.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return location, fmt.Errorf("Accession ID parse error %d %s", location.ID, err)
	}
	return location, nil
}

// UpdateLocation updates an existing location record in an ArchivesSpace instance
func (api *ArchivesSpaceAPI) UpdateLocation(location *Location) (*ResponseMsg, error) {
	u := api.URL
	u.Path = location.URI
	return api.UpdateAPI(u.String(), location)
}

// DeleteLocation deletes a location from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) DeleteLocation(location *Location) (*ResponseMsg, error) {
	u := api.URL
	u.Path = location.URI
	return api.DeleteAPI(u.String(), location)
}

// ListLocations return a list of Location IDs from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) ListLocations() ([]int, error) {
	u := api.URL
	u.Path = fmt.Sprintf(`/locations`)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return api.ListAPI(u.String())
}

// CreateDigitalObject - return a new digital object
func (api *ArchivesSpaceAPI) CreateDigitalObject(repoID int, obj *DigitalObject) (*ResponseMsg, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/repositories/%d/digital_object", repoID)
	return api.CreateAPI(u.String(), obj)
}

// GetDigitalObject - return a given digital object
func (api *ArchivesSpaceAPI) GetDigitalObject(repoID, objID int) (*DigitalObject, error) {
	u := *api.URL
	u.Path = fmt.Sprintf("/repositories/%d/digital_objects/%d", repoID, objID)

	obj := new(DigitalObject)
	err := api.GetAPI(u.String(), obj)
	if err != nil {
		return nil, fmt.Errorf("GetDigitalObject() %s, error, %s", u.String(), err)
	}
	obj.ID = URIToID(obj.URI)
	return obj, nil
}

// UpdateDigitalObject - returns an updated digital
func (api *ArchivesSpaceAPI) UpdateDigitalObject(obj *DigitalObject) (*ResponseMsg, error) {
	u := api.URL
	u.Path = obj.URI
	return api.UpdateAPI(u.String(), obj)
}

// DeleteDigitalObject - return the results of deleting a digital object
func (api *ArchivesSpaceAPI) DeleteDigitalObject(obj *DigitalObject) (*ResponseMsg, error) {
	u := api.URL
	u.Path = obj.URI
	return api.DeleteAPI(u.String(), obj)
}

// ListDigitalObjects - return a list of digital object ids
func (api *ArchivesSpaceAPI) ListDigitalObjects(repoID int) ([]int, error) {
	u := api.URL
	u.Path = fmt.Sprintf(`/repositories/%d/digital_objects`, repoID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return api.ListAPI(u.String())
}


// Search return a JSON content from search results from an ArchivesSpace instance
func (api *ArchivesSpaceAPI) Search(opt *SearchQuery) ([]byte, error) {
	u := api.URL
	if opt.URI != "" {
		u.Path = opt.URI
	} else {
		u.Path = "/search"
	}
	q := u.Query()
	//FIXME: Need to walk the struct provided by opt..
	if opt.Q != "" {
		q.Set("q", opt.Q)
	}
	if opt.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", opt.Page))
	} else {
		q.Set("page", "1")
	}
	if opt.PageSize > 0 {
		q.Set("page_size", fmt.Sprintf("%d", opt.PageSize))
	}
	/*
		if opt.Type != "" {
			q.Set("type", opt.Type)
		}
		if opt.Sort != "" {
			q.Set("sort", opt.Sort)
		}
		//FIXME: Need to understand how to express facits in a URL
		if len(opt.Facets) > 0 {

		}
		for k, v := range opt.FilterTerm {
			q.Set(k, v)
		}
		if len(opt.Exclude) > 0{
			q.Set("exclude", IntListToString(opt.Exclude, ","))
		}
		//Skipping RootRecord and RESTHelpers
		if len(opt.IDSet) > 0 {
			q.Set("id_set", IntListToString(opt.IDSet, ","))
		}
		if opt.AllIDs == true {
			q.Set("all_ids", "true")
		}
	*/
	u.RawQuery = q.Encode()
	searchResults := new(SearchResults)
	return api.API("GET", u.String(), &searchResults)
}


//FIXME: need Create, Get, Update, Delete, List functions for Resources, Extents, Instances, Group, Users
