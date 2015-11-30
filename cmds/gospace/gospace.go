/**
 * cmds/gospace.go - A command line utility using the gospace package to work
 * with ArchivesSpace's REST API.
 */
package main

import (
	"fmt"
	"os"
)

func usage(msg string, exitCode int) {
	USAGE_TEXT := fmt.Sprintf(`
  USAGE: %s SUBJECT ACTION TARGET [OPTIONS]

  Synopsis: %s is a command line utility for interacting with an ArchivesSpace
  instance.  The command is tructure around an SUBJECT, ACTION and a TARGET.

  SUBJECT can be repository, agent, or accession.

  ACTION can be one of create, list, update, and delete.

  TARGET is the object of the subject and action (e.g. a Repository ID, Agent ID,
	  Accession ID, etc)

  OPTIONS can be any additional values need to specify the TARGET

  Example:

  	%s repository create 2 "My Test Repository"

  The subject is "repository", the action is "create", the target is "2"
  and the options are "My Test Repository".

  This would create a test repository with the ID of 2 and the description of
  "My Test Repistory"

`, os.Args[0], os.Args[0], os.Args[0])

	fmt.Fprintln(os.Stderr, USAGE_TEXT)
	if msg != "" {
		fmt.Fprintf(os.Stderr, " %s\n\n", msg)
	}
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 2 {
		usage("gospace is a general command that runs a subcomment on an ArchivesSpace installation.", 1)
	}
	fmt.Println(os.Args)
}
