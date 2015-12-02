/**
 * gospace_test.go - Test routines for gospace.go
 */
package gospace

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Get the environment variables needed for testing.
var (
	aspaceProtocol = os.Getenv("ASPACE_PROTOCOL")
	aspaceHost     = os.Getenv("ASPACE_HOST")
	aspacePort     = os.Getenv("ASPACE_PORT")
	aspaceUsername = os.Getenv("ASPACE_USERNAME")
	aspacePassword = os.Getenv("ASPACE_PASSWORD")
)

func checkConfig(t *testing.T) bool {
	isSetup := true
	if aspaceProtocol == "" {
		t.Error("ASPACE_PROTOCOL environment variable not set.", aspaceProtocol)
		isSetup = false
	}
	if aspaceHost == "" {
		t.Error("ASPACE_HOST environment variable not set.", aspaceHost)
		isSetup = false
	}
	if aspacePort == "" {
		t.Error("ASPACE_PORT environment variable not set.", aspacePort)
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
		t.Error("Environment variables needed to run tests not configured")
		t.FailNow()
	}
}

func TestArchiveSpaceAPI(t *testing.T) {
	// Get the environment variables needed for testing.
	isSetup := checkConfig(t)
	if isSetup == false {
		t.Error("Environment variables needed to run tests not configured", isSetup)
		t.SkipNow()
	}

	aspace := New(aspaceProtocol, aspaceHost, aspacePort, aspaceUsername, aspacePassword)
	if aspace.URL == nil {
		t.Errorf("%s\t%s://%s:%s", aspace.URL.String(), aspaceProtocol, aspaceHost, aspacePort)
	}
	if strings.Compare(aspace.URL.String(), fmt.Sprintf("%s://%s:%s", aspaceProtocol, aspaceHost, aspacePort)) != 0 {
		t.Errorf("%s != %s://%s:%s\n", aspace.URL.String(), aspaceProtocol, aspaceHost, aspacePort)
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

	aspace := New(aspaceProtocol, aspaceHost, aspacePort, aspaceUsername, aspacePassword)
	tm := time.Now()
	repoCode := fmt.Sprintf("%d", tm.Unix())

	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	repo0 := new(Repository)
	repo0.RepoCode = repoCode
	repo0.Name = "This is a test generated from go_test"
	repo1, err := aspace.CreateRepository(repo0)
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if repo1 == nil {
		t.Errorf("Repository should not be nil")
	}

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
	_, err = aspace.UpdateRepository(repo2)
	if err != nil {
		t.Errorf("UpdateRepository failed for %v: %s", repo2, err)
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
