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

func TestArchiveSpaceAPI(t *testing.T) {
	// Get the environment variables needed for testing.
	aspaceProtocol := os.Getenv("ASPACE_PROTOCOL")
	aspaceHost := os.Getenv("ASPACE_HOST")
	aspacePort := os.Getenv("ASPACE_PORT")
	aspaceUsername := os.Getenv("ASPACE_USERNAME")
	aspacePassword := os.Getenv("ASPACE_PASSWORD")

	if aspaceProtocol == "" {
		t.Error("ASPACE_PROTOCOL environment variable not set.")
		t.Fail()
	}
	if aspaceHost == "" {
		t.Error("ASPACE_HOST environment variable not set.")
		t.Fail()
	}
	if aspacePort == "" {
		t.Error("ASPACE_PORT environment variable not set.")
		t.Fail()
	}
	if aspaceUsername == "" {
		t.Error("ASPACE_USERNAME environment variable not set.")
		t.Fail()
	}
	if aspacePassword == "" {
		t.Error("ASPACE_PASSWORD environment variable not set.")
		t.Fail()
	}

	aspace := New(aspaceProtocol, aspaceHost, aspacePort, aspaceUsername, aspacePassword)
	if aspace.URL == nil {
		t.Errorf("%s\t%s://%s:%s", aspace.URL.String(), aspaceProtocol, aspaceHost, aspacePort)
		t.FailNow()
	}
	if strings.Compare(aspace.URL.String(), fmt.Sprintf("%s://%s:%s", aspaceProtocol, aspaceHost, aspacePort)) != 0 {
		t.Errorf("%s != %s://%s:%s\n", aspace.URL.String(), aspaceProtocol, aspaceHost, aspacePort)
		t.FailNow()
	}

	if aspace.IsAuth() == true {
		t.Error("aspace.IsAuth() returning true before authentication")
	}
	err := aspace.Login()
	if err != nil {
		t.Errorf("%s\t%s", err, aspace.URL.String())
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
	aspaceProtocol := os.Getenv("ASPACE_PROTOCOL")
	aspaceHost := os.Getenv("ASPACE_HOST")
	aspacePort := os.Getenv("ASPACE_PORT")
	aspaceUsername := os.Getenv("ASPACE_USERNAME")
	aspacePassword := os.Getenv("ASPACE_PASSWORD")
	aspace := New(aspaceProtocol, aspaceHost, aspacePort, aspaceUsername, aspacePassword)
	tm := time.Now()
	repoCode := fmt.Sprintf("%v", tm.Unix())

	err := aspace.Login()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	repo1, err := aspace.CreateRepository(repoCode, "This is a test generated from go_test")
	if err != nil {
		t.Errorf("Error from CreateRepository() %s", err)
	}
	if repo1 == nil {
		t.Errorf("Repository should not be nil")
	}
	fmt.Printf("DEBUG repo1: |%v|\n", repo1)

	repo2, err := aspace.GetRepository(repo1.ID)
	if err != nil {
		t.Errorf("GetRepository() error: %s", err)
		t.FailNow()
	}
	fmt.Printf("DEBUG repo2: |%v|\n", repo2)
	if repo1.ID != repo2.ID {
		t.Errorf("GetRepository() returned different IDs: %d != %d", repo1.ID, repo2.ID)
	}
	if strings.Compare(repo1.RepoCode, repo2.RepoCode) != 0 {
		t.Errorf("GetRepository() returned different RepoCode: %s != %s\n", repo1.RepoCode, repo2.RepoCode)
	}
	if strings.Compare(repo1.Name, repo2.Name) != 0 {
		t.Errorf("GetRepository() returned different RepoCode: %s != %s\n", repo1.Name, repo2.Name)
	}

	//FIXME: Need to add function and test for UpdateRepostiory, DeleteRepository, ListRepositories
}
