/**
 * aspace_test.go - Test routines for aspace.go
 */
package aspace

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// Get the environment variables needed for testing.
var (
	aspaceURL      = os.Getenv("ASPACE_API_URL")
	aspaceUsername = os.Getenv("ASPACE_USERNAME")
	aspacePassword = os.Getenv("ASPACE_PASSWORD")
)

func checkConfig(t *testing.T) bool {
	isSetup := true
	if aspaceURL == "" {
		t.Error("ASPACE_API_URL environment variable not set.", aspaceURL)
		isSetup = false
	}
	if aspaceUsername == "" {
		t.Error("ASPACE_USERNAME environment variable not set.", aspaceUsername)
		isSetup = false
	}
	if aspacePassword == "" {
		t.Error("ASPACE_PASSWORD environment variable not set.", aspacePassword)
		isSetup = false
	}
	return isSetup
}

func TestSetup(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		log.Fatalf("Environment variables needed to run tests not configured")
	}
	// Make sure we're not talking to a named system (should be localhost)
	u, err := url.Parse(aspaceURL)
	if err != nil {
		log.Fatalf("aspaceURL value doesn't make sense %s %s", aspaceURL, err)
	}
	if strings.Contains(u.Host, "localhost:") == false {
		log.Fatalf("Tests expect to run on http://localhost:8089 not %s", aspaceURL)
	}
	log.Printf("Test setup completed\n")
}

func TestArchiveSpaceAPI(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.SkipNow()
	}

	aspace := New(aspaceURL, aspaceUsername, aspacePassword)
	if aspace.URL == nil {
		t.Errorf("%s\t%s", aspace.URL.String(), aspaceURL)
	}
	if strings.Compare(aspace.URL.String(), fmt.Sprintf("%s", aspaceURL)) != 0 {
		t.Errorf("%s != %s\n", aspace.URL.String(), aspaceURL)
	}

	if aspace.IsAuth() == true {
		t.Error("aspace.IsAuth() returning true before authentication")
	}
	err := aspace.Login()
	if err != nil {
		t.Errorf("%s\t%s", err, aspace.URL.String())
		t.FailNow()
	}
	if aspace.IsAuth() == false {
		t.Error("aspace.IsAuth() return false after authentication")
	}

	err = aspace.Logout()
	if err != nil {
		t.Errorf("Logout() %s", err)
	}
}

func TestRepository(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	aspace := New(aspaceURL, aspaceUsername, aspacePassword)
	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())

	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	repo1 := new(Repository)
	repo1.RepoCode = repoCode
	repo1.Name = "This is a test generated from go_test"
	response, err := aspace.CreateRepository(repo1)
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if response == nil {
		t.Errorf("CeateRepository() should not have a nil response")
	}
	if response.Status != "Created" {
		t.Errorf("CreatedRepository() should return a Created response %s", response)
	}
	repo1.ID = response.ID

	repo2, err := aspace.GetRepository(repo1.ID)
	if err != nil {
		t.Errorf("GetRepository() error: %s", err)
	}
	if repo1.ID != repo2.ID {
		t.Errorf("GetRepository() returned different IDs: %d != %d", repo1.ID, repo2.ID)
	}
	if strings.Compare(repo1.RepoCode, repo2.RepoCode) != 0 {
		t.Errorf("GetRepository() returned different RepoCode: %s != %s\n", repo1.RepoCode, repo2.RepoCode)
	}
	if strings.Compare(repo1.Name, repo2.Name) != 0 {
		t.Errorf("GetRepository() returned different RepoCode: %s != %s\n", repo1.Name, repo2.Name)
	}

	repo2.Name = fmt.Sprintf("Modified Name: %s", repo2.Name)
	repo2.URL = `http://www.archive.example.edu`
	repo2.ImageURL = `http://www.archive.example.edu/logo.svg`
	response, err = aspace.UpdateRepository(repo2)
	if err != nil {
		t.Errorf("UpdateRepository failed for %v: %s", repo2, err)
	}
	if response.Status != "Updated" {
		t.Errorf("UpdateRepository() should return a response.Status of Updated %s", response)
	}
	isOK := true
	repo1, err = aspace.GetRepository(repo2.ID)
	if err != nil {
		t.Errorf("GetRepository() %d after update failed %s", repo2.ID, err)
		isOK = false
	}
	if strings.Compare(repo2.Name, repo1.Name) != 0 {
		t.Errorf("Name [%s] != [%s]", repo1.Name, repo2.Name)
		isOK = false
	}
	if strings.Compare(repo2.URL, repo1.URL) != 0 {
		t.Errorf("URL [%s] != [%s]", repo1.Name, repo2.Name)
		isOK = false
	}
	if strings.Compare(repo2.ImageURL, repo1.ImageURL) != 0 {
		t.Errorf("ImageURL [%s] != [%s]", repo1.Name, repo2.Name)
		isOK = false
	}
	if isOK == false {
		t.Logf("Auth Token: %s", aspace.AuthToken)
		t.FailNow()
	}

	repos, err := aspace.ListRepositories()
	if err != nil {
		t.Errorf("ListRepostiories failed for %v : %s", aspace, err)
	} else if len(repos) == 0 {
		t.Errorf("Expected one or more in repository list: %v", repos)
	}

	_, err = aspace.DeleteRepository(repo2)
	if err != nil {
		t.Errorf("DeleteRepository failed for %v: %s", repo2, err)
		t.FailNow()
	}

	_, err = aspace.GetRepository(repo1.ID)
	if err == nil {
		t.Errorf("GetRepository() should return an error after a deleting repo id %d: %s", repo1.ID, err)
		t.FailNow()
	}
}

func TestAgent(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	aspace := New(aspaceURL, aspaceUsername, aspacePassword)
	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Test the listing of agents by Type or individually by type/id
	for _, aType := range []string{"people", "families", "corporate_entities", "software"} {
		if agentIDs, err := aspace.ListAgents(aType); err != nil {
			t.Errorf(`ListAgents("%s") error: %s`, aType, err)
		} else if len(agentIDs) > 0 {
			for _, id := range agentIDs {
				if agentInfo, err := aspace.GetAgent(aType, id); err != nil {
					t.Errorf(`GetAgent("%s", %d) error: %s`, aType, id, err)
				} else {
					if agentInfo.ID != id {
						t.Errorf("Returned Agent info id does not match requested %d, returened record %d", id, agentInfo.ID)
					}
					uri := fmt.Sprintf("/agents/%s/%d", aType, id)
					if agentInfo.URI != uri {
						t.Errorf("Returned Agent Info URI does not match %s != %s", uri, agentInfo.URI)
					}
					//FIXME: should add more tests for additional fields.
				}
			}
		}
	}

	name0 := new(NamePerson)
	name0.PrimaryName = "Topper"
	name0.RestOfName = "Cosmo"
	name0.NameOrder = "direct"
	name0.SortName = "Topper, Cosmo"
	name0.Source = "local"
	agent0 := new(Agent)
	agent0.Names = append(agent0.Names, name0)

	aType := "people"
	response, err := aspace.CreateAgent(aType, agent0)
	if err != nil {
		t.Errorf(`CreateAgent("%s", %s) error: %s`, aType, agent0, err)
		t.FailNow()
	}
	if response.Status != "Created" {
		t.Errorf(`CreateAgent("%s", %s) status error: %s`, aType, agent0, response)
		t.FailNow()
	}
	agent0.ID = response.ID

	agent1, err := aspace.GetAgent(aType, agent0.ID)
	if err != nil {
		t.Errorf(`GetAgent(%d) failed %s`, agent0.ID, err)
		t.FailNow()
	}

	if agent1.Names[0].PrimaryName != name0.PrimaryName {
		t.Errorf(`CreateAgent("%s", %s), error: Names[0] does not match %s != %s `, aType, agent0, agent0.Names[0].PrimaryName, agent1.Names[0].PrimaryName)
		t.FailNow()
	}
	agent1.Names[0].NameOrder = "inverted"
	response, err = aspace.UpdateAgent(agent1)
	if err != nil {
		t.Errorf(`UpdateAgent(%s), error: %s`, agent1, err)
		t.FailNow()
	}
	if response.Status != "Updated" {
		t.Errorf(`UpdateAgent(%s), status error: %s`, agent1, response)
		t.FailNow()
	}
	agent2, _ := aspace.GetAgent(aType, agent1.ID)
	if strings.Compare(agent2.Names[0].NameOrder, "inverted") != 0 {
		t.Errorf("UpdateAgent(%s), error: Failed to update Names[0].NameOrder [%s] != [%s]", agent1, agent1.Names[0].NameOrder, agent2.Names[0].NameOrder)
		t.FailNow()
	}
	response, err = aspace.DeleteAgent(agent2)
	if err != nil {
		t.Errorf("DaleteAgent(%s), error: %s", agent2, err)
		t.FailNow()
	}
	if response.Status != "Deleted" {
		t.Errorf("DeleteAgent(%s), error: unexpected response status: %s", agent2, response)
		t.FailNow()
	}

	//Do we need to specific test for families, corporate_entities, and software?
	//if so....
	//FIXME: Create an Agent/families
	//FIXME: Update an Agent/families
	//FIXME: Delete an Agent/families
	//FIXME: Create an Agent/corporate_entities
	//FIXME: Update an Agent/corporate_entities
	//FIXME: Delete an Agent/corporate_entities
	//FIXME: Create an Agent/software
	//FIXME: Update an Agent/software
	//FIXME: Delete an Agent/software
}

func TestAccession(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	aspace := New(aspaceURL, aspaceUsername, aspacePassword)
	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())
	repo := new(Repository)
	repo.RepoCode = repoCode
	repo.Name = fmt.Sprintf("This is a test data generated from go_test for working with Accession data in test Repository %s", repoCode)
	response, err := aspace.CreateRepository(repo)
	defer aspace.DeleteRepository(repo)
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if response.Status != "Created" {
		t.Errorf("Erro from CreateRepository() %s", response)
	}
	repo.ID = response.ID
	repo, err = aspace.GetRepository(repo.ID)
	if repo == nil {
		t.Errorf("Repository should not be nil")
	}
	if repo.ID == 0 {
		t.Errorf("Failed to create a test repository for accession testing")
		t.FailNow()
	}

	// Test the listing of accessions
	accessionIDs, err := aspace.ListAccessions(repo.ID)
	if err != nil {
		t.Errorf(`ListAccessions() error: %s`, err)
		t.FailNow()
	}
	if len(accessionIDs) != 0 {
		t.Errorf(`ListAccessions() should return zero accessions for test repository %d, found %d`, repo.ID, len(accessionIDs))
	}

	//FIXME: Need to insert some accessions to test with.
	for i := 1; i <= 10; i++ {
		// This is an minimal Accession record.
		accession1 := new(Accession)
		accession1.ID0 = fmt.Sprintf("%04d", tm.Year())
		accession1.ID1 = fmt.Sprintf("%04d", i)
		accession1.AccessionDate = fmt.Sprintf("%d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
		response, err = aspace.CreateAccession(repo.ID, accession1)
		if err != nil {
			t.Errorf("Can't create accession %v, %s", accession1, err)
			t.FailNow()
		}
		if response.Status != "Created" {
			t.Errorf(`CreateAccession(%s) return unexpected status %s`, accession1, response)
			t.FailNow()
		}
		accession1.ID = response.ID
		accession1.URI = response.URI
		accession2, err := aspace.GetAccession(repo.ID, accession1.ID)
		if err != nil {
			t.Errorf("GetAccession(%d, %d) error %s", repo.ID, accession1.ID, err)
		}
		if accession1.ID0 != accession2.ID0 {
			t.Errorf("Accesion ID0 should be %s, found %s", accession1.ID0, accession2.ID0)
		}
		if accession1.ID1 != accession2.ID1 {
			t.Errorf("Accesion ID1 should be %s, found %s", accession1.ID1, accession2.ID1)
		}
		if accession1.AccessionDate != accession2.AccessionDate {
			t.Errorf("AccessionDate should be %s, found %s", accession1.AccessionDate, accession2.AccessionDate)
		}
	}

	accessionIDs, err = aspace.ListAccessions(repo.ID)
	if err != nil {
		t.Errorf(`ListAccessions() error: %s`, err)
		t.FailNow()
	}
	if len(accessionIDs) != 10 {
		t.Errorf(`ListAccessions() should return tenaccessions for test repository %d, found %d`, repo.ID, len(accessionIDs))
		t.FailNow()
	}

	for _, id := range accessionIDs {
		accessionInfo, err := aspace.GetAccession(repo.ID, id)
		if err != nil {
			t.Errorf(`GetAccession(%d) error: %s`, id, err)
		}
		if accessionInfo.ID != id {
			t.Errorf("Returned Agent info id does not match requested %d, returned record %d", id, accessionInfo.ID)
		}
		uri := fmt.Sprintf(`/repositories/%d/accessions/%d`, repo.ID, id)
		if accessionInfo.URI != uri {
			t.Errorf("Returned Agent Info URI does not match %s != %s", uri, accessionInfo.URI)
		}
		//FIXME: Need tests for UpdateAccession() and DeleteAccession()
		//FIXME: should add more tests for additional fields.
	}
}

func TestSubjects(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	aspace := New(aspaceURL, aspaceUsername, aspacePassword)
	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	term := new(Term)
	term.Term = "Hello World"
	term.TermType = "topical"
	term.Vocabulary = "/vocabularies/1"
	subject := new(Subject)
	subject.Source = "local"
	subject.Terms = append(subject.Terms, term)
	subject.Vocabulary = "/vocabularies/1"
	response, err := aspace.CreateSubject(subject)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("DEBUG aspace.CreateSubject() --> %s\n", response)
}
