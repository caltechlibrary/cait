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
	description = `
 USAGE: searchdataset [OPTIONS]

 SYNOPSIS

 searchdataset is a command line utility to search the dataset directory.
 It produces a Bleve search index used by servepages web service.
 Configuration is done through environmental variables.

 OPTIONS
`

	configuration = `

 CONFIGURATION

 searchdataset relies on the following environment variables for
 configuration when overriding the defaults:

    CAIT_DATASET_INDEX	This is the directory that will contain all the Bleve
                        indexes.

`
	help      bool
	explain   bool
	size      int
	from      int
	indexName string
	// q match
	// q_required match all
	// q_exact match phrase
	// q_excluded disjunct with match
	q         string
	qRequired string
	qExact    string
	qExcluded string
	qAll      bool
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func init() {
	from = 0
	indexName = "dataset.bleve"
	size = 10
	flag.StringVar(&indexName, "-i", indexName, "use this index")
	flag.BoolVar(&help, "h", false, "display this message")
	flag.BoolVar(&explain, "e", false, "explain the query")
	flag.IntVar(&size, "s", size, "display n results for per response")
	flag.IntVar(&from, "f", from, "display results from number")
	flag.StringVar(&q, "q", q, "use query string query")
	flag.StringVar(&qRequired, "q_required", qRequired, "use match term query")
	flag.StringVar(&qExact, "q_exact", qExact, "use match phrase query")
	flag.StringVar(&qExcluded, "q_exclude", qExcluded, "use disjunct query")
	flag.BoolVar(&qAll, "q_all", false, "use match all query")
}

func main() {
	flag.Parse()

	if help == true {
		usage()
	}

	args := flag.Args()
	if len(args) == 0 && qAll == false && q == "" && qRequired == "" && qExact == "" && qExcluded == "" {
		fmt.Fprintln(os.Stderr, "USAGE: search [-h, OPTIONS] QUERY_TERMS")
		os.Exit(1)
	}

	index, err := bleve.Open(indexName)
	if err != nil {
		log.Fatal(err)
	}

	var (
		conQry []bleve.Query
	)
	if len(args) > 0 {
		conQry = append(conQry, bleve.NewQueryStringQuery(strings.Join(args, " ")))
	}
	if q != "" {
		conQry = append(conQry, bleve.NewQueryStringQuery(q))
	}
	if qExact != "" {
		conQry = append(conQry, bleve.NewMatchPhraseQuery(qExact))
	}
	if qRequired != "" {
		for _, s := range strings.Fields(qRequired) {
			conQry = append(conQry, bleve.NewQueryStringQuery(fmt.Sprintf("+%s", s)))
		}
	}
	if qExcluded != "" {
		for _, s := range strings.Fields(qExcluded) {
			conQry = append(conQry, bleve.NewQueryStringQuery(fmt.Sprintf("-%s", s)))
		}
	}
	if qAll == true {
		conQry = append(conQry, bleve.NewMatchAllQuery())
	}

	query := bleve.NewConjunctionQuery(conQry)
	search := bleve.NewSearchRequestOptions(query, size, from, explain)

	search.Highlight = bleve.NewHighlight()
	search.Highlight.AddField("title")
	search.Highlight.AddField("content_description")
	search.Highlight.AddField("subjects")
	search.Highlight.AddField("extents")
	search.Highlight.AddField("digital_objects.title")
	search.Highlight.AddField("digital_objects.files_uris")

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
	fmt.Printf("Available fields: %s\n", strings.Join(fields, "|"))
}
