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
	"fmt"
)

// ImportRepository import into ArchivesSpace a repository defined by JSON file via the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) ImportRepository(filepath string) error {
	return fmt.Errorf("ImportRepository() not implemented.")
}

// ImportRepositories import into ArchivesSpace all repositories in a directory defined by JSON files via the ArchivesSpace REST API
func (api *ArchivesSpaceAPI) ImportRepositories(dir string) error {
	return fmt.Errorf("ImportRepositories() not implemented.")
}
