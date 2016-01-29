//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2016
// Caltech Library
//
package aspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// WriteJSON write out an ArchivesSpace data structure as a JSON file.
func WriteJSON(data interface{}, dir string, fname string) error {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s/%s, %s", dir, fname, err)
	}
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s/%s, %s", dir, fname, err)
	}
	return ioutil.WriteFile(path.Join(dir, fname), src, 0664)
}

// ExportRepository for specific id to a JSON file.
func (api *ArchivesSpaceAPI) ExportRepository(id int, dir string, fname string) error {
	data, err := api.GetRepository(id)
	if err != nil {
		return fmt.Errorf("Can't get repository %d data, %s", id, err)
	}
	err = WriteJSON(data, dir, fname)
	if err != nil {
		return fmt.Errorf("Can't write repository %d data, %s", id, err)
	}
	return nil
}

// ExportRepositories exports all repositories record to a JSON file by ID.
func (api *ArchivesSpaceAPI) ExportRepositories() error {
	ids, err := api.ListRepositoryIDs()
	if err != nil {
		return fmt.Errorf("Can't get list of repository ids, %s", err)
	}
	for _, id := range ids {
		dir := path.Join(api.DataSet, "repositories")
		fname := fmt.Sprintf("%d.json", id)
		err = api.ExportRepository(id, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't export repository %d data, %s", id, err)
		}
	}
	return nil
}

// ExportAgents exports all agent records of a given type to JSON files by id.
func (api *ArchivesSpaceAPI) ExportAgents(agentType string) error {
	ids, err := api.ListAgents(agentType)
	if err != nil {
		log.Fatalf("Can't get agent ids for %s, %s", agentType, err)
	}
	for _, id := range ids {
		data, err := api.GetAgent(agentType, id)
		if err != nil {
			return fmt.Errorf("Can't get agents/%s/%d, %s", agentType, id, err)
		}
		dir := path.Join(api.DataSet, "agents", agentType)
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write agents/%s/%d.json, %s", agentType, id, err)
		}
	}
	return nil
}

// ExportAccessions exports all accessions by id to JSON files.
func (api *ArchivesSpaceAPI) ExportAccessions(repoID int) error {
	ids, err := api.ListAccessions(repoID)
	if err != nil {
		return fmt.Errorf("Can't list accession ids from repository %d, %s", repoID, err)
	}
	for _, id := range ids {
		data, err := api.GetAccession(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get repositories/%d/accession/%d, %s", repoID, id, err)
		}
		dir := path.Join(api.DataSet, "repositories", fmt.Sprintf("%d", repoID), "accessions")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write repositories/%d/accession/%d.json, %s", repoID, id, err)
		}
	}
	return nil
}

// ExportSubjects exports all subjects by id to JSON files.
func (api *ArchivesSpaceAPI) ExportSubjects() error {
	ids, err := api.ListSubjects()
	if err != nil {
		return fmt.Errorf("Can't list subject ids, %s", err)
	}
	for _, id := range ids {
		data, err := api.GetSubject(id)
		if err != nil {
			return fmt.Errorf("Can't get subjects/%d, %s", id, err)
		}
		dir := path.Join(api.DataSet, "subjects")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write subjects/%d.json, %s", id, err)
		}
	}
	return nil
}

// ExportVocabularies exports all the vocabularies by ids to JSON files.
func (api *ArchivesSpaceAPI) ExportVocabularies() error {
	ids, err := api.ListVocabularies()
	if err != nil {
		return fmt.Errorf("Can't list vocabulary ids, %s", err)
	}
	for _, id := range ids {
		data, err := api.GetVocabulary(id)
		if err != nil {
			return fmt.Errorf("Can't get vocabularies/%d, %s", id, err)
		}
		dir := path.Join(api.DataSet, "vocabularies")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write vocabularies/%d.json, %s", id, err)
		}
	}
	return nil
}

// ExportTerms export all terms by voc id, term id to JSON files.
func (api *ArchivesSpaceAPI) ExportTerms() error {
	vocIDs, err := api.ListVocabularies()
	if err != nil {
		return fmt.Errorf("Can't list vocabulary ids, %s", err)
	}
	for _, vocID := range vocIDs {
		terms, err := api.ListTerms(vocID)
		if err != nil {
			return fmt.Errorf("Can't list term ids for vocabularies/%d, %s", vocID, err)
		}
		for _, term := range terms {
			dir := path.Join(api.DataSet, "vocabularies", fmt.Sprintf("%d", vocID), "terms")
			fname := fmt.Sprintf("%d.json", term.ID)
			err = WriteJSON(&term, dir, fname)
			if err != nil {
				return fmt.Errorf("Can't write vocabularies/%d/terms/%d.json, %s", vocID, term.ID, err)
			}
		}
	}
	return nil
}

// ExportLocations export all locations by id to JSON files.
func (api *ArchivesSpaceAPI) ExportLocations() error {
	ids, err := api.ListLocations()
	if err != nil {
		return fmt.Errorf("Can't list location ids, %s", err)
	}
	for _, id := range ids {
		data, err := api.GetLocation(id)
		if err != nil {
			return fmt.Errorf("Can't get locations/%d, %s", id, err)
		}
		dir := path.Join(api.DataSet, "locations")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write locations/%d.json, %s", id, err)
		}
	}
	return nil
}

// ExportDigitalObjects export all digital objects by id to JSON files.
func (api *ArchivesSpaceAPI) ExportDigitalObjects(repoID int) error {
	ids, err := api.ListDigitalObjects(repoID)
	if err != nil {
		return fmt.Errorf("Can't list digital_object ids, %s", err)
	}
	for _, id := range ids {
		data, err := api.GetDigitalObject(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get /repositories/%d/digial_object/%d, %s", repoID, id, err)
		}
		dir := path.Join(api.DataSet, "repositories", fmt.Sprintf("%d", repoID), "digital_objects")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write %d/%d.json, %s", dir, id, err)
		}
	}
	return nil
}
