/**
 * cmds/aspacepage/aspacepage.go - A command line utility that builds pages from the results of aspace.go
 * command.
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
	"path/filepath"
	"strings"
	"text/template"

	"../../../aspace"
)

var (
	description = `
 USAGE: aspacepage [OPTIONS]

 OVERVIEW

	aspacepage generates HTML pages based on the JSON output form
	aspace and templates associated with the command.

 OPTIONS
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
)

func usage(exitCode int) {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println("")
	os.Exit(exitCode)
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
	fname = path.Join(templateDir, aJSONTmplName)
	aJSONTmpl, err = aJSONTmpl.ParseFiles(fname)
	if err != nil {
		return fmt.Errorf("Can't parse template %s, %s", fname, err)
	}
	return nil
}

func walkRepositories(p string, f os.FileInfo, err error) error {
	// Process accession records
	if strings.Contains(p, "accessions") == true && strings.HasSuffix(p, ".json") {
		src, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		accession := new(aspace.Accession)
		err = json.Unmarshal(src, &accession)
		if err != nil {
			return err
		}
		if accession.Publish == true && accession.Suppressed == false {
			// If the accession is published and the accession is not suppressed then generate the webpage
			fname := path.Join(htdocsDir, fmt.Sprintf("%s.html", accession.URI))
			dname := path.Dir(fname)
			err := os.MkdirAll(dname, 0775)
			if err != nil {
				return fmt.Errorf("Can't create %s, %s", dname, err)
			}
			// Process HTML file
			fp, err := os.Create(fname)
			if err != nil {
				return fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			log.Printf("Writing %s", fname)
			err = aHTMLTmpl.Execute(fp, accession)
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
			err = aIncTmpl.Execute(fp, accession)
			if err != nil {
				log.Fatalf("template execute error %s, %s", "accession.html", err)
				return err
			}
			fp.Close()

			// Process JSON file (an abridged version of the JSON output in data)
			fname = path.Join(htdocsDir, fmt.Sprintf("%s.json", accession.URI))
			fp, err = os.Create(fname)
			if err != nil {
				return fmt.Errorf("Problem creating %s, %s", fname, err)
			}
			log.Printf("Writing %s", fname)
			err = aJSONTmpl.Execute(fp, accession)
			if err != nil {
				log.Fatalf("template execute error %s, %s", "accession.html", err)
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
	flag.StringVar(&htdocsDir, "htdocs", "htdocs", "specify where to write the HTML files to")
	flag.StringVar(&dataDir, "data", "data", "specify where to read the JSON files from")
	flag.StringVar(&templateDir, "templates", path.Join("templates", "default"), "specify where to read the templates from")
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
}

func main() {
	flag.Parse()
	if help == true {
		usage(0)
	}

	log.Printf("Reading templates from %s\n", templateDir)
	err := loadTemplates(templateDir)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Processing data in %s\n", dataDir)
	err = processData()
	if err != nil {
		log.Fatalf("%s", err)
	}
}
