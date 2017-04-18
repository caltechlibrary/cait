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
func (api *ArchivesSpaceAPI) ExportRepository(id int, fname string) error {
	dir := "repository"
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s/%s, %s", api.Dataset, dir, err)
	}
	defer c.Close()

	data, err := api.GetRepository(id)
	if err != nil {
		return fmt.Errorf("Can't get repository %d data, %s", id, err)
	}
	err = WriteJSON(c, fname, data)
	if err != nil {
		return fmt.Errorf("Can't write repository %d data, %s", id, err)
	}
	return nil
}

// ExportRepositories exports all repositories record to a JSON file by ID.
func (api *ArchivesSpaceAPI) ExportRepositories(verbose bool) error {
	ids, err := api.ListRepositoryIDs()
	if err != nil {
		return fmt.Errorf("Can't get list of repository ids, %s", err)
	}
	for i, id := range ids {
		fname := fmt.Sprintf("%d.json", id)
		err = api.ExportRepository(id, fname)
		if err != nil {
			return fmt.Errorf("Can't export repository %d data, %s", id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d repository definitions exported\n", i)
		}
	}
	return nil
}

// ExportAgents exports all agent records of a given type to JSON files by id.
func (api *ArchivesSpaceAPI) ExportAgents(agentType string, verbose bool) error {
	dir := path.Join("agents", agentType)
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()
	ids, err := api.ListAgents(agentType)
	if err != nil {
		log.Fatalf("Can't get agent ids for %s, %s", agentType, err)
	}
	for i, id := range ids {
		data, err := api.GetAgent(agentType, id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, agentType, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d agents/%s exported\n", i, agentType)
		}
	}
	return nil
}

// ExportAccessions exports all accessions by id to JSON files.
func (api *ArchivesSpaceAPI) ExportAccessions(repoID int, verbose bool) error {
	dir := fmt.Sprintf("repositories/%d/accessions", repoID)
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	ids, err := api.ListAccessions(repoID)
	if err != nil {
		return fmt.Errorf("Can't list accession ids from repository %d, %s", repoID, err)
	}
	if verbose == true {
		log.Printf("Exporting %s\n", dir)
	}
	for i, id := range ids {
		data, err := api.GetAccession(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d accessions exported from repository no. %d\n", i, repoID)
		}
	}
	return nil
}

// ExportSubjects exports all subjects by id to JSON files.
func (api *ArchivesSpaceAPI) ExportSubjects(verbose bool) error {
	dir := "subjects"
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s/%s, %s", api.Dataset, dir, err)
	}
	defer c.Close()

	ids, err := api.ListSubjects()
	if err != nil {
		return fmt.Errorf("Can't list subject ids, %s", err)
	}
	for i, id := range ids {
		data, err := api.GetSubject(id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d subjects exported\n", i)
		}
	}
	return nil
}

// ExportVocabularies exports all the vocabularies by ids to JSON files.
func (api *ArchivesSpaceAPI) ExportVocabularies(verbose bool) error {
	dir := "vocabularies"
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()
	ids, err := api.ListVocabularies()
	if err != nil {
		return fmt.Errorf("Can't list vocabulary ids, %s", err)
	}
	for i, id := range ids {
		data, err := api.GetVocabulary(id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d vocabulary terms exported\n", i)
		}
	}
	return nil
}

func (api *ArchivesSpaceAPI) collectTerms(vocID int, verbose bool) error {
	dir := path.Join("vocabularies", fmt.Sprintf("%d", vocID), "terms")
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	terms, err := api.ListTerms(vocID)
	if err != nil {
		return fmt.Errorf("Can't list term ids for %s, %s", dir, err)
	}
	for i, term := range terms {
		fname := fmt.Sprintf("%d.json", term.ID)
		err = WriteJSON(c, fname, &term)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, term.ID, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d Vocabulary terms exported\n", i)
		}
	}
	return nil
}

// ExportTerms export all terms by voc id, term id to JSON files.
func (api *ArchivesSpaceAPI) ExportTerms(verbose bool) error {
	vocIDs, err := api.ListVocabularies()
	if err != nil {
		return fmt.Errorf("Can't list vocabulary ids, %s", err)
	}

	for i, vocID := range vocIDs {
		err = api.collectTerms(vocID, verbose)
		if err != nil {
			return err
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d Vocabulary terms exported\n", i)
		}
	}
	return nil
}

// ExportLocations export all locations by id to JSON files.
func (api *ArchivesSpaceAPI) ExportLocations(verbose bool) error {
	dir := "locations"
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	ids, err := api.ListLocations()
	if err != nil {
		return fmt.Errorf("Can't list location ids, %s", err)
	}
	for i, id := range ids {
		data, err := api.GetLocation(id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d locations exported\n", i)
		}
	}
	return nil
}

// ExportDigitalObjects export all digital objects by id to JSON files.
func (api *ArchivesSpaceAPI) ExportDigitalObjects(repoID int, verbose bool) error {
	dir := path.Join("repositories", fmt.Sprintf("%d", repoID), "digital_objects")
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	ids, err := api.ListDigitalObjects(repoID)
	if err != nil {
		return fmt.Errorf("Can't list digital_object ids, %s", err)
	}
	for i, id := range ids {
		data, err := api.GetDigitalObject(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d digital objects exported\n", i)
		}
	}
	return nil
}

// ExportResources export all resources by id to JSON files.
func (api *ArchivesSpaceAPI) ExportResources(repoID int, verbose bool) error {
	dir := path.Join("repositories", fmt.Sprintf("%d", repoID), "resources")
	c, err := ApiCollection(api, dir)
	if err != nil {
		return fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	ids, err := api.ListResources(repoID)
	if err != nil {
		return fmt.Errorf("Can't list resource ids, %s", err)
	}
	for i, id := range ids {
		data, err := api.GetResource(repoID, id)
		if err != nil {
			return fmt.Errorf("Can't get %s/%d, %s", dir, id, err)
		}
		fname := fmt.Sprintf("%d.json", id)
		err = WriteJSON(c, fname, &data)
		if err != nil {
			return fmt.Errorf("Can't write %s/%d.json, %s", dir, id, err)
		}
		if verbose == true && i > 0 && (i%100) == 0 {
			log.Printf("%d resources exported\n", i)
		}
	}
	return nil
}

// ExportArchivesSpace exports all content currently support by the Golang API implementation
func (api *ArchivesSpaceAPI) ExportArchivesSpace(verbose bool) error {
	var err error

	log.Println("Exporting repositories")
	err = api.ExportRepositories(verbose)
	if err != nil {
		return fmt.Errorf("Can't export repositories, %s", err)
	}

	log.Printf("Exporting subjects\n")
	err = api.ExportSubjects(verbose)
	if err != nil {
		return fmt.Errorf("Can't export subjects, %s", err)
	}

	log.Printf("Exporting vocabularies\n")
	err = api.ExportVocabularies(verbose)
	if err != nil {
		return fmt.Errorf("Can't export vocabularies, %s", err)
	}

	log.Printf("Exporting terms")
	err = api.ExportTerms(verbose)
	if err != nil {
		return fmt.Errorf("Can't export terms, %s", err)
	}

	log.Printf("Exporting locations")
	err = api.ExportLocations(verbose)
	if err != nil {
		return fmt.Errorf("Can't export locations, %s", err)
	}

	for _, agentType := range []string{"people", "corporate_entities", "families", "software"} {
		log.Printf("Exporting agents/%s\n", agentType)
		err = api.ExportAgents(agentType, verbose)
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
		err = api.ExportDigitalObjects(id, verbose)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/digital_objects, %s", id, err)
		}
		log.Printf("Exporting repositories/%d/resources\n", id)
		err = api.ExportResources(id, verbose)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/accessions, %s", id, err)
		}
		log.Printf("Exporting repositories/%d/accessions\n", id)
		err = api.ExportAccessions(id, verbose)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/accessions, %s", id, err)
		}
	}
	log.Printf("Export complete")

	//FIXME: Add other types as we start to use them
	//FIXME: E.g. Resources, Extents, Instances, Group, Users
	return nil
}
