//
// cmds/indexpages/indexpages.go - Create/update a bleve index the htdocs contents generated with the genpages utility.
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
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	// 3rd Party packages
	"github.com/blevesearch/bleve"

	// Caltech Libraries packages
	"github.com/caltechlibrary/cait"
	"github.com/caltechlibrary/cli"
)

var (
	usage = `USAGE: %s [OPTIONS] [BLEVE_INDEX]`

	description = `
SYNOPSIS

%s is a command line utility to create a new index of content in the htdocs directory.
It produces a Bleve search index used by servepages web service.
Configuration is done through environmental variables.

CONFIGURATION

%s relies on the following environment variables for
configuration when overriding the defaults:

    CAIT_HTDOCS   This should be the path to the directory tree
                  containings the content (e.g. JSON files) to be index.
                  This is generally populated with the genpages command.

    CAIT_BLEVE	  A colon delimited list of the Bleve indexes (for swapping)
`

	// Standard Options
	showHelp    bool
	showVersion bool
	showLicense bool

	// App Options
	htdocs      string
	bleveNames  string
	dirCount    int
	fileCount   int
	showVerbose bool
)

func handleSignals() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			//handle SIGINT
			log.Println("SIGINT received shutting down")
			os.Exit(0)
		case syscall.SIGTERM:
			//handle SIGTERM
			log.Println("SIGTERM received shutting down")
			os.Exit(0)
		}
	}()
}

func newIndex(indexName string) (bleve.Index, error) {
	log.Printf("Creating Bleve index at %s\n", indexName)
	log.Println("Setting up index...")
	indexMapping := bleve.NewIndexMapping()
	// Add Accession as a specific document map
	accessionMapping := bleve.NewDocumentMapping()

	// Now add specific accession fields
	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = "en"
	titleMapping.Store = true
	titleMapping.Index = true
	accessionMapping.AddFieldMappingsAt("title", titleMapping)

	descriptionMapping := bleve.NewTextFieldMapping()
	descriptionMapping.Analyzer = "en"
	descriptionMapping.Store = true
	descriptionMapping.Index = true
	accessionMapping.AddFieldMappingsAt("content_description", descriptionMapping)

	subjectsMapping := bleve.NewTextFieldMapping()
	subjectsMapping.Analyzer = "en"
	subjectsMapping.Store = true
	subjectsMapping.Index = true
	subjectsMapping.IncludeTermVectors = true
	accessionMapping.AddFieldMappingsAt("subjects", subjectsMapping)

	subjectsFunctionMapping := bleve.NewTextFieldMapping()
	subjectsFunctionMapping.Analyzer = "en"
	subjectsFunctionMapping.Store = true
	subjectsFunctionMapping.Index = true
	subjectsFunctionMapping.IncludeTermVectors = true
	accessionMapping.AddFieldMappingsAt("subjects_function", subjectsFunctionMapping)

	subjectsTopicalMapping := bleve.NewTextFieldMapping()
	subjectsTopicalMapping.Analyzer = "en"
	subjectsTopicalMapping.Store = true
	subjectsTopicalMapping.Index = true
	subjectsTopicalMapping.IncludeTermVectors = true
	accessionMapping.AddFieldMappingsAt("subjects_topical", subjectsTopicalMapping)

	objectTitleMapping := bleve.NewTextFieldMapping()
	objectTitleMapping.Analyzer = "en"
	objectTitleMapping.Store = true
	objectTitleMapping.Index = false
	accessionMapping.AddFieldMappingsAt("digital_objects.title", objectTitleMapping)

	objectFileURIMapping := bleve.NewTextFieldMapping()
	objectFileURIMapping.Analyzer = ""
	objectFileURIMapping.Store = true
	objectFileURIMapping.Index = false
	accessionMapping.AddFieldMappingsAt("digital_objects.file_uris", objectFileURIMapping)

	extentsMapping := bleve.NewTextFieldMapping()
	extentsMapping.Analyzer = "en"
	extentsMapping.Store = true
	extentsMapping.Index = true
	accessionMapping.AddFieldMappingsAt("extents", extentsMapping)

	accessionDateMapping := bleve.NewTextFieldMapping()
	accessionDateMapping.Analyzer = "en"
	accessionDateMapping.Store = true
	accessionDateMapping.Index = false
	accessionMapping.AddFieldMappingsAt("accession_date", accessionDateMapping)

	datesMapping := bleve.NewTextFieldMapping()
	datesMapping.Store = true
	datesMapping.Index = false
	accessionMapping.AddFieldMappingsAt("date_expression", datesMapping)

	createdMapping := bleve.NewDateTimeFieldMapping()
	createdMapping.Store = true
	createdMapping.Index = false
	accessionMapping.AddFieldMappingsAt("created", createdMapping)

	// Finally add this mapping to the main index mapping
	indexMapping.AddDocumentMapping("accession", accessionMapping)

	log.Printf("Opening a new Bleve index at %s", indexName)
	index, err := bleve.New(indexName, indexMapping)
	if err != nil {
		return nil, fmt.Errorf("Can't create new bleve index %s, %s", indexName, err)
	}
	return index, nil
}

func indexSite(index bleve.Index, maxBatchSize int, verbose bool) error {
	startT := time.Now()
	count := 0
	batch := index.NewBatch()
	batchSize := 10
	log.Printf("Walking %s", path.Join(htdocs, "repositories"))
	err := filepath.Walk(path.Join(htdocs, "repositories"), func(p string, f os.FileInfo, err error) error {
		if strings.Contains(p, "/accessions/") == true && strings.HasSuffix(p, ".json") == true {
			src, err := ioutil.ReadFile(p)
			if err != nil {
				log.Printf("Can't read %s, %s", p, err)
				return nil
			}
			if utf8.Valid(src) == false {
				log.Printf("%s not valid UTF-8", p)
				return nil
			}
			view := new(cait.NormalizedAccessionView)
			err = json.Unmarshal(src, &view)
			if err != nil {
				log.Printf("Can't parse %s, %s", p, err)
				return nil
			}
			// Trim the htdocs and trailing .json extension
			//log.Printf("Queued %s", p)
			err = batch.Index(strings.TrimSuffix(strings.TrimPrefix(p, htdocs), "json"), view)
			if err != nil {
				log.Printf("Indexing error %s, %s", p, err)
				return nil
			}
			if batch.Size() >= batchSize {
				err := index.Batch(batch)
				if err != nil {
					log.Fatal(err)
				}
				count += batch.Size()
				batch = index.NewBatch()
				if batchSize < maxBatchSize {
					batchSize = batchSize * 2
				}
				if batchSize > maxBatchSize {
					batchSize = maxBatchSize
				}
			}
			if verbose == true && count > 0 && (count%100) == 0 {
				log.Printf("Indexed: %d items, batch size %d, running %s\n", count, batchSize, time.Now().Sub(startT))
			}
		}
		return nil
	})
	if batch.Size() > 0 {
		err := index.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
		count += batch.Size()
		if verbose == true {
			log.Printf("Indexed: %d items, batch size %d, running %s\n", count, batchSize, time.Now().Sub(startT))
		}
	}
	log.Printf("Total indexed: %d items, total run time %s\n", count, time.Now().Sub(startT))
	return err
}

func check(cfg *cli.Config, key, value string) string {
	if value == "" {
		log.Fatalf("Missing %s_%s", cfg.EnvPrefix, strings.ToUpper(key))
		return ""
	}
	return value
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

	// App Options
	flag.StringVar(&htdocs, "htdocs", "", "The document root for the website")
	flag.StringVar(&bleveNames, "bleve", "", "a colon delimited list of Bleve index db names")
	flag.BoolVar(&showVerbose, "verbose", false, "more verbose logging")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, "CAIT", fmt.Sprintf(cait.LicenseText, appName, cait.Version), cait.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.OptionsText = "OPTIONS\n"

	if showHelp == true {
		fmt.Println(cfg.Usage())
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

	if len(args) > 0 {
		bleveNames = strings.Join(args, ":")
	}

	htdocs = cfg.CheckOption("htdocs", cfg.MergeEnv("htdocs", htdocs), true)
	names := cfg.CheckOption("bleve", cfg.MergeEnv("bleve", bleveNames), true)

	log.Printf("%s %s\n", appName, cait.Version)

	handleSignals()

	for _, indexName := range strings.Split(names, ":") {
		index, err := newIndex(indexName)
		if err != nil {
			log.Printf("Skipping %s, %s", indexName, err)
		} else {
			defer index.Close()

			// Walk our data import tree and index things
			log.Printf("Start indexing of %s in %s\n", htdocs, indexName)
			err = indexSite(index, 1000, showVerbose)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Printf("Finished")
}
