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
	"path"
)

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
		dir := path.Join(api.Dataset, "repositories")
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
		dir := path.Join(api.Dataset, "agents", agentType)
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
		dir := path.Join(api.Dataset, "repositories", fmt.Sprintf("%d", repoID), "accessions")
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
		dir := path.Join(api.Dataset, "subjects")
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
		dir := path.Join(api.Dataset, "vocabularies")
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
			dir := path.Join(api.Dataset, "vocabularies", fmt.Sprintf("%d", vocID), "terms")
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
		dir := path.Join(api.Dataset, "locations")
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
		dir := path.Join(api.Dataset, "repositories", fmt.Sprintf("%d", repoID), "digital_objects")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
	}
	return nil
}

// ExportResources export all resources by id to JSON files.
func (api *ArchivesSpaceAPI) ExportResources(repoID int) error {
	ids, err := api.ListResources(repoID)
	if err != nil {
		return fmt.Errorf("Can't list resource ids, %s", err)
	}
	for _, id := range ids {
		data, err := api.GetResource(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get /repositories/%d/resources/%d, %s", repoID, id, err)
		}
		dir := path.Join(api.Dataset, "repositories", fmt.Sprintf("%d", repoID), "resources")
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(&data, dir, fname)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
	}
	return nil
}

// ExportArchivesSpace exports all content currently support by the Golang API implementation
func (api *ArchivesSpaceAPI) ExportArchivesSpace() error {
	var err error

	log.Println("Exporting repositories")
	err = api.ExportRepositories()
	if err != nil {
		return fmt.Errorf("Can't export repositories, %s", err)
	}

	log.Printf("Exporting subjects\n")
	err = api.ExportSubjects()
	if err != nil {
		return fmt.Errorf("Can't export subjects, %s", err)
	}

	log.Printf("Exporting vocabularies\n")
	err = api.ExportVocabularies()
	if err != nil {
		return fmt.Errorf("Can't export vocabularies, %s", err)
	}

	log.Printf("Exporting terms")
	err = api.ExportTerms()
	if err != nil {
		return fmt.Errorf("Can't export terms, %s", err)
	}

	log.Printf("Exporting locations")
	err = api.ExportLocations()
	if err != nil {
		return fmt.Errorf("Can't export locations, %s", err)
	}

	for _, agentType := range []string{"people", "corporate_entities", "families", "software"} {
		log.Printf("Exporting agents/%s\n", agentType)
		err = api.ExportAgents(agentType)
		if err != nil {
			return fmt.Errorf("Can't export agents, %s", err)
		}
	}

	ids, err := api.ListRepositoryIDs()
	if err != nil {
		return fmt.Errorf("Can't get a list of repository ids, %s", err)
	}
	for _, id := range ids {
		log.Printf("Exporting repositories/%d/digital_objects\n", id)
		err = api.ExportDigitalObjects(id)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/digital_objects, %s", id, err)
		}
		log.Printf("Exporting repositories/%d/resources\n", id)
		err = api.ExportResources(id)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/accessions, %s", id, err)
		}
		log.Printf("Exporting repositories/%d/accessions\n", id)
		err = api.ExportAccessions(id)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/accessions, %s", id, err)
		}
	}
	log.Printf("Export complete")

	//FIXME: Add other types as we start to use them
	//FIXME: E.g. Resources, Extents, Instances, Group, Users
	return nil
}
