//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
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
	caitURL      = os.Getenv("CAIT_API_URL")
	caitToken    = os.Getenv("CAIT_API_TOKEN")
	caitUsername = os.Getenv("CAIT_USERNAME")
	caitPassword = os.Getenv("CAIT_PASSWORD")
)

func checkConfig(t *testing.T) bool {
	isSetup := true
	if caitURL == "" {
		t.Error("CAIT_API_URL environment variable not set.", caitURL)
		isSetup = false
	}
	if caitToken != "" {
		t.Error("CAIT_API_TOKEN already set, should be empty for tests.", caitToken)
		isSetup = false
	}
	if caitUsername == "" {
		t.Error("CAIT_USERNAME environment variable not set.", caitUsername)
		isSetup = false
	}
	if caitPassword == "" {
		t.Error("CAIT_PASSWORD environment variable not set.", caitPassword)
		isSetup = false
	}
	return isSetup
}

func TestSetup(t *testing.T) {
	log.Printf("Checking configuration before starting tests.\n")
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		log.Println("Environment variables needed to run tests not configured")
		t.FailNow()
	}
	// Make sure we're not talking to a named system (should be localhost)
	u, err := url.Parse(caitURL)
	if err != nil {
		log.Printf("caitURL value doesn't make sense %s %s", caitURL, err)
		t.FailNow()
	}
	if strings.Contains(u.Host, "localhost:") == false {
		log.Printf("Tests expect to run on http://localhost:8089 not %s", caitURL)
		t.FailNow()
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

	cait := New(caitURL, caitUsername, caitPassword)
	if cait.URL == nil {
		t.Errorf("%s\t%s", cait.URL.String(), caitURL)
	}
	if strings.Compare(cait.URL.String(), fmt.Sprintf("%s", caitURL)) != 0 {
		t.Errorf("%s != %s\n", cait.URL.String(), caitURL)
	}

	if cait.IsAuth() == true {
		t.Error("cait.IsAuth() returning true before authentication")
	}
	err := cait.Login()
	if err != nil {
		t.Errorf("%s\t%s", err, cait.URL.String())
		t.FailNow()
	}
	if cait.IsAuth() == false {
		t.Error("cait.IsAuth() return false after authentication")
	}

	err = cait.Logout()
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

	cait := New(caitURL, caitUsername, caitPassword)
	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())

	err := cait.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	repo1 := new(Repository)
	repo1.RepoCode = repoCode
	repo1.Name = "This is a test generated from go_test"
	response, err := cait.CreateRepository(repo1)
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

	repo2, err := cait.GetRepository(repo1.ID)
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
	response, err = cait.UpdateRepository(repo2)
	if err != nil {
		t.Errorf("UpdateRepository failed for %v: %s", repo2, err)
	}
	if response.Status != "Updated" {
		t.Errorf("UpdateRepository() should return a response.Status of Updated %s", response)
	}
	isOK := true
	repo1, err = cait.GetRepository(repo2.ID)
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
		t.Logf("Auth Token: %s", cait.AuthToken)
		t.FailNow()
	}

	repos, err := cait.ListRepositories()
	if err != nil {
		t.Errorf("ListRepostiories failed for %v : %s", cait, err)
	} else if len(repos) == 0 {
		t.Errorf("Expected one or more in repository list: %v", repos)
	}

	_, err = cait.DeleteRepository(repo2)
	if err != nil {
		t.Errorf("DeleteRepository failed for %v: %s", repo2, err)
		t.FailNow()
	}

	_, err = cait.GetRepository(repo1.ID)
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

	cait := New(caitURL, caitUsername, caitPassword)
	err := cait.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Test the listing of agents by Type or individually by type/id
	for _, aType := range []string{"people", "families", "corporate_entities", "software"} {
		if agentIDs, err := cait.ListAgents(aType); err != nil {
			t.Errorf(`ListAgents("%s") error: %s`, aType, err)
		} else if len(agentIDs) > 0 {
			for _, id := range agentIDs {
				if agentInfo, err := cait.GetAgent(aType, id); err != nil {
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
	response, err := cait.CreateAgent(aType, agent0)
	if err != nil {
		t.Errorf(`CreateAgent("%s", %s) error: %s`, aType, agent0, err)
		t.FailNow()
	}
	if response.Status != "Created" {
		t.Errorf(`CreateAgent("%s", %s) status error: %s`, aType, agent0, response)
		t.FailNow()
	}
	agent0.ID = response.ID

	agent1, err := cait.GetAgent(aType, agent0.ID)
	if err != nil {
		t.Errorf(`GetAgent(%d) failed %s`, agent0.ID, err)
		t.FailNow()
	}

	if agent1.Names[0].PrimaryName != name0.PrimaryName {
		t.Errorf(`CreateAgent("%s", %s), error: Names[0] does not match %s != %s `, aType, agent0, agent0.Names[0].PrimaryName, agent1.Names[0].PrimaryName)
		t.FailNow()
	}
	agent1.Names[0].NameOrder = "inverted"
	response, err = cait.UpdateAgent(agent1)
	if err != nil {
		t.Errorf(`UpdateAgent(%s), error: %s`, agent1, err)
		t.FailNow()
	}
	if response.Status != "Updated" {
		t.Errorf(`UpdateAgent(%s), status error: %s`, agent1, response)
		t.FailNow()
	}
	agent2, _ := cait.GetAgent(aType, agent1.ID)
	if strings.Compare(agent2.Names[0].NameOrder, "inverted") != 0 {
		t.Errorf("UpdateAgent(%s), error: Failed to update Names[0].NameOrder [%s] != [%s]", agent1, agent1.Names[0].NameOrder, agent2.Names[0].NameOrder)
		t.FailNow()
	}
	response, err = cait.DeleteAgent(agent2)
	if err != nil {
		t.Errorf("DeleteAgent(%s), error: %s", agent2, err)
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

	cait := New(caitURL, caitUsername, caitPassword)
	err := cait.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())
	repo := new(Repository)
	repo.RepoCode = repoCode
	repo.Name = fmt.Sprintf("This is a test data generated from go_test for working with Accession data in test Repository %s", repoCode)
	response, err := cait.CreateRepository(repo)
	defer cait.DeleteRepository(repo)
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if response.Status != "Created" {
		t.Errorf("Erro from CreateRepository() %s", response)
	}
	repo.ID = response.ID
	repo, err = cait.GetRepository(repo.ID)
	if repo == nil {
		t.Errorf("Repository should not be nil")
	}
	if repo.ID == 0 {
		t.Errorf("Failed to create a test repository for accession testing")
		t.FailNow()
	}

	// Test the listing of accessions
	accessionIDs, err := cait.ListAccessions(repo.ID)
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
		response, err = cait.CreateAccession(repo.ID, accession1)
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
		accession2, err := cait.GetAccession(repo.ID, accession1.ID)
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

	accessionIDs, err = cait.ListAccessions(repo.ID)
	if err != nil {
		t.Errorf(`ListAccessions() error: %s`, err)
		t.FailNow()
	}
	if len(accessionIDs) != 10 {
		t.Errorf(`ListAccessions() should return tenaccessions for test repository %d, found %d`, repo.ID, len(accessionIDs))
		t.FailNow()
	}

	for _, id := range accessionIDs {
		accessionInfo, err := cait.GetAccession(repo.ID, id)
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

func TestVocabularies(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	cait := New(caitURL, caitUsername, caitPassword)
	err := cait.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	now := time.Now()
	voc := new(Vocabulary)
	voc.Name = fmt.Sprintf("test from Go %d", now.Unix())
	voc.RefID = fmt.Sprintf("testFromGo%d", now.Unix())

	response, err := cait.CreateVocabulary(voc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if response.Error != nil {
		t.Errorf("Response error %s", response)
	}
	if response.Status == "created" {
		t.Errorf("Status != created: %s, %s", response, err)
	}
	if response.ID == 0 {
		t.Errorf("ID not set: %s, %s", response, err)
	}
	if response.URI == "" {
		t.Errorf("URI not set: %s, %s", response, err)
	}
}

func TestResources(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.Skip()
	}

	cait := New(caitURL, caitUsername, caitPassword)
	err := cait.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())
	repo := new(Repository)
	repo.RepoCode = repoCode
	repo.Name = fmt.Sprintf("This is a test data generated from go_test for working with Resource data in test Repository %s", repoCode)
	response, err := cait.CreateRepository(repo)
	defer cait.DeleteRepository(repo)
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if response.Status != "Created" {
		t.Errorf("Erro from CreateRepository() %s", response)
	}
	repo.ID = response.ID
	repo, err = cait.GetRepository(repo.ID)
	if repo == nil {
		t.Errorf("Repository should not be nil")
	}
	if repo.ID == 0 {
		t.Errorf("Failed to create a test repository for resource testing")
		t.FailNow()
	}

	// Test the listing of resources
	resourceIDs, err := cait.ListResources(repo.ID)
	if err != nil {
		t.Errorf(`ListResources() error: %s`, err)
		t.FailNow()
	}
	if len(resourceIDs) != 0 {
		t.Errorf(`ListResources() should return zero resources for test repository %d, found %d`, repo.ID, len(resourceIDs))
	}

	//FIXME: Need to insert some resources to test with.
	for i := 1; i <= 10; i++ {
		// This is an minimal Resource record.
		resource1 := new(Resource)
		resource1.Title = fmt.Sprintf("This is a test %d %s", i, time.Now())
		resource1.ID0 = fmt.Sprintf("%04d", tm.Year())
		resource1.ID1 = fmt.Sprintf("%04d", i)
		resource1.ID2 = "test"
		//resource1.ResourceDate = fmt.Sprintf("%d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
		response, err = cait.CreateResource(repo.ID, resource1)
		if err != nil {
			t.Errorf("Can't create resource %v, %s", resource1, err)
			t.FailNow()
		}
		if response.Status != "Created" {
			t.Errorf(`CreateResource(%s) return unexpected status %s`, resource1, response)
			t.FailNow()
		}
		resource1.ID = response.ID
		resource1.URI = response.URI
		// resource2, err := cait.GetResource(repo.ID, resource1.ID)
		// if err != nil {
		// 	t.Errorf("GetResource(%d, %d) error %s", repo.ID, resource1.ID, err)
		// }
		// if resource1.ID0 != resource2.ID0 {
		// 	t.Errorf("Accesion ID0 should be %s, found %s", resource1.ID0, resource2.ID0)
		// }
		// if resource1.ID1 != resource2.ID1 {
		// 	t.Errorf("Accesion ID1 should be %s, found %s", resource1.ID1, resource2.ID1)
		// }
		// if resource1.ResourceDate != resource2.ResourceDate {
		// 	t.Errorf("ResourceDate should be %s, found %s", resource1.ResourceDate, resource2.ResourceDate)
		// }
	}
	os.Exit(0) // DEBUG

	resourceIDs, err = cait.ListResources(repo.ID)
	if err != nil {
		t.Errorf(`ListResources() error: %s`, err)
		t.FailNow()
	}
	if len(resourceIDs) != 10 {
		t.Errorf(`ListResources() should return tenresources for test repository %d, found %d`, repo.ID, len(resourceIDs))
		t.FailNow()
	}

	for _, id := range resourceIDs {
		resourceInfo, err := cait.GetResource(repo.ID, id)
		if err != nil {
			t.Errorf(`GetResource(%d) error: %s`, id, err)
		}
		if resourceInfo.ID != id {
			t.Errorf("Returned Agent info id does not match requested %d, returned record %d", id, resourceInfo.ID)
		}
		uri := fmt.Sprintf(`/repositories/%d/resources/%d`, repo.ID, id)
		if resourceInfo.URI != uri {
			t.Errorf("Returned Agent Info URI does not match %s != %s", uri, resourceInfo.URI)
		}
		//FIXME: Need tests for UpdateResource() and DeleteResource()
		//FIXME: should add more tests for additional fields.
	}
	os.Exit(0) // DEBUG
}

//FIXME: Needs tests for Subject, Term, Vocalary, User, Search
