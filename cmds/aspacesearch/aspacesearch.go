/**
 * cmds/aspacesearch/aspacesearch.go - A command line utility runs search
 * for the site defined by ASPACE_HTDOCS using the index identified by
 * ASPACE_BLEVE_INDEX.
 */
package main

import (
	"flag"
	"fmt"
	"os"

	//"../../../aspace"
)

var (
	description = `
 USAGE: aspacesearch [OPTIONS]

 OVERVIEW

	aspacesearch provides search services defined by ASPACE_SEARCH_URL for the
	website content defined by ASPACE_HTDOCS using the index defined
	by ASPACE_BLEVE_INDEX.

 OPTIONS
`
	configuration = `
 CONFIGURATION

 aspacesearch can be configured through environment variables. The following
 variables are supported-

   ASPACE_SEARCH_URL

   ASPACE_BLEVE_INDEX

   ASPACE_TEMPLATES

`
	help         bool
	indexName    string
	htdocsDir    string
	templatesDir string
	serviceURL   string
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func init() {
	serviceURL = os.Getenv("ASPACE_SEARCH_URL")
	indexName = os.Getenv("ASPACE_BLEVE_INDEX")
	htdocsDir = os.Getenv("ASPACE_HTDOCS")
	templatesDir = os.Getenv("ASPACE_TEMPLATES")
	flag.StringVar(&serviceURL, "search", "http://localhost:8501", "The URL to listen on for search requests")
	flag.StringVar(&indexName, "index", "index.bleve", "specify the Bleve index to use")
	flag.StringVar(&htdocsDir, "htdocs", "htdocs", "specify where to write the HTML files to")
	flag.StringVar(&templatesDir, "templates", "templates/default", "The directory path for templates")
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
}

func main() {
	flag.Parse()
	if help == true {
		usage()
	}

	fmt.Printf("aspacesearch not implemented yet: %s, %s, %s, %s", serviceURL, indexName, htdocsDir, templatesDir)
}
