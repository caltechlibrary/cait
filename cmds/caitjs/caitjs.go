//
// cmds/caitjs/caitjs.go - A command line JavaScript interpreter making the full cait API
// scriptable in JavaScript.
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
	// standard library
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// 3rd Party packages
	"github.com/chzyer/readline"

	// Caltech Library Pacakges
	"github.com/caltechlibrary/cait"
)

var (
	showHelp    bool
	showVersion bool
	runRepl     bool
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "display this message")
	flag.BoolVar(&showVersion, "v", false, "display version information")
	flag.BoolVar(&runRepl, "i", false, "run interactively in a REPL")
}

func main() {
	flag.Parse()

	jsFilename := ""
	jsArgs := flag.Args()

	if runRepl == false && len(jsArgs) == 0 {
		fmt.Println(`
 USAGE: caitjs [OPTIONS] JAVASCRIPT_FILENAME [OPTIONS_PASSED_TO_JAVASCRIPT_FILE]

 OPTIONS

`)
		flag.PrintDefaults()
		fmt.Printf("\nVersion %s\n", cait.Version)
		os.Exit(1)
	}

	if showHelp == true {
		fmt.Println(`
 USAGE: caitjs [OPTIONS] JAVASCRIPT_FILENAME [OPTIONS_PASSED_TO_JAVASCRIPT_FILE]

 OPTIONS

`)
		flag.PrintDefaults()
		fmt.Printf("\nVersion %s\n", cait.Version)
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Printf("Version %s\n", cait.Version)
		os.Exit(0)
	}

	caitAPIURL := os.Getenv("CAIT_API_URL")
	caitUsername := os.Getenv("CAIT_USERNAME")
	caitPassword := os.Getenv("CAIT_PASSWORD")
	if caitAPIURL == "" {
		log.Fatalf("You need to setup your environment vairables to use caitjs.")
	}

	api := cait.New(caitAPIURL, caitUsername, caitPassword)
	vm := cait.NewJavaScript(api, jsArgs)
	// if we have a script run it.
	if len(jsArgs) > 0 {
		jsFilename, jsArgs = jsArgs[0], jsArgs[1:]

		jsSrc, err := ioutil.ReadFile(jsFilename)
		if err != nil {
			log.Fatalf("Can't ready %s, %s", jsFilename, err)
		}
		script, err := vm.Compile(jsFilename, jsSrc)
		if err != nil {
			log.Fatalf("Compile error, %s", err)
		}
		_, err = vm.Run(script)
		if err != nil {
			log.Fatalf("Runtime error, %s", err)
		}
	}

	// if we need a repl run it
	if runRepl == true {
		fmt.Printf(`

   Welcome to cait %s, Caltech Archives Integration Tools
   Accessing ArchivesSpace REST API %q
   With username %q

   Autocomplete enabled for os, http and api objects
   Press Ctrl+I for line completion choices
   Press Ctrl+G cancel

   Exit catijs with Ctrl+D or type: os.exit(0);

		`, cait.Version, caitAPIURL, caitUsername)

		completer := readline.NewPrefixCompleter(
			// os object completions
			readline.PcItem("os.args()"),
			readline.PcItem("os.exit(exitCode)"),
			readline.PcItem("os.getEnv(envvar)"),
			readline.PcItem("os.readFile(filename)"),
			readline.PcItem("os.writeFile(filename, data)"),
			readline.PcItem("os.rename(oldname, newname)"),
			readline.PcItem("os.remove(filename)"),
			readline.PcItem("os.chmod(filename, perms)"),
			readline.PcItem("os.find(filename)"),
			readline.PcItem("os.mkdir(dirname)"),
			readline.PcItem("os.mkdirAll(pathDirname)"),
			readline.PcItem("os.rmdir(dirname)"),
			readline.PcItem("os.rmdirAll(pathDirname)"),
			// http object completions
			readline.PcItem("http.get(url, headers)"),
			readline.PcItem("http.post(url, headers, payload)"),
			// api object completions
			readline.PcItem("api.login()"),
			readline.PcItem("api.logout()"),
			readline.PcItem("api.createRepository(repoObject)"),
			readline.PcItem("api.getRepository(repoID)"),
			readline.PcItem("api.updateRepository(repoObject)"),
			readline.PcItem("api.deleteRepository(repoObject)"),
			readline.PcItem("api.listRepositories()"),
			readline.PcItem("api.createAgent(agentObject)"),
			readline.PcItem("api.getAgent(agentType, agentID)"),
			readline.PcItem("api.updateAgent(agentObject)"),
			readline.PcItem("api.deleteAgent(agentObject)"),
			readline.PcItem("api.listAgents(agentType)"),
			readline.PcItem("api.createAccession(accessionObject)"),
			readline.PcItem("api.getAccession(repoID, accessionID)"),
			readline.PcItem("api.updateAccession(accessionObject)"),
			readline.PcItem("api.deleteAccession(accessionObject)"),
			readline.PcItem("api.listAccessions(repoID)"),
			readline.PcItem("api.createSubject(subjectObject)"),
			readline.PcItem("api.getSubject(subjectID)"),
			readline.PcItem("api.updateSubject(subjectObject)"),
			readline.PcItem("api.deleteSubject(subjectObject)"),
			readline.PcItem("api.listSubjects()"),
			readline.PcItem("api.createDigitalObject(digitalObjectObject)"),
			readline.PcItem("api.getDigitalObject(repoID, digitalObjectID)"),
			readline.PcItem("api.updateDigitalObject(digitalObjectObject)"),
			readline.PcItem("api.deleteDigitalObject(digitalObjectObject)"),
			readline.PcItem("api.listDigitalObjects(repoID)"),
		)

		rl, err := readline.NewEx(&readline.Config{
			Prompt:       "> ",
			AutoComplete: completer,
		})
		if err != nil {
			panic(err)
		}
		defer rl.Close()

		for {
			jsSrc, err := rl.Readline()
			if err != nil { // io.EOF, readline.ErrInterrupt
				break
			}
			if len(strings.Trim(jsSrc, " ")) > 0 {
				if script, err := vm.Compile("repl", jsSrc); err != nil {
					fmt.Printf("Compile error, %s\n", err)
				} else {
					out, err := vm.Eval(script)
					switch {
					case err != nil:
						fmt.Printf("Runtime error, %s\n", err)
					default:
						fmt.Println(out.String())
					}
				}
			}
		}
	}
}
