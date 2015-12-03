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
	"os"
	"path"
	"strings"

	"../../../gospace"
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

func usage(msg string, exitCode int) {
	appName := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, `
  USAGE: %s [OPTIONS] SUBJECT ACTION [PAYLOAD]

  %s is a command line utility for interacting with an ArchivesSpace
  instance.  The command is tructure around an SUBJECT, ACTION and an optional PAYLOAD

  OPTIONS addition flags based parameters appropriate apply to the SUBJECT, ACTION or PAYLOAD

  SUBJECT can be %s.

  ACTION can be %s.

  PAYLOAD is a JSON expression appropriate to SUBJECT and ACTION.

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

	ASPACE_PROTOCOL          %s
	ASPACE_HOST              %s
	ASPACE_PORT              %s
	ASPACE_USERNAME          %s
	ASPACE_PASSWORD          %s


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

`,
		appName,
		os.Getenv("ASPACE_PROTOCOL"),
		os.Getenv("ASPACE_HOST"),
		os.Getenv("ASPACE_PORT"),
		os.Getenv("ASPACE_USERNAME"),
		os.Getenv("ASPACE_PASSWORD"),
		appName,
		appName,
		appName)

	if msg != "" {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", msg)
	}

	os.Exit(exitCode)
}

func configureApp() (map[string]string, error) {
	envKeys := []string{
		"ASPACE_PROTOCOL",
		"ASPACE_HOST",
		"ASPACE_PORT",
		"ASPACE_USERNAME",
		"ASPACE_PASSWORD",
	}
	conf := make(map[string]string)
	for _, ky := range envKeys {
		conf[ky] = os.Getenv(ky)
		if conf[ky] == "" {
			return nil, fmt.Errorf("%s is undefined in the enviroment (e.g. try export %s=SOME_VALUE_FOR_%s)", ky, ky, ky)
		}
	}
	return conf, nil
}

func containsElement(src []string, elem string) bool {
	for _, item := range src {
		if strings.Compare(item, elem) == 0 {
			return true
		}
	}
	return false
}

func parseCmd(args []string) (*command, error) {
	cmd := new(command)

	if len(args) < 2 {
		return nil, fmt.Errorf("Commands have the form SUBJECT ACTION [OBJECT] [OPTIONS]")
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

func runInstanceCmd(cmd *command, config map[string]string) (string, error) {
	api := gospace.New(config["ASPACE_PROTOCOL"], config["ASPACE_HOST"], config["ASPACE_PORT"], config["ASPACE_USERNAME"], config["ASPACE_PASSWORD"])
	if err := api.Login(); err != nil {
		return "", err
	}
	switch cmd.Action {
	case "export":
		return "", api.ExportInstance(cmd.Payload)
	case "import":
		return "", api.ImportInstance(cmd.Payload)
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runRepoCmd(cmd *command, config map[string]string) (string, error) {
	api := gospace.New(config["ASPACE_PROTOCOL"], config["ASPACE_HOST"], config["ASPACE_PORT"], config["ASPACE_USERNAME"], config["ASPACE_PASSWORD"])
	if err := api.Login(); err != nil {
		return "", err
	}
	switch cmd.Action {
	case "create":
		repo := new(gospace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), repo)
		repo, err = api.CreateRepository(repo)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(repo)
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
		repo := new(gospace.Repository)
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
		repo := new(gospace.Repository)
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
		repo := new(gospace.Repository)
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
		return "", api.ExportInstance(cmd.Payload)
	case "import":
		return "", api.ImportInstance(cmd.Payload)
	}
	return "", fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.Subject)
}

func runAgentCmd(cmd *command, config map[string]string) (string, error) {
	api := gospace.New(config["ASPACE_PROTOCOL"], config["ASPACE_HOST"], config["ASPACE_PORT"], config["ASPACE_USERNAME"], config["ASPACE_PASSWORD"])
	if err := api.Login(); err != nil {
		return "", err
	}
	//FIXME: figure out how I want to pass in agent type
	agent := new(gospace.Agent)
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
		agent, err = api.CreateAgent(aType, agent)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(agent)
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
		agent := new(gospace.Agent)
		err := json.Unmarshal([]byte(cmd.Payload), &agent)
		if err != nil {
			return "", err
		}
		responseMsg, err := api.UpdateAgent(agent)
		if err != nil {
			return "", err
		}
		src, err := json.Marshal(responseMsg)
		return string(src), err
	case "delete":
		agent := new(gospace.Agent)
		err := json.Unmarshal([]byte(cmd.Payload), &agent)
		if err != nil {
			return "", err
		}
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

func runCmd(cmd *command, config map[string]string) (string, error) {
	switch cmd.Subject {
	case "instance":
		return runInstanceCmd(cmd, config)
	case "repository":
		return runRepoCmd(cmd, config)
	case "agent":
		return runAgentCmd(cmd, config)
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

func init() {
	flag.BoolVar(help, "h", false, "Display the help page")
	flag.StringVar(payload, "i", "", "Use this filepath for the payload")
	flag.BoolVar(version, "v", false, "Display version info")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if *help == true {
		usage("", 0)
	}

	if *version == true {
		fmt.Printf("Version: %s\n", gospace.Version)
		os.Exit(0)
	}

	if len(args) < 2 {
		usage("aspace is a command line tool for interacting with an ArchivesSpace installation.", 1)
	}
	config, err := configureApp()
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	cmd, err := parseCmd(args)
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
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
	src, err := runCmd(cmd, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(src)
	os.Exit(0)
}
