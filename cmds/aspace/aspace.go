/**
 * cmds/aspace/aspace.go - A command line utility using the aspace package to work
 * with ArchivesSpace's REST API.
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"../../../aspace"
	"github.com/blevesearch/bleve"
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
		"instance",
		"repository",
		"agent",
		"accession",
		"subject",
		"vocabulary",
		"term",
		"location",
		"search",
	}
	actions = []string{
		"create",
		"list",
		"update",
		"delete",
		"export",
		"import",
	}
)

// These are the global environment variables defaults used by various combinations of subjects and actions
var (
	aspaceAPIURL     = `http://localhost:8089`
	aspaceUsername   = ``
	aspacePassword   = ``
	aspaceDataSet    = `data`
	aspaceHtdocs     = `htdocs`
	aspaceTemplates  = `templates`
	aspaceBleveIndex = `index.bleve`
)

func usage(msg string, exitCode int) {
	appName := "aspace" //path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, `
  USAGE: %s SUBJECT ACTION [OPTIONS|PAYLOAD]

  %s is a command line utility for interacting with an ArchivesSpace
  instance.  The command is tructure around an SUBJECT, ACTION and an optional PAYLOAD

  SUBJECT can be %s.

  ACTION can be %s.

  PAYLOAD is a JSON expression appropriate to SUBJECT and ACTION.

  OPTIONS addition flags based parameters appropriate apply to the SUBJECT, ACTION or PAYLOAD

`,
		appName,
		appName,
		strings.Join(subjects, ", "),
		strings.Join(actions, ", "))

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "\t-%s\t%s\n", f.Name, f.Usage)
	})

	fmt.Fprintf(os.Stderr, `
  %s also relies on the shell environment for information about connecting
  to the ArchivesSpace instance. The following shell variables are used

	ASPACE_API_URL           (e.g. http://localhost:8089)
	ASPACE_API_TOKEN         (e.g. long token string of letters and numbers)

  If ASPACE_API_TOKEN is not set then ASPACE_USERNAME and ASPACE_PASSWORD
  are used if available.

  EXAMPLES:

  	%s repository create '{"repo_code":"MyTest","name":"My Test Repository"}'

  The subject is "repository", the action is "create", the target is "MyTest"
  and the options are "My Test Repository".

  This would create a test repository with a repo code of "MyTest" and a name of
  "My Test Repository".

  You can check to see what repositories exists with

    %s repository list

  Or for a specific repository by ID with

    %s repository list '{"id": 2}'

  Other SUBJECTS and ACTIONS work in a similar fashion.

`,
		appName,
		appName,
		appName,
		appName)

	if msg != "" {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", msg)
	}

	os.Exit(exitCode)
}

func containsElement(src []string, elem string) bool {
	for _, item := range src {
		if strings.Compare(item, elem) == 0 {
			return true
		}
	}
	return false
}

func exportInstance(api *aspace.ArchivesSpaceAPI) error {
	var err error

	log.Println("Logging into ", api.URL)
	err = api.Login()
	if err != nil {
		return fmt.Errorf("%s, error %s", api.URL, err)
	}
	log.Printf("export TOKEN=%s\n", api.AuthToken)

	err = api.ExportRepositories()
	if err != nil {
		return fmt.Errorf("Can't export repositories, %s", err)
	}

	for _, agentType := range []string{"people", "corporate_entities", "families", "software"} {
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
		err = api.ExportAccessions(id)
		if err != nil {
			return fmt.Errorf("Can't export repositories/%d/accessions, %s", id, err)
		}
	}

	err = api.ExportSubjects()
	if err != nil {
		return fmt.Errorf("Can't export subjects, %s", err)
	}

	err = api.ExportVocabularies()
	if err != nil {
		return fmt.Errorf("Can't export vocabularies, %s", err)
	}

	err = api.ExportTerms()
	if err != nil {
		return fmt.Errorf("Can't export terms, %s", err)
	}

	err = api.ExportLocations()
	if err != nil {
		return fmt.Errorf("Can't export locations, %s", err)
	}

	//FIXME: Add other types as we start to use them
	//FIXME: E.g. DigitalObject, Instances, Extents, Resource, Group, Users
	return nil
}

func importInstance(api *aspace.ArchivesSpaceAPI) error {
	return fmt.Errorf(`importInstance("%s") not implemented`, api)
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

	if cmd.Subject == "search" {
		cmd.Action = ""
		cmd.Payload = strings.Join(args[1:], " ")
		return cmd, nil
	}

	if containsElement(actions, args[1]) == false {
		return nil, fmt.Errorf("%s is not an action (e.g. %s)", args[1], strings.Join(actions, ", "))
	}
	cmd.Action = args[1]
	if len(args) > 2 {
		cmd.Payload = strings.Join(args[2:], " ")
	}
	return cmd, nil
}

func runInstanceCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	switch cmd.Action {
	case "export":
		return "", exportInstance(api)
	case "import":
		return "", importInstance(api)
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runRepoCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	switch cmd.Action {
	case "create":
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), repo)
		response, err := api.CreateRepository(repo)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		// repo, err = api.GetRepository(response.ID)
		// src, err := json.Marshal(repo)
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if cmd.Payload == "" {
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
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		repoID := repo.ID
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot convert %s to a number %s"}`, cmd.Payload, err)
		}
		repo, err = api.GetRepository(repoID)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		src, err := json.Marshal(repo)
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot find %s %s"}`, cmd.Payload, err)
		}
		return string(src), nil
	case "update":
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.UpdateRepository(repo)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		repo, err = api.GetRepository(repo.ID)
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
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		err = api.ExportRepository(
			repo.ID,
			path.Join(api.DataSet, "repositories"),
			fmt.Sprintf("%d.json", repo.ID),
		)
		if err != nil {
			return "", err
		}
		return `{"status": "ok"}`, nil
	case "import":
		repo := new(aspace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		err = api.ImportRepository(path.Join(api.DataSet, "repositories", fmt.Sprintf("%d.json", repo.ID)))
		if err != nil {
			return "", err
		}
		return `{"status": "ok"}`, nil
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runAgentCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	//Agent Type Payload as JSON encoded objects
	agent := new(aspace.Agent)
	err := json.Unmarshal([]byte(cmd.Payload), &agent)
	if err != nil {
		return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
	}
	p := strings.Split(agent.URI, "/")
	if len(p) < 3 {
		return "", fmt.Errorf(`Agent commands require a uri in the JSON payload, %s`, cmd.Payload)
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
			return "", fmt.Errorf("%s", response)
		}
		// agent, err = api.GetAgent(aType, response.ID)
		// if err != nil {
		// 	return "", err
		// }
		// src, err := json.Marshal(agent)
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if agent.ID == 0 {
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
		agentID := agent.ID
		if err != nil {
			return "", fmt.Errorf(`{"error": "Cannot convert %s to a number %s"}`, cmd.Payload, err)
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
		agent, err = api.GetAgent(aType, agent.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteAgent(agent)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runAccessionCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	// Repo ID is passed as a JSON object
	accession := new(aspace.Accession)
	err := json.Unmarshal([]byte(cmd.Payload), &accession)
	if err != nil {
		return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
	}
	repoID := 0
	accessionID := accession.ID
	ref, ok := accession.Repository["ref"]
	if ok {
		p := strings.Split(ref, "/")
		repoID, err = strconv.Atoi(p[len(p)-1])
		if err != nil {
			repoID = 0
		}
	}
	if repoID == 0 {
		p := strings.Split(accession.URI, "/")
		if len(p) > 2 {
			repoID, err = strconv.Atoi(p[2])
			if err != nil {
				return "", fmt.Errorf(`{"error":"Could not determine repository id"}`)
			}
		}
	}
	if accessionID == 0 {
		p := strings.Split(accession.URI, "/")
		if len(p) > 4 {
			accessionID, err = strconv.Atoi(p[4])
			if err != nil {
				return "", fmt.Errorf(`{"error":"Could not determine accession id"}`)
			}
		}
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateAccession(repoID, accession)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if accession.ID == 0 {
			accessions, err := api.ListAccessions(repoID)
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
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
		accession, err = api.GetAccession(repoID, accession.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteAccession(accession)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runSubjectCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	subject := new(aspace.Subject)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &subject)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateSubject(subject)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if subject.ID == 0 {
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
		subject, err := api.GetSubject(subject.ID)
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
		subject, err := api.GetSubject(subject.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteSubject(subject)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runLocationCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	location := new(aspace.Location)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &location)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateLocation(location)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if location.ID == 0 {
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
		location, err := api.GetLocation(location.ID)
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
		location, err := api.GetLocation(location.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteLocation(location)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runVocabularyCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	vocabulary := new(aspace.Vocabulary)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &vocabulary)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateVocabulary(vocabulary)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		if vocabulary.ID == 0 {
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
		vocabulary, err := api.GetVocabulary(vocabulary.ID)
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
		vocabulary, err := api.GetVocabulary(vocabulary.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteVocabulary(vocabulary)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runTermCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	if err := api.Login(); err != nil {
		return "", err
	}
	term := new(aspace.Term)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &term)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	var (
		vocabularyID int
		err          error
	)
	p := strings.Split(term.URI, "/")
	if len(p) > 3 {
		vocabularyID, err = strconv.Atoi(p[2])
		if err != nil {
			return "", fmt.Errorf(`{"error":%q}`, err)
		}
	}
	switch cmd.Action {
	case "create":
		response, err := api.CreateTerm(vocabularyID, term)
		if err != nil {
			return "", err
		}
		if response.Status != "Created" {
			return "", fmt.Errorf("%s", response)
		}
		src, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(src), nil
	case "list":
		//FIXME: calculate the vocabulary ID
		if term.ID == 0 {
			var ids []int
			ids, err := api.ListTerms(vocabularyID)
			if err != nil {
				return "", fmt.Errorf(`{"error": "%s"}`, err)
			}
			src, err := json.Marshal(ids)
			if err != nil {
				return "", fmt.Errorf(`{"error": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		term, err := api.GetTerm(vocabularyID, term.ID)
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
		term, err := api.GetTerm(vocabularyID, term.ID)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.DeleteTerm(term)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runSearchCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	var (
		opt *aspace.SearchQuery
		err error
	)
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &opt)
		if err != nil {
			return "", fmt.Errorf("Could not decode %s, error: %s", cmd.Payload, err)
		}
	}
	bleveIndex := os.Getenv("ASPACE_BLEVE_INDEX")
	if bleveIndex == "" {
		// Fall back to the ArchivesSpace search API
		if err := api.Login(); err != nil {
			return "", err
		}
		results, err := api.Search(opt)
		if err != nil {
			return "", fmt.Errorf(`{"error": "%s"}`, err)
		}
		return string(results), nil
	}

	// search for some text
	index, err := bleve.Open(bleveIndex)
	if err != nil {
		return "", fmt.Errorf("Can't open index %s, %s", bleveIndex, err)
	}
	defer index.Close()
	query := bleve.NewMatchQuery(opt.Q)
	if opt.PageSize == 0 {
		opt.PageSize = 10
	}
	request := bleve.NewSearchRequestOptions(query, opt.PageSize, opt.Page, opt.Explain)

	results, err := index.Search(request)
	if err != nil {
		return "", fmt.Errorf("Search error, terms [%s], %s", opt.Q, err)
	}
	return fmt.Sprintf("%s", results), nil
}

func runCmd(api *aspace.ArchivesSpaceAPI, cmd *command) (string, error) {
	switch cmd.Subject {
	case "instance":
		return runInstanceCmd(api, cmd)
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
	case "search":
		return runSearchCmd(api, cmd)
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

func main() {
	flag.BoolVar(help, "h", false, "Display the help page")
	flag.StringVar(payload, "i", "", "Use this filepath for the payload")
	flag.BoolVar(version, "v", false, "Display version info")

	aspaceAPIURL = aspace.MergeEnv("ASPACE_API_URL", aspaceAPIURL)
	aspaceUsername = aspace.MergeEnv("ASPACE_USERNAME", aspaceUsername)
	aspacePassword = aspace.MergeEnv("ASPACE_PASSWORD", aspacePassword)
	aspaceDataSet = aspace.MergeEnv("ASPACE_DATASET", aspaceDataSet)
	aspaceHtdocs = aspace.MergeEnv("ASPACE_HTDOCS", aspaceHtdocs)
	aspaceTemplates = aspace.MergeEnv("ASPACE_TEMPLATES", aspaceTemplates)
	aspaceBleveIndex = aspace.MergeEnv("ASPACE_BLEVE_INDEX", aspaceBleveIndex)

	api := aspace.New(aspaceAPIURL, aspaceUsername, aspacePassword)

	flag.Parse()
	if *help == true {
		usage("", 0)
	}

	if *version == true {
		fmt.Printf("Version: %s\n", aspace.Version)
		os.Exit(0)
	}

	args := os.Args[1:]
	if len(args) < 2 {
		usage("aspace is a command line tool for interacting with an ArchivesSpace installation.", 1)
	}

	cmd, err := parseCmd(args)
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	os.Args = args[1:]

	flag.Parse()

	if *help == true {
		usage("", 0)
	}

	if *version == true {
		fmt.Printf("Version: %s\n", aspace.Version)
		os.Exit(0)
	}

	if *payload != "" {
		src, err := ioutil.ReadFile(*payload)
		if err != nil {
			usage(fmt.Sprintf("Cannot read %s", *payload), 1)
		}
		cmd.Payload = string(src)
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
