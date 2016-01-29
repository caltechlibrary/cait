//
// cmds/caitpage/caitpage.go - A command line utility that builds pages from the exported results of cait.go
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2016
// Caltech Library
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
                      'cait instance export' command. Defaults to ./data

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
	aJSONTmplName = "accession.json"
	aJSONTmpl     = template.New(aJSONTmplName)
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
	fname := path.Join(templateDir, aHTMLTmplName)
	aHTMLTmpl, err = aHTMLTmpl.ParseFiles(fname)
	if err != nil {
		return fmt.Errorf("Can't parse template %s, %s", fname, err)
	}
	fname = path.Join(templateDir, aIncTmplName)
	aIncTmpl, err = aIncTmpl.ParseFiles(fname)
	if err != nil {
		return fmt.Errorf("Can't parse template %s, %s", fname, err)
	}
	return nil
}

func loadSubjects(subjectDir string) error {
	var err error
	subjects, err = cait.MakeSubjectMap(subjectDir)
	return err
}

func walkRepositories(p string, f os.FileInfo, err error) error {
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
			view, err := accession.NormalizeView(subjects)
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
}

func processData() error {
	return filepath.Walk(path.Join(dataDir, "repositories"), walkRepositories)
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

	subjectDir := path.Join(dataDir, "subjects")
	log.Printf("Reading Subjects from %s", subjectDir)
	err := loadSubjects(subjectDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Reading templates from %s\n", templateDir)
	err = loadTemplates(templateDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Processing data in %s\n", dataDir)
	err = processData()
	if err != nil {
		log.Fatalf("%s", err)
	}
}
