/**
 * aspaceindexer.go - A search indexer for [Bleve Search](https://github.com/blevesearch/bleve)
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

	"github.com/blevesearch/bleve"

	"../../../aspace"
)

var (
	description = `
 USAGE: aspaceindexer [-h|--help]

 SYNOPSIS

 aspaceindexer is a command line utility to index content fetched from
 an ArchivesSpace instance via the ArchivesSpace REST API (e.g. with
 aspace tool). It indexes content for the Bleve search library.
 Configuration is done through environmental variables.

 OPTIONS
`

	configuration = `

 CONFIGURATION

 aspaceindexer relies on the following environment variables for
 configuration when overriding the defaults:

    ASPACE_HTDOCS       This should be the path to the directory tree
                        containings the content (e.g. JSON files) to be index.
                        This is generally populated with the aspacepage command.
						Defaults to ./htdocs.

    ASPACE_BLEVE_INDEX	This is the directory that will contain all the Bleve
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
	htdocsDir = os.Getenv("ASPACE_HTDOCS")
	indexName = os.Getenv("ASPACE_BLEVE_INDEX")
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
		view := new(aspace.NormalizedAccessionView)
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
	log.Printf("Finsihed")
}
