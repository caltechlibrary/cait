//
// cmds/cait/cait.go - A command line utility using the cait package to work
// with ArchivesSpace's REST API.
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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"../../../cait"
)

type command struct {
	Subject string
	Action  string
	Payload string
	Options []string
}

var (
	help    = flag.Bool("help", false, "Display the help page")
	payload = flag.String("input", "", "Use this filepath for the payload")
	version = flag.Bool("version", false, "Display the version info")
)

var (
	subjects = []string{
		"archivesspace",
		"repository",
		"agent",
		"accession",
		"subject",
		"vocabulary",
		"term",
		"location",
		"digital_object",
	}
	actions = []string{
		"create",
		"list",
		"update",
		"delete",
		"export",
	}
)

// These are the global environment variables defaults used by various combinations of subjects and actions
var (
	description = `
  USAGE: cait SUBJECT ACTION [OPTIONS|PAYLOAD]

  cait is a command line utility for interacting with ArchivesSpace.
  The command is tructure around an SUBJECT, ACTION and an optional PAYLOAD

    SUBJECT can be %s.

    ACTION can be %s.

    PAYLOAD is a JSON expression appropriate to SUBJECT and ACTION.

    OPTIONS addition flags based parameters appropriate apply to the SUBJECT,
	        ACTION or PAYLOAD

`

	configuration = `
 CONFIGURATION

  cait also relies on the shell environment for information about connecting
  to ArchivesSpace. The following shell variables are used

    CAIT_API_URL           (e.g. http://localhost:8089)
    CAIT_API_TOKEN         (e.g. long token string of letters and numbers)

  If CAIT_API_TOKEN is not set then CAIT_USERNAME and CAIT_PASSWORD
  are used if available.

  EXAMPLES:

  	cait repository create '{"repo_code":"MyTest","name":"My Test Repository"}'

  The subject is "repository", the action is "create", the target is "MyTest"
  and the options are "My Test Repository".

  This would create a test repository with a repo code of "MyTest" and a name of
  "My Test Repository".

  You can check to see what repositories exists with

    cait repository list

  Or for a specific repository by ID with

    cait repository list '{"uri": "/repositories/2"}'

  Other SUBJECTS and ACTIONS work in a similar fashion.

`
	caitAPIURL       = `http://localhost:8089`
	caitUsername     = ``
	caitPassword     = ``
	caitDataset      = `dataset`
	caitDatasetIndex = `dataset.bleve`
	caitHtdocs       = `htdocs`
	caitHtdocsIndex  = `htdocs.bleve`
	caitTemplates    = `templates`
)

func usage() {
	fmt.Printf(description,
		strings.Join(subjects, ", "),
		strings.Join(actions, ", "))
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func containsElement(src []string, elem string) bool {
	for _, item := range src {
		if strings.Compare(item, elem) == 0 {
			return true
		}
	}
	return false
}

func exportArchivesSpace(api *cait.ArchivesSpaceAPI) error {
	log.Println("Logging into ", api.URL)
	err := api.Login()
	if err != nil {
		return fmt.Errorf("%s, error %s", api.URL, err)
	}
	//log.Printf("export TOKEN=%s\n", api.AuthToken)
	err = api.ExportArchivesSpace()
	if err != nil {
		return fmt.Errorf("Failed to export ArchivesSpace, %s", err)
	}
	return nil
}

func parseCmd(args []string) (*command, error) {
	cmd := new(command)

	if len(args) < 2 {
		return nil, fmt.Errorf("Commands have the form SUBJECT ACTION [OPTIONS] [PAYLOAD]")
	}

	if containsElement(subjects, args[0]) == false {
		return nil, fmt.Errorf("%s is not a subject (e.g. %s)", args[0], strings.Join(subjects, ", "))
	}
	cmd.Subject = args[0]

	if containsElement(actions, args[1]) == false {
		return nil, fmt.Errorf("%s is not an action (e.g. %s)", args[1], strings.Join(actions, ", "))
	}
	cmd.Action = args[1]
	if len(args) > 2 {
		cmd.Payload = strings.Join(args[2:], " ")
	}
	return cmd, nil
}

func runArchivesSpaceCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	switch cmd.Action {
	case "export":
		return "", exportArchivesSpace(api)
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runRepoCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	repoID := 0
	repo := new(cait.Repository)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), repo)
		if err != nil {
			return "", fmt.Errorf("Problem unmashalling JSON repository request, %s", err)
		}
		repoID = cait.URIToID(repo.URI)
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateRepository(repo)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create repo status %s, %s", repo.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if repoID == 0 {
			repos, err := api.ListRepositories()
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(repos)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		repo, err := api.GetRepository(repoID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(repo)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateRepository(repo)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		repo, err := api.GetRepository(repoID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteRepository(repo)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportRepository(
			repoID,
			path.Join(api.Dataset, "repositories"),
			fmt.Sprintf("%d.json", repoID),
		)
		if err != nil {
			return "", fmt.Errorf("Exporting repositories, %s", err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runAgentCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	//Agent Type Payload as JSON encoded objects
	agent := new(cait.Agent)
	err := json.Unmarshal([]byte(cmd.Payload), &agent)
	if err != nil {
		return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
	}
	agentID := cait.URIToID(agent.URI)
	p := strings.Split(agent.URI, "/")
	if len(p) < 3 {
		return "", fmt.Errorf(`Agent commands require a uri in the JSON payload, e.g. {"uri":"/agents/people"} or {"uri":/"agents/poeple/3"}, %s`, cmd.Payload)
	} else if len(p) == 3 {
		agent.URI = ""
	}
	aType := p[2]
	switch cmd.Action {
	case "create":
		response, err := api.CreateAgent(aType, agent)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create agent status %s, %s", agent.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if agentID == 0 {
			agents, err := api.ListAgents(aType)
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(agents)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		agent, err = api.GetAgent(aType, agentID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(agent)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateAgent(agent)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		agent, err = api.GetAgent(aType, agentID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteAgent(agent)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportAgents(aType)
		if err != nil {
			return "", fmt.Errorf("Exporting /agents/%s, %s", aType, err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runAccessionCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	// Repo ID is passed as a JSON object
	accession := new(cait.Accession)
	err := json.Unmarshal([]byte(cmd.Payload), &accession)
	if err != nil {
		return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
	}
	accessionID := cait.URIToID(accession.URI)
	repoID := cait.URIToRepoID(accession.URI)
	if repoID == 0 {
		return "", fmt.Errorf(`{"error":"Could not determine repository id from uri"}`)
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateAccession(repoID, accession)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create accession status %s, %s", accession.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if accessionID == 0 {
			accessions, err := api.ListAccessions(repoID)
			if err != nil {
				return "", fmt.Errorf(`{"uri": "/repositories/%d/accessions","error": "%s"}`, repoID, err)
			}
			src, err := json.Marshal(accessions)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		accession, err = api.GetAccession(repoID, accessionID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(accession)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateAccession(accession)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		accession, err = api.GetAccession(repoID, accessionID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteAccession(accession)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportAccessions(repoID)
		if err != nil {
			return "", fmt.Errorf("Exporting repositories/%d/accessions, %s", repoID, err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runSubjectCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	subject := new(cait.Subject)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &subject)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	subjectID := cait.URIToID(subject.URI)
	switch cmd.Action {
	case "create":
		response, err := api.CreateSubject(subject)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create Subject status %s, %s", subject.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if subjectID == 0 {
			subjects, err := api.ListSubjects()
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(subjects)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		subject, err := api.GetSubject(subjectID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(subject)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateSubject(subject)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		subject, err := api.GetSubject(subjectID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteSubject(subject)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportSubjects()
		if err != nil {
			return "", fmt.Errorf("Exporting /subjects, %s", err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runLocationCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	location := new(cait.Location)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &location)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	locationID := cait.URIToID(location.URI)
	switch cmd.Action {
	case "create":
		response, err := api.CreateLocation(location)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create location status %s, %s", location.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if locationID == 0 {
			locations, err := api.ListLocations()
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(locations)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		location, err := api.GetLocation(locationID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(location)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateLocation(location)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		location, err := api.GetLocation(locationID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteLocation(location)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportLocations()
		if err != nil {
			return "", fmt.Errorf("Exporting /locations, %s", err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runVocabularyCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	vocabulary := new(cait.Vocabulary)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &vocabulary)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	vocabularyID := cait.URIToID(vocabulary.URI)
	switch cmd.Action {
	case "create":
		response, err := api.CreateVocabulary(vocabulary)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create vocabulary status %s, %s", vocabulary.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if vocabularyID == 0 {
			var ids []int
			ids, err := api.ListVocabularies()
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(ids)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		vocabulary, err := api.GetVocabulary(vocabularyID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(vocabulary)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateVocabulary(vocabulary)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		vocabulary, err := api.GetVocabulary(vocabularyID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteVocabulary(vocabulary)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportVocabularies()
		if err != nil {
			return "", fmt.Errorf("Exporting /vocabularies, %s", err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runTermCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	term := new(cait.Term)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &term)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	termID := cait.URIToID(term.URI)
	vocabularyID := cait.URIToVocabularyID(term.URI)
	if vocabularyID == 0 {
		return "", fmt.Errorf(`Can't determine vocabulary ID from uri, e.g. {"uri":"/vocabularies/1/terms"} or {"uri":"/vocabularies/1/terms/2"}, %s`, cmd.Payload)
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateTerm(vocabularyID, term)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create term status %s, %s", term.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		//FIXME: calculate the vocabulary ID
		if termID == 0 {
			var ids []int
			ids, err := api.ListTermIDs(vocabularyID)
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(ids)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		term, err := api.GetTerm(vocabularyID, termID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(term)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateTerm(term)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		//FIXME: calculate the vocabulary ID
		term, err := api.GetTerm(vocabularyID, termID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteTerm(term)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportTerms()
		if err != nil {
			return "", fmt.Errorf("Exporting /terms, %s", err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runDigitalObjectCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	obj := new(cait.DigitalObject)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &obj)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	objID := cait.URIToID(obj.URI)
	repoID := cait.URIToRepoID(obj.URI)
	if repoID == 0 {
		return "", fmt.Errorf(`Can't determine repository ID from uri, e.g. {"uri":"/repositories/2/digital_objects"} or {"uri":"/repositories/2/digital_objects/3"}`)
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateDigitalObject(repoID, obj)
		if err != nil {
			return "", fmt.Errorf("Create digital_object fialed %s, %s", obj.URI, err)
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("Create digital_object status %s, %s", obj.URI, response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", fmt.Errorf("Create digital object response %s, %s", obj.URI, err)
		}
		return string(src), nil
	case "list":
		if objID == 0 {
			objs, err := api.ListDigitalObjects(repoID)
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s", "uri": "/repositories/%d/digital_objects"}`, err, repoID)
			}
			src, err := json.Marshal(objs)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		obj, err := api.GetDigitalObject(repoID, objID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		responseMsg, err := api.UpdateDigitalObject(obj)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		obj, err := api.GetDigitalObject(repoID, objID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteDigitalObject(obj)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "export":
		err := api.ExportDigitalObjects(repoID)
		if err != nil {
			return "", fmt.Errorf("Exporting repositories/%d/digital_objects, %s", repoID, err)
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("runDigitalObjectCmd() action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runCmd(api *cait.ArchivesSpaceAPI, cmd *command) (string, error) {
	switch cmd.Subject {
	case "archivesspace":
		return runArchivesSpaceCmd(api, cmd)
	case "repository":
		return runRepoCmd(api, cmd)
	case "agent":
		return runAgentCmd(api, cmd)
	case "accession":
		return runAccessionCmd(api, cmd)
	case "subject":
		return runSubjectCmd(api, cmd)
	case "location":
		return runLocationCmd(api, cmd)
	case "vocabulary":
		return runVocabularyCmd(api, cmd)
	case "term":
		return runTermCmd(api, cmd)
	case "digital_object":
		return runDigitalObjectCmd(api, cmd)
	}
	return "", fmt.Errorf("%s %s not implemented", cmd.Subject, cmd.Action)
}

func (c *command) String() string {
	src, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return string(src)
}

func getenv(envvar, defaultValue string) string {
	tmp := os.Getenv(envvar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func main() {
	flag.BoolVar(help, "h", false, "Display the help page")
	flag.StringVar(payload, "i", "", "Use this filepath for the payload")
	flag.BoolVar(version, "v", false, "Display version info")

	caitAPIURL = getenv("CAIT_API_URL", caitAPIURL)
	caitUsername = getenv("CAIT_USERNAME", caitUsername)
	caitPassword = getenv("CAIT_PASSWORD", caitPassword)
	caitDataset = getenv("CAIT_DATASET", caitDataset)
	caitDatasetIndex = getenv("CAIT_DATASET_INDEX", caitDatasetIndex)
	caitHtdocs = getenv("CAIT_HTDOCS", caitHtdocs)
	caitHtdocsIndex = getenv("CAIT_HTDOCS_INDEX", caitHtdocsIndex)
	caitTemplates = getenv("CAIT_TEMPLATES", caitTemplates)

	api := cait.New(caitAPIURL, caitUsername, caitPassword)

	flag.Parse()
	if *help == true {
		usage()
	}

	if *version == true {
		fmt.Printf("Version: %s\n", cait.Version)
		os.Exit(0)
	}

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("Missing commands options. For more info try: cait -h")
	}

	cmd, err := parseCmd(args)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	os.Args = args[1:]

	flag.Parse()

	if *help == true {
		usage()
	}

	if *version == true {
		log.Printf("Version: %s\n", cait.Version)
		os.Exit(0)
	}

	if *payload != "" {
		src, err := ioutil.ReadFile(*payload)
		if err != nil {
			log.Fatalf("Cannot read %s", *payload)
		}
		cmd.Payload = fmt.Sprintf("%s", src)
	}

	if cmd.Subject == "agent" && len(args) > 2 {
		cmd.Options = []string{args[2]}
	}
	src, err := runCmd(api, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(src)
	os.Exit(0)
}
