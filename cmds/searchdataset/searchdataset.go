//
// cmds/searchdataset/searchdataset.go - a command line utility to search the dataset indexed with indexdataset utility.
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
