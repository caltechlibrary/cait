//
// caitindexer.go - A search indexer for [Bleve Search](https://github.com/blevesearch/bleve)
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

	"github.com/blevesearch/bleve"

	"../../../cait"
)

var (
	description = `
 USAGE: caitindexer [-h|--help]

 SYNOPSIS

 caitindexer is a command line utility to index content fetched from
 an ArchivesSpace via the ArchivesSpace REST API (e.g. with
 cait tool). It indexes content for the Bleve search library.
 Configuration is done through environmental variables.

 OPTIONS
`

	configuration = `

 CONFIGURATION

 caitindexer relies on the following environment variables for
 configuration when overriding the defaults:

    CAIT_HTDOCS       This should be the path to the directory tree
                        containings the content (e.g. JSON files) to be index.
                        This is generally populated with the caitpage command.
						Defaults to ./htdocs.

    CAIT_BLEVE_INDEX	This is the directory that will contain all the Bleve
                        indexes. Defaults to ./index.bleve

`
	help      bool
	htdocsDir string
	indexName string
	dirCount  int
	fileCount int
	index     bleve.Index
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func init() {
	htdocsDir = "htdocs"
	indexName = "index.bleve"
	htdocsDir = os.Getenv("CAIT_HTDOCS")
	indexName = os.Getenv("CAIT_BLEVE_INDEX")
	flag.StringVar(&htdocsDir, "htdocs", "htdocs", "The document root for the website")
	flag.StringVar(&indexName, "index", "index.bleve", "The name of the Bleve index")
	flag.BoolVar(&help, "h", false, "this help message")
	flag.BoolVar(&help, "help", false, "this help message")
}

func getIndex(indexName string) (bleve.Index, error) {
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		log.Printf("Creating Bleve index at %s\n", indexName)
		mapping := bleve.NewIndexMapping()
		mapping.DefaultAnalyzer = "en"
		index, err := bleve.New(indexName, mapping)
		if err != nil {
			return nil, fmt.Errorf("Can't create new bleve index %s, %s", indexName, err)
		}
		return index, nil
	}
	log.Printf("Opening Bleve index at %s\n", indexName)
	index, err := bleve.Open(indexName)
	if err != nil {
		return nil, fmt.Errorf("Can't create new bleve index %s, %s", indexName, err)
	}
	return index, nil
}

func walkHtdocs(p string, f os.FileInfo, err error) error {
	if strings.Contains(p, "/accessions/") == true && strings.HasSuffix(p, ".json") == true {
		src, err := ioutil.ReadFile(p)
		if err != nil {
			log.Printf("Can't read %s, %s", p, err)
			return nil
		}
		view := new(cait.NormalizedAccessionView)
		err = json.Unmarshal(src, &view)
		if err != nil {
			log.Printf("Can't parse %s, %s", p, err)
			return nil
		}
		// Trim the htdocsDir and trailing .json extension
		i := len(htdocsDir)
		j := len(p)
		err = index.Index(p[i:j-5], view)
		if err != nil {
			log.Printf("Indexing error %s, %s", p, err)
			return nil
		}
		log.Printf("Indexed %s", p)
	}
	return nil
}

func indexSite() error {
	//FIXME: Need to switch this for indexing on batch.
	return filepath.Walk(path.Join(htdocsDir, "repositories"), walkHtdocs)
}

func main() {
	var err error

	flag.Parse()
	if help == true {
		usage()
	}

	index, err = getIndex(indexName)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer index.Close()

	// Walk our data import tree and index things
	log.Printf("Start indexing of %s in %s\n", htdocsDir, indexName)
	indexSite()
	log.Printf("Finished")
}
