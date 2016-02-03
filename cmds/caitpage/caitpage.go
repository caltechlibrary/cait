//
// cmds/caitpage/caitpage.go - A command line utility that builds pages from the exported results of cait.go
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
	"text/template"

	"../../../cait"
)

var (
	description = `
 USAGE: caitpage [OPTIONS]

 OVERVIEW

	caitpage generates HTML pages based on the JSON output form
	cait and templates associated with the command.

 OPTIONS
`
	configuration = `
 CONFIGURATION

    caitpages can be configured through setting the following environment
	variables-

    CAIT_DATASET    this is the directory that contains the output of the
                      'cait archivesspace export' command. Defaults to ./data

    CAIT_TEMPLATES  this is the directory that contains the templates
                      used used to generate the static content of the website.
                      Defaults to ./templates/default.

    CAIT_HTDOCS     this is the directory where the HTML files are written.
                      Defaults to ./htdocs

`

	help          bool
	htdocsDir     string
	dataDir       string
	templateDir   string
	aHTMLTmplName = "accession.html"
	aHTMLTmpl     = template.New(aHTMLTmplName)
	aIncTmplName  = "accession.include"
	aIncTmpl      = template.New(aIncTmplName)

	subjects = make(map[string]*cait.Subject)
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func loadTemplates(templateDir string) error {
	var err error
	aHTMLTmpl, err = template.ParseFiles(path.Join(templateDir, aHTMLTmplName), path.Join(templateDir, aIncTmplName))
	if err != nil {
		return fmt.Errorf("Can't parse template %s, %s, %s", aHTMLTmplName, aIncTmplName, err)
	}
	aIncTmpl, err = template.ParseFiles(path.Join(templateDir, aIncTmplName))
	if err != nil {
		return fmt.Errorf("Can't parse template %s, %s", aIncTmplName, err)
	}
	return nil
}

func loadSubjects(subjectDir string) error {
	var err error
	subjects, err = cait.MakeSubjectMap(subjectDir)
	return err
}

func processData(titleIndex map[string]*cait.NavRecord) error {
	return filepath.Walk(path.Join(dataDir, "repositories"), func(p string, f os.FileInfo, err error) error {
		// Process accession records
		if strings.Contains(p, "accessions") == true && strings.HasSuffix(p, ".json") {
			src, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			accession := new(cait.Accession)
			err = json.Unmarshal(src, &accession)
			if err != nil {
				return err
			}
			if accession.Publish == true && accession.Suppressed == false {
				// Create a normalized view of the accession to make it easier to work with
				navRecord, _ := titleIndex[accession.URI]
				view, err := accession.NormalizeView(subjects, navRecord)
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
					log.Fatalf("template execute error %s, %s", "accession.html", err)
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
					log.Fatalf("template execute error %s, %s", "accession.include", err)
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

func init() {
	dataDir = os.Getenv("CAIT_DATASET")
	templateDir = os.Getenv("CAIT_TEMPLATES")
	htdocsDir = os.Getenv("CAIT_HTDOCS")
	flag.StringVar(&htdocsDir, "htdocs", "htdocs", "specify where to write the HTML files to")
	flag.StringVar(&dataDir, "data", "data", "specify where to read the JSON files from")
	flag.StringVar(&templateDir, "templates", path.Join("templates", "default"), "specify where to read the templates from")
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
}

func main() {
	flag.Parse()
	if help == true {
		usage()
	}
	titleIndex, err := cait.MakeAccessionTitleIndex(dataDir)
	if err != nil {
		log.Fatalf("Can't make a title index %s, %s", dataDir, err)
	}
	subjectDir := path.Join(dataDir, "subjects")
	log.Printf("Reading Subjects from %s", subjectDir)
	err = loadSubjects(subjectDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Reading templates from %s\n", templateDir)
	err = loadTemplates(templateDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Processing data in %s\n", dataDir)
	err = processData(titleIndex)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
