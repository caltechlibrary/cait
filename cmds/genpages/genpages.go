//
// cmds/genpages/genpages.go - A command line utility that builds pages from the exported results of cait.go
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
	"path/filepath"
	"strings"
	//"html/template"
	"text/template"

	// Caltech Library packages
	"github.com/caltechlibrary/cait"
)

var (
	description = `
 USAGE: genpages [OPTIONS]

 OVERVIEW

	genpages generates HTML, .include pages and normalized JSON based on the JSON output form
	cait and templates associated with the command.

 OPTIONS
`
	configuration = `
 CONFIGURATION

    genpages can be configured through setting the following environment
	variables-

    CAIT_DATASET    this is the directory that contains the output of the
                      'cait archivesspace export' command.

    CAIT_TEMPLATES  this is the directory that contains the templates
                      used used to generate the static content of the website.

    CAIT_HTDOCS     this is the directory where the HTML files are written.

`

	help        bool
	htdocsDir   string
	datasetDir  string
	templateDir string
)

func usage() {
	fmt.Println(description)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Println(configuration)
	os.Exit(0)
}

func loadTemplates(templateDir, aHTMLTmplName, aIncTmplName string) (*template.Template, *template.Template, error) {
	aHTMLTmpl, err := cait.AssembleTemplate(path.Join(templateDir, aHTMLTmplName), path.Join(templateDir, aIncTmplName))
	if err != nil {
		return nil, nil, fmt.Errorf("Can't parse template %s, %s, %s", aHTMLTmplName, aIncTmplName, err)
	}
	aIncTmpl, err := cait.Template(path.Join(templateDir, aIncTmplName))
	if err != nil {
		return aHTMLTmpl, nil, fmt.Errorf("Can't parse template %s, %s", aIncTmplName, err)
	}
	return aHTMLTmpl, aIncTmpl, nil
}

func processAccessions(templateDir string, aHTMLTmplName string, aIncTmplName string, agents []*cait.Agent, subjects map[string]*cait.Subject, digitalObjects map[string]*cait.DigitalObject) error {
	log.Printf("Reading templates from %s\n", templateDir)
	aHTMLTmpl, aIncTmpl, err := loadTemplates(templateDir, aHTMLTmplName, aIncTmplName)
	check(err)

	return filepath.Walk(path.Join(datasetDir, "repositories"), func(p string, f os.FileInfo, err error) error {
		// Process accession records
		if strings.Contains(p, "accessions") == true && strings.HasSuffix(p, ".json") == true {
			src, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			accession := new(cait.Accession)
			err = json.Unmarshal(src, &accession)
			if err != nil {
				return err
			}
			// FIXME: which restrictions do we care about--
			//        accession.Publish, accession.Suppressed, accession.AccessRestrictions,
			//        accession.RestrictionsApply, accession.UseRestrictions
			if accession.Publish == true && accession.Suppressed == false && accession.RestrictionsApply == false {
				// Create a normalized view of the accession to make it easier to work with
				view, err := accession.NormalizeView(agents, subjects, digitalObjects)
				if err != nil {
					return fmt.Errorf("Could not generate normalized view, %s", err)
				}

				// If the accession is published and the accession is not suppressed then generate the webpage
				fname := path.Join(htdocsDir, fmt.Sprintf("%s.html", accession.URI))
				dname := path.Dir(fname)
				err = os.MkdirAll(dname, 0775)
				if err != nil {
					return fmt.Errorf("Can't create %s, %s", dname, err)
				}
				// Process HTML file
				fp, err := os.Create(fname)
				if err != nil {
					return fmt.Errorf("Problem creating %s, %s", fname, err)
				}
				log.Printf("Writing %s", fname)
				err = aHTMLTmpl.Execute(fp, view)
				if err != nil {
					log.Fatalf("template execute error %s, %s", aHTMLTmplName, err)
					return err
				}
				fp.Close()

				// Process Include file (just the HTML content)
				fname = path.Join(htdocsDir, fmt.Sprintf("%s.include", accession.URI))
				fp, err = os.Create(fname)
				if err != nil {
					return fmt.Errorf("Problem creating %s, %s", fname, err)
				}
				log.Printf("Writing %s", fname)
				err = aIncTmpl.Execute(fp, view)
				if err != nil {
					log.Fatalf("template execute error %s, %s", aIncTmplName, err)
					return err
				}
				fp.Close()

				// Process JSON file (an abridged version of the JSON output in data)
				fname = path.Join(htdocsDir, fmt.Sprintf("%s.json", accession.URI))
				src, err := json.Marshal(view)
				if err != nil {
					return fmt.Errorf("Could not JSON encode %s, %s", fname, err)
				}
				log.Printf("Writing %s", fname)
				err = ioutil.WriteFile(fname, src, 0664)
				if err != nil {
					log.Fatalf("could not write JSON view %s, %s", fname, err)
					return err
				}
				fp.Close()
			}
		}
		return nil
	})
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getenv(envar, defaultValue string) string {
	tmp := os.Getenv(envar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func init() {
	datasetDir = getenv("CAIT_DATASET", "dataset")
	templateDir = getenv("CAIT_TEMPLATES", path.Join("templates", "default"))
	htdocsDir = getenv("CAIT_HTDOCS", "htdocs")
	flag.StringVar(&htdocsDir, "htdocs", htdocsDir, "specify where to write the HTML files to")
	flag.StringVar(&datasetDir, "dataset", datasetDir, "specify where to read the JSON files from")
	flag.StringVar(&templateDir, "templates", templateDir, "specify where to read the templates from")
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
}

func main() {
	flag.Parse()
	if help == true {
		usage()
	}

	if htdocsDir != "" {
		if _, err := os.Stat(htdocsDir); os.IsNotExist(err) {
			os.MkdirAll(htdocsDir, 0775)
		}
	}

	//
	// Setup directories relationships
	//
	digitalObjectDir := ""
	filepath.Walk(datasetDir, func(p string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() == true && strings.HasSuffix(p, "digital_objects") {
			digitalObjectDir = p
			return nil
		}
		return err
	})
	if digitalObjectDir == "" {
		check(fmt.Errorf("Can't find the digital object directory in %s", datasetDir))
	}
	subjectDir := path.Join(datasetDir, "subjects")
	agentsDir := path.Join(datasetDir, "agents", "people")

	//
	// Setup Maps and generate the accessions pages
	//
	log.Printf("Reading Subjects from %s", subjectDir)
	subjectsMap, err := cait.MakeSubjectMap(subjectDir)
	check(err)

	log.Printf("Reading Digital Objects from %s", digitalObjectDir)
	digitalObjectsMap, err := cait.MakeDigitalObjectMap(digitalObjectDir)
	check(err)

	agentsList, err := cait.MakeAgentList(agentsDir)
	check(err)

	log.Printf("Processing accessions in %s\n", datasetDir)
	err = processAccessions(templateDir, "accession.html", "accession.include", agentsList, subjectsMap, digitalObjectsMap)
	check(err)
}
