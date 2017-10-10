//
// cmds/genpages/genpages.go - A command line utility that builds pages from the exported results of cait.go
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"text/template"

	// Caltech Library packages
	"github.com/caltechlibrary/cait"
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/tmplfn"
)

var (
	usage = `USAGE: %s [OPTIONS]`

	description = `SYNOPSIS

%s generates HTML, .include pages and normalized JSON based on the JSON output form
cait and templates associated with the command.

CONFIGURATION

%s can be configured through setting the following environment
variables-

    CAIT_DATASET    this is the directory that contains the output of the
                      'cait archivesspace export' command.

    CAIT_TEMPLATES  this is the directory that contains the templates
                      used used to generate the static content of the website.

    CAIT_HTDOCS     this is the directory where the HTML files are written.`

	// Standard Options
	showHelp    bool
	showVersion bool
	showLicense bool

	// App Options
	showVerbose bool
	htdocsDir   string
	datasetDir  string
	repoNo      string
	templateDir string
)

func loadTemplates(templateDir, aHTMLTmplName, aIncTmplName string) (*template.Template, *template.Template, error) {
	tmplFuncs := tmplfn.Join(tmplfn.AllFuncs(), cait.TmplMap)
	t := tmplfn.New(tmplFuncs)
	if err := t.ReadFiles(path.Join(templateDir, aHTMLTmplName), path.Join(templateDir, aIncTmplName)); err != nil {
		return nil, nil, fmt.Errorf("Can't read template %s, %s, %s", aHTMLTmplName, aIncTmplName, err)
	}
	aHTMLTmpl, err := t.Assemble()
	if err != nil {
		return nil, nil, fmt.Errorf("Can't parse template %s, %s, %s", aHTMLTmplName, aIncTmplName, err)
	}

	t = tmplfn.New(tmplFuncs)
	if err := t.ReadFiles(path.Join(templateDir, aIncTmplName)); err != nil {
		return nil, nil, fmt.Errorf("Can't read template %s, %s", aIncTmplName, err)
	}
	aIncTmpl, err := t.Assemble()
	if err != nil {
		return nil, nil, fmt.Errorf("Can't parse template %s, %s", aIncTmplName, err)
	}
	return aHTMLTmpl, aIncTmpl, nil
}

func processAgentsPeople(api *cait.ArchivesSpaceAPI, templateDir string, aHTMLTmplName string, aIncTmplName string, agentsPeopleDir string) (int, error) {
	log.Printf("Reading templates from %s\n", templateDir)
	aHTMLTmpl, aIncTmpl, err := loadTemplates(templateDir, aHTMLTmplName, aIncTmplName)
	if err != nil {
		return 0, fmt.Errorf("template error %q, %q: %s", aHTMLTmplName, aIncTmplName, err)
	}
	c, err := cait.OpenCollection(api, agentsPeopleDir)
	if err != nil {
		return 0, fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	keys := cait.GetKeys(c)
	cnt := 0
	for i, key := range keys {
		// Process accession records
		src, err := cait.ReadJSON(c, key)
		if err != nil {
			return cnt, err
		}
		person := new(cait.Agent)
		err = json.Unmarshal(src, &person)
		if err != nil {
			return cnt, err
		}

		// FIXME: which restrictions do we care about agent/people?--
		//        agent.Published, person.DisplayName.IsDisplayName, person.DisplayName.Authorized
		if person.Published == true && person.IsLinkedToPublishedRecord == true && person.DisplayName.IsDisplayName == true && person.DisplayName.Authorized == true {
			// Create a normalized view of the accession to make it easier to work with

			// If the accession is published and the accession is not suppressed then generate the webpage
			fname := path.Join(htdocsDir, fmt.Sprintf("%s.html", person.URI))
			dname := path.Dir(fname)
			err = os.MkdirAll(dname, 0775)
			if err != nil {
				return cnt, fmt.Errorf("Can't create %s, %s", dname, err)
			}

			// Process HTML file
			fp, err := os.Create(fname)
			if err != nil {
				return cnt, fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = aHTMLTmpl.Execute(fp, person)
			if err != nil {
				log.Fatalf("template execute error %s, %s", aHTMLTmplName, err)
				return cnt, err
			}
			fp.Close()

			// Process Include file (just the HTML content)
			fname = path.Join(htdocsDir, fmt.Sprintf("%s.include", person.URI))
			fp, err = os.Create(fname)
			if err != nil {
				return cnt, fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = aIncTmpl.Execute(fp, person)
			if err != nil {
				log.Fatalf("template execute error %s, %s", aIncTmplName, err)
				return cnt, err
			}
			fp.Close()

			// Process JSON file (an abridged version of the JSON output in data)
			fname = path.Join(htdocsDir, fmt.Sprintf("%s.json", person.URI))
			src, err := json.Marshal(person)
			if err != nil {
				return cnt, fmt.Errorf("Could not JSON encode %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = ioutil.WriteFile(fname, src, 0664)
			if err != nil {
				log.Fatalf("could not write JSON view %s, %s", fname, err)
				return cnt, err
			}
			fp.Close()
		}
		cnt = i
		if cnt > 0 && (cnt%100) == 0 {
			log.Printf("%d Agents/People\n", cnt)
		}
	}
	return cnt, nil
}

func processAccessions(api *cait.ArchivesSpaceAPI, templateDir string, aHTMLTmplName string, aIncTmplName string, accessionsDir string, agents []*cait.Agent, subjects map[string]*cait.Subject, digitalObjects map[string]*cait.DigitalObject) (int, error) {
	log.Printf("Reading templates from %s\n", templateDir)
	aHTMLTmpl, aIncTmpl, err := loadTemplates(templateDir, aHTMLTmplName, aIncTmplName)
	if err != nil {
		return 0, fmt.Errorf("template error %q, %q: %s", aHTMLTmplName, aIncTmplName, err)
	}
	c, err := cait.OpenCollection(api, accessionsDir)
	if err != nil {
		return 0, fmt.Errorf("Can't open collection %s, %s", api.Dataset, err)
	}
	defer c.Close()

	keys := cait.GetKeys(c)
	cnt := 0
	for i, key := range keys {
		// Process accession records
		src, err := cait.ReadJSON(c, key)
		if err != nil {
			return cnt, err
		}
		accession := new(cait.Accession)
		err = json.Unmarshal(src, &accession)
		if err != nil {
			return cnt, err
		}
		// FIXME: which restrictions do we care about--
		//        accession.Publish, accession.Suppressed, accession.AccessRestrictions,
		//        accession.RestrictionsApply, accession.UseRestrictions
		if accession.Publish == true && accession.Suppressed == false && accession.RestrictionsApply == false {
			// Create a normalized view of the accession to make it easier to work with
			view, err := accession.NormalizeView(agents, subjects, digitalObjects)
			if err != nil {
				return cnt, fmt.Errorf("Could not generate normalized view, %s", err)
			}

			// If the accession is published and the accession is not suppressed then generate the webpage
			fname := path.Join(htdocsDir, fmt.Sprintf("%s.html", accession.URI))
			dname := path.Dir(fname)
			err = os.MkdirAll(dname, 0775)
			if err != nil {
				return cnt, fmt.Errorf("Can't create %s, %s", dname, err)
			}

			// Process HTML file
			fp, err := os.Create(fname)
			if err != nil {
				return cnt, fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = aHTMLTmpl.Execute(fp, view)
			if err != nil {
				log.Fatalf("template execute error %s, %s", aHTMLTmplName, err)
				return cnt, err
			}
			fp.Close()

			// Process Include file (just the HTML content)
			fname = path.Join(htdocsDir, fmt.Sprintf("%s.include", accession.URI))
			fp, err = os.Create(fname)
			if err != nil {
				return cnt, fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = aIncTmpl.Execute(fp, view)
			if err != nil {
				log.Fatalf("template execute error %s, %s", aIncTmplName, err)
				return cnt, err
			}
			fp.Close()

			// Process JSON file (an abridged version of the JSON output in data)
			fname = path.Join(htdocsDir, fmt.Sprintf("%s.json", accession.URI))
			src, err := json.Marshal(view)
			if err != nil {
				return cnt, fmt.Errorf("Could not JSON encode %s, %s", fname, err)
			}
			if showVerbose == true {
				log.Printf("Writing %s", fname)
			}
			err = ioutil.WriteFile(fname, src, 0664)
			if err != nil {
				log.Fatalf("could not write JSON view %s, %s", fname, err)
				return cnt, err
			}
			fp.Close()
		}
		cnt = i
		if cnt > 0 && (cnt%100) == 0 {
			log.Printf("%d Accessions processed\n", cnt)
		}

	}
	return cnt, nil
}

func init() {
	// We are going to log to standard out rather than standard err
	log.SetOutput(os.Stdout)

	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	//flag.BoolVar(&showExamples, "example", false, "display example(s)")

	// App Options
	flag.BoolVar(&showVerbose, "verbose", false, "more verbose logging")
	flag.StringVar(&htdocsDir, "htdocs", "", "specify where to write the HTML files to")
	flag.StringVar(&datasetDir, "dataset", "", "specify where to read the JSON files from")
	flag.StringVar(&repoNo, "repo-no", "2", "specify a repository number to use, default is 2")
	flag.StringVar(&templateDir, "templates", "", "specify where to read the templates from")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, "CAIT", cait.Version)
	cfg.LicenseText = fmt.Sprintf(cait.LicenseText, appName, cait.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.OptionText = "OPTIONS\n\n"

	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	datasetDir = cfg.CheckOption("dataset", cfg.MergeEnv("dataset", datasetDir), true)
	repoNo = cfg.CheckOption("repo-no", cfg.MergeEnv("repo_no", repoNo), true)
	templateDir = cfg.CheckOption("templates", cfg.MergeEnv("templates", templateDir), true)
	htdocsDir = cfg.CheckOption("htdocs", cfg.MergeEnv("htdocs", htdocsDir), true)

	if htdocsDir != "" {
		if _, err := os.Stat(htdocsDir); os.IsNotExist(err) {
			os.MkdirAll(htdocsDir, 0775)
		}
	}

	// create our API object
	api := cait.New("", "", "", datasetDir)

	//
	// Setup directories relationships
	//
	accessionsDir := path.Join("repositories", repoNo, "accessions")
	digitalObjectDir := path.Join("repositories", repoNo, "digital_objects")
	subjectDir := path.Join("subjects")
	agentsPeopleDir := path.Join("agents", "people")

	log.Printf("%s %s\n", appName, cait.Version)

	//
	// Setup Maps and generate the accessions pages
	//
	log.Printf("Reading Subjects from %s\n", subjectDir)
	subjectsMap, err := api.MakeSubjectMap(subjectDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Mapped %d subjects\n", len(subjectsMap))

	log.Printf("Reading Digital Objects from %s\n", digitalObjectDir)
	digitalObjectsMap, err := api.MakeDigitalObjectMap(digitalObjectDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Mapped %d Digital Objects\n", len(digitalObjectsMap))

	log.Printf("Reading Agents/People from %s\n", agentsPeopleDir)
	agentsList, err := api.MakeAgentList(agentsPeopleDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Mapped %d Agents/People\n", len(agentsList))

	log.Printf("Processing Agents/People in %s\n", agentsPeopleDir)
	cnt, err := processAgentsPeople(api, templateDir, "agents-people.html", "agents-people.include", agentsPeopleDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Processed %d Agents/Peoples\n", cnt)

	log.Printf("Processing accessions in %s\n", datasetDir)
	cnt, err = processAccessions(api, templateDir, "accession.html", "accession.include", accessionsDir, agentsList, subjectsMap, digitalObjectsMap)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Processed %d Accessoins\n", cnt)
}
