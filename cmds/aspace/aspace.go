/**
 * cmds/aspace/aspace.go - A command line utility using the aspace package to work
 * with ArchivesSpace's REST API.
 */
package main

import (
	"encoding/json"
	"fmt"
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
	usageText := fmt.Sprintf(`
  USAGE: %s SUBJECT ACTION [PAYLOAD] [OPTIONS]

  Synopsis: %s is a command line utility for interacting with an ArchivesSpace
  instance.  The command is tructure around an SUBJECT, ACTION and an optional PAYLOAD

  SUBJECT can be %s.

  ACTION can be %s.

  PAYLOAD is a JSON expression appropriate to the ACTION on SUBJECT.

  OPTIONS addition flag based options appropriate to the SUBJECT, ACTION and DATA (e.g. -h, --help for help)

  %s also relies on the shell environment for information about connecting
  to the ArchivesSpace instance. The following shell variables are used

	ASPACE_PROTOCOL          %s
	ASPACE_HOST              %s
	ASPACE_PORT              %s
	ASPACE_USERNAME          %s
	ASPACE_PASSWORD          %s


  Example:

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
		appName,
		strings.Join(subjects, ", "),
		strings.Join(actions, ", "),
		appName,
		os.Getenv("ASPACE_PROTOCOL"),
		os.Getenv("ASPACE_HOST"),
		os.Getenv("ASPACE_PORT"),
		os.Getenv("ASPACE_USERNAME"),
		os.Getenv("ASPACE_PASSWORD"),
		appName,
		appName,
		appName)

	fmt.Fprintln(os.Stderr, usageText)
	if msg != "" {
		fmt.Fprintf(os.Stderr, " %s\n\n", msg)
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
				return "", fmt.Errorf(`{"status": "error", "message": "%s"}`, err)
			}
			src, err := json.Marshal(repos)
			if err != nil {
				return "", fmt.Errorf(`{"status": "error", "message": "Cannot JSON encode %s %s"}`, cmd.Payload, err)
			}
			return string(src), nil
		}
		//FIXME: need to extract something userful like the repoID form cmd.Payload
		repo := new(gospace.Repository)
		err := json.Unmarshal([]byte(cmd.Payload), &repo)
		if err != nil {
			return "", err
		}
		repoID := repo.ID
		if err != nil {
			return "", fmt.Errorf(`{"status": "error", "message": "Cannot convert %s to a number %s"}`, cmd.Payload, err)
		}
		repo, err = api.GetRepository(repoID)
		if err != nil {
			return "", fmt.Errorf(`{"status": "error", "message": "%s"}`, err)
		}
		src, err := json.Marshal(repo)
		if err != nil {
			return "", fmt.Errorf(`{"status": "error", "message": "Cannot find %s %s"}`, cmd.Payload, err)
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

func runCmd(cmd *command, config map[string]string) (string, error) {
	switch cmd.Subject {
	case "repository":
		return runRepoCmd(cmd, config)
	case "instance":
		return runRepoCmd(cmd, config)
	}
	return "", fmt.Errorf("%s %s not implemented", cmd.Subject, cmd.Action)
}

func main() {
	if len(os.Args) < 2 {
		usage("aspace is a command line tool for interacting with an ArchivesSpace installation.", 1)
	}
	config, err := configureApp()
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	cmd, err := parseCmd(os.Args[1:])
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}

	src, err := runCmd(cmd, config)
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	fmt.Println(src)
	os.Exit(0)
}
