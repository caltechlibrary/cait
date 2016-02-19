package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/bleve"
)

var (
	help    bool
	explain bool
	size    int
	from    int
)

func init() {
	size = 10
	flag.BoolVar(&help, "h", true, "display this message")
	flag.BoolVar(&explain, "e", true, "explain the query")
	flag.IntVar(&size, "s", 10, "display n results for per response")
	flag.IntVar(&from, "f", 0, "display results from number")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "USAGE: search QUERY_TERMS")
		flag.PrintDefaults()
		os.Exit(1)
	}
	indexName := os.Getenv("CAIT_HTDOCS_INDEX")
	index, err := bleve.Open(indexName)
	if err != nil {
		log.Fatal(err)
	}
	terms := strings.Join(args, " ")
	query := bleve.NewQueryStringQuery(terms)
	if from < 0 {
		from = 0
	}
	if size < 0 {
		size = 10
	}
	search := bleve.NewSearchRequestOptions(query, size, from, explain)

	search.Highlight = bleve.NewHighlight()
	//search.Highlight = bleve.NewHighlightWithStyle("ansi")
	search.Highlight.AddField("title")
	search.Highlight.AddField("content_description")
	search.Highlight.AddField("subjects")
	search.Highlight.AddField("extents")
	//	search.Highlight.AddField("digital_objects.title")
	//	search.Highlight.AddField("digital_objects.files_uris")

	subjectFacet := bleve.NewFacetRequest("subjects", 3)
	search.AddFacet("subjects", subjectFacet)
	// Return all fields for each result.
	search.Fields = []string{"*"}

	results, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
	fields, _ := index.Fields()
	fmt.Printf("DEBUG fields: %s\n", strings.Join(fields, "|"))
}
