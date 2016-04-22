//
// cmds/indexdataset/indexdataset.go - A command line utility that builds a bleve index of the raw contents exported with cait utility
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
	"time"

	// 3rd Party Packages
	"github.com/blevesearch/bleve"

	// Caltech Library Packages
	"github.com/caltechlibrary/cait"
)

var (
	description = `
 USAGE: indexdataset [OPTIONS]

 SYNOPSIS

 indexdataset is a command line utility to indexes content in the dataset directory.
 It produces a Bleve search index used by servepages web service.
 Configuration is done through environmental variables.

 OPTIONS
`

	configuration = `

 CONFIGURATION

 indexdataset relies on the following environment variables for
 configuration when overriding the defaults:

    CAIT_DATASET       This should be the path to the directory tree
                       containings the imported content (e.g. JSON files) to be index.
                       This is generally populated by the cait command.

    CAIT_DATASET_INDEX	This is the directory that will contain all the Bleve
                        indexes.

`
	help         bool
	replaceIndex bool
	indexName    string
	datasetDir   string
)

func usage() {
	fmt.Println(description)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t(defaults to %s) %s\n", f.Name, f.DefValue, f.Usage)
	})
	fmt.Println(configuration)
	os.Exit(0)
}

func openIndex(indexName string, indexMapping *bleve.IndexMapping) (bleve.Index, error) {
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		return bleve.New(indexName, indexMapping)
	}
	return bleve.Open(indexName)
}

func indexSite(index bleve.Index, batchSize int, dataSet map[string]interface{}) error {
	batch := index.NewBatch()
	startT := time.Now()
	i := 0
	for id, data := range dataSet {
		batch.Index(id, data)
		if batch.Size() >= batchSize {
			err := index.Batch(batch)
			if err != nil {
				return err
			}
			i += batch.Size()
			log.Printf("Index %d accessions", i)
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		err := index.Batch(batch)
		i += batch.Size()
		log.Printf("Index %d accessions in %s", i, time.Now().Sub(startT))
		return err
	}
	return nil
}

func getenv(envvar, defaultValue string) string {
	tmp := os.Getenv(envvar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func init() {
	datasetDir = getenv("CAIT_DATASET", "dataset")
	indexName = getenv("CAIT_DATASET_INDEX", "dataset.bleve")
	flag.StringVar(&datasetDir, "dataset", datasetDir, "The document root for the dataset")
	flag.StringVar(&indexName, "index", indexName, "The name of the Bleve index")
	flag.BoolVar(&replaceIndex, "r", false, "Replace the index if it exists")
	flag.BoolVar(&help, "h", false, "this help message")
	flag.BoolVar(&help, "help", false, "this help message")
}

func main() {
	flag.Parse()

	if help == true {
		usage()
		os.Exit(0)
	}
	if replaceIndex == true {
		os.RemoveAll(indexName)
	}

	log.Println("Building subject map...")
	subjectMap, _ := cait.MakeSubjectMap(path.Join(datasetDir, "subjects"))
	log.Println("Building digital object map...")
	digitalObjectMap, _ := cait.MakeDigitalObjectMap(path.Join(datasetDir, "repositories/2/digital_objects"))

	log.Println("Building agent list...")
	agentList, _ := cait.MakeAgentList(path.Join(datasetDir, "people"))

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

	// objectTitleMapping := bleve.NewTextFieldMapping()
	// // objectTitleMapping.Analyzer = "en"
	// objectTitleMapping.Store = true
	// objectTitleMapping.Index = false
	// objectTitleMapping.IncludeInAll = false
	// objectTitleMapping.IncludeTermVectors = false
	// accessionMapping.AddFieldMappingsAt("digital_objects.title", objectTitleMapping)
	//
	// objectFileURIMapping := bleve.NewTextFieldMapping()
	// //objectFileURIMapping.Analyzer = "Simple"
	// objectFileURIMapping.Store = true
	// objectFileURIMapping.Index = false
	// objectFileURIMapping.IncludeInAll = false
	// objectFileURIMapping.IncludeTermVectors = false
	// accessionMapping.AddFieldMappingsAt("digital_objects.file_uris", objectFileURIMapping)
	//
	// extentsMapping := bleve.NewTextFieldMapping()
	// extentsMapping.Analyzer = "en"
	// extentsMapping.Store = true
	// extentsMapping.Index = true
	// accessionMapping.AddFieldMappingsAt("extents", extentsMapping)
	//
	// createdMapping := bleve.NewDateTimeFieldMapping()
	// accessionMapping.AddFieldMappingsAt("created", createdMapping)

	// Finally add this mapping to the main index mapping
	indexMapping.AddDocumentMapping("accession", accessionMapping)

	index, _ := openIndex(indexName, indexMapping)
	log.Println("Start indexing...")
	startT := time.Now()
	indexSite(index, 1000, (func() map[string]interface{} {
		i := 0
		m := make(map[string]interface{})
		filepath.Walk(path.Join(datasetDir, "repositories/2/accessions"), func(p string, _ os.FileInfo, _ error) error {
			if strings.HasSuffix(p, ".json") {
				if (i % 100) == 0 {
					log.Printf("Read %d accessions", i)
				}
				src, _ := ioutil.ReadFile(p)
				data := new(cait.Accession)
				err := json.Unmarshal(src, &data)
				if err == nil {
					i++
					m[data.URI], _ = data.NormalizeView(agentList, subjectMap, digitalObjectMap)
				}
				return err
			}
			return nil
		})
		log.Printf("Read %d accessions in %s", i, time.Now().Sub(startT))
		return m
	})())
	log.Printf("Done! %s", time.Now().Sub(startT))
}
