/**
 * cmds/gospace.go - A command line utility using the gospace package to work
 * with ArchivesSpace's REST API.
 */
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"path"
	"../../../gospace"
	"time"
	"encoding/json"
)

func usage(msg string, exitCode int) {
	appName := path.Base(os.Args[0])
	USAGE_TEXT := fmt.Sprintf(`
  USAGE: %s SUBJECT ACTION OBJECT [OPTIONS]

  Synopsis: %s is a command line utility for interacting with an ArchivesSpace
  instance.  The command is tructure around an SUBJECT, ACTION and a OBJECT.

  SUBJECT can be repository, agent, or accession.

  ACTION can be one of create, list, update, and delete.

  OBJECT is the object of the subject and action (e.g. a Repository ID, Agent ID,
	  Accession ID, etc)

  OPTIONS can be any additional values need to specify or clarify the OBJECT

  %s also relies on the shell environment for information about connecting
  to the ArchivesSpace instance. The following shell variables are used

	ASPACE_PROTOCOL          %s
	ASPACE_HOST              %s
	ASPACE_PORT              %s
	ASPACE_USERNAME          %s
	ASPACE_PASSWORD          %s
	ASPACE_REPOSITORY_ID     %s
	ASPACE_REPOSITORY_NAME   %s


  Example:

  	%s repository create 2 "My Test Repository"

  The subject is "repository", the action is "create", the target is "2"
  and the options are "My Test Repository".

  This would create a test repository with the ID of 2 and the description of
  "My Test Repository"

`,
	appName,
	appName,
	appName,
	os.Getenv("ASPACE_PROTOCOL"),
	os.Getenv("ASPACE_HOST"),
	os.Getenv("ASPACE_PORT"),
	os.Getenv("ASPACE_USERNAME"),
	os.Getenv("ASPACE_PASSWORD"),
	os.Getenv("ASPACE_REPOSITORY_ID"),
	os.Getenv("ASPACE_REPOSITORY_NAME"),
	appName)

	fmt.Fprintln(os.Stderr, USAGE_TEXT)
	if msg != "" {
		fmt.Fprintf(os.Stderr, " %s\n\n", msg)
	}
	os.Exit(exitCode)
}


func configureApp() (map[string]string, error){
	envKeys := []string{
		"ASPACE_PROTOCOL",
		"ASPACE_HOST",
		"ASPACE_PORT",
		"ASPACE_USERNAME",
		"ASPACE_PASSWORD",
		"ASPACE_REPOSITORY_ID",
		"ASPACE_REPOSITORY_NAME",
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

type command struct {
	Subject string
	Action string
	Object string
	Options []string
}

func containsElement(src []string, elem string) bool {
	for _, item := range src {
		if strings.Compare(item, elem) == 0 {
			return true
		}
	}
	return false
}

func parseCmd(args []string)(*command, error) {
	if len(args) < 3 {
		return  nil, fmt.Errorf("Commands have the form SUBJECT ACTION OBJECT [OPTIONS]")
	}
	subjects := []string{
		"repository",
		"agent",
		"accession",
	}
	actions := []string{
		"create",
		"list",
		"update",
		"delete",
	}

	cmd := new(command)

	if containsElement(subjects, args[0]) == false {
		return nil, fmt.Errorf("%s is not a subject (e.g. %s)", args[0], strings.Join(subjects, ", "))
	}

	if containsElement(actions, args[1]) == false {
		return nil, fmt.Errorf("%s is not an action (e.g. %s)", args[1], strings.Join(actions, ", "))
	}

	cmd.Subject = args[0]
	cmd.Action = args[1]
	cmd.Object = args[2]
	if len(args) > 2 {
		cmd.Options = args[3:]
	}
	return cmd, nil
}

func runRepoCmd(cmd *command, config map[string]string) error {
	api := gospace.New(config["ASPACE_PROTOCOL"], config["ASPACE_HOST"], config["ASPACE_PORT"], config["ASPACE_USERNAME"], config["ASPACE_PASSWORD"])
	if err := api.Login(); err != nil {
		return err
	}
	switch cmd.Action {
	case "create":
		_, err := api.CreateRepository(cmd.Object, strings.Join(cmd.Options, " "))
		return err
	case "list":
		if strings.Compare(cmd.Object, "all") == 0 {
			repos, err := api.ListRepositories()
			if err != nil {
				fmt.Fprintf(os.Stderr, `{"status": "error", "message": "%s"}`, err)
				os.Exit(1)
			}
			src, err := json.Marshal(repos)
			if err != nil {
				fmt.Fprintf(os.Stderr, `{"status": "error", "message": "Cannot JSON encode %s %s"}`, cmd.Object, err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", src)
			os.Exit(0)
		} else {
			repoID, err := strconv.Atoi(cmd.Object)
			if err != nil {
				fmt.Fprintf(os.Stderr, `{"status": "error", "message": "Cannot convert %s to a number %s"}`, cmd.Object, err)
				os.Exit(1)
			}
			repo, err := api.GetRepository(repoID)
			if err != nil {
				fmt.Fprintf(os.Stderr, `{"status": "error", "message": "%s"}`, err)
				os.Exit(1)
			}
			src, err := json.Marshal(repo)
			if err != nil {
				fmt.Fprintf(os.Stderr, `{"status": "error", "message": "Cannot find %s %s"}`, cmd.Object, err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", src)
			os.Exit(0)
		}
	case "update":
		repo := new(gospace.Repository)
		err := json.Unmarshal([]byte(cmd.Object), &repo);
		if err != nil {
			return err
		}
		return api.UpdateRepository(repo)
	case "delete":
		repoID, err := strconv.Atoi(cmd.Object)
		if err != nil {
			return err
		}
		repo, err := api.GetRepository(repoID)
		if err != nil {
			return err
		}
		return api.DeleteRepository(repo)
	}
	return fmt.Errorf("action %s not implemented for %s", cmd.Action, cmd.)
}

func runCmd(cmd *command, config map[string]string) error {
	if strings.Compare(cmd.Subject, "repository") == 0 {
		return runRepoCmd(cmd, config)
	}
	return fmt.Errorf("%v not implemented.", cmd)
}

func main() {
	t1 := time.Now()
	if len(os.Args) < 2 {
		usage("gospace is a command line tool for interacting with an ArchivesSpace installation.", 1)
	}
	config, err := configureApp()
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	cmd, err := parseCmd(os.Args[1:])
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}

	err = runCmd(cmd, config)
	if err != nil {
		usage(fmt.Sprintf("%s", err), 1)
	}
	fmt.Printf("Done. %s seconds.\n\n", time.Since(t1).String())
	os.Exit(0)
}
