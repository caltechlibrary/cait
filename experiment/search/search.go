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
	indexName := os.Getenv("CAIT_BLEVE_INDEX")
	index, err := bleve.Open(indexName)
	if err != nil {
		log.Fatal(err)
	}
	terms := strings.Join(args, " ")
	query := bleve.NewQueryStringQuery(terms)
	if from < 1 {
		from = 1
	}
	if size < 1 {
		size = 10
	}
	search := bleve.NewSearchRequestOptions(query, size, from-1, explain)
	subjectFacet := bleve.NewFacetRequest("subjects", 5)
	search.AddFacet("subjects", subjectFacet)

	search.Highlight = bleve.NewHighlight()
	//search.Highlight = bleve.NewHighlightWithStyle("ansi")
	search.Highlight.AddField("title")
	search.Highlight.AddField("content_description")
	search.Highlight.AddField("extents")
	search.Highlight.AddField("digital_objects.title")
	search.Highlight.AddField("digital_objects.files_uris")

	results, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
	fields, _ := index.Fields()
	fmt.Printf("DEBUG fields: %s\n", strings.Join(fields, "|"))
}
