/**
 * cmds/aspacesearch/aspacesearch.go - A command line utility runs search
 * for the site defined by ASPACE_HTDOCS using the index identified by
 * ASPACE_BLEVE_INDEX.
 */
package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"text/template"// Using text template because I am not HTML escaping results...
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"../../../aspace"
	"github.com/blevesearch/bleve"
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
	serviceURL   *url.URL

	advancedPage []byte
	basicPage    []byte

	index bleve.Index
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

// SearchForm holds the expected values for both Basic and Advanced search
type SearchForm struct {
	Method   string `json:"method"`
	Action   string `json:"action"`
	AllIDs   bool   `json:"all_ids,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Page     int    `json:"page,omitempty"`
	// Simple Search
	Query string `json:"q,omitempty"`
	// Advanced Search
	QueryRequired string `json:"q_required,omitempty"`
	QueryPhrase   string `json:"q_phrase,omitempty"`
	QueryExcluded string `json:"q_exclude,omitempty"`
	// Subjects can be a comma delimited list of subjects (e.g. Manuscript Collection, Image Archive)
	Subjects string `json:"q_subjects,omitempty"`
}

// Records are the return structure with all search results and metadata to navigate them
type Records struct {
	Prefix string
	// SearchTerms resolves string to a search expression (Strange Attraction+subjects:Manuscript Collection-chemestry)
	SearchTerms string
	FirstPage   int `json:"first_page,omitempty"`
	LastPage    int `json:"last_page,omitempty"`
	ThisPage    int `json:"this_page,omitempty"`
	OffsetFirst int `json:"offset_first,omitempty"`
	OffsetLast  int `json:"offset_last"`
	TotalHits   int `json:"total_hits,omitempty"`

	//Facets map[string]map[string]interface{} `json:"facets,omitempty"`
	//{"facet_queries":{},"facet_fields":{},"facet_dates":{},"facet_ranges":{}

	Records []*aspace.NormalizedAccessionView `json:"results,omitemty"`
}

func mapToSearchQuery(m map[string]string) (*aspace.SearchQuery, error) {
	var err error
	fmt.Printf("DEBUG starting mapToSearchQuery: %v\n", m)
	q := new(aspace.SearchQuery)
	if _, ok := m["uri"]; ok == true {
		q.URI = m["uri"]
		fmt.Printf("DEBUG q.URI: %s\n", q.URI)
	}
	if _, ok := m["q"]; ok == true {
		q.Q = m["q"]
	}
	if _, ok := m["page"]; ok == true {
		fmt.Printf("DEBUG converting page: %s\n", m["page"])
		q.Page, err = strconv.Atoi(m["page"])
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
	}
	if _, ok := m["page_size"]; ok == true {
		fmt.Printf("DEBUG converting page_size: %s\n", m["page_size"])
		q.PageSize, err = strconv.Atoi(m["page_size"])
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
	}

	/*
		if _, ok := m["repo_id"]; ok == true {
			q.RepoID, err = strconv.Atoi(m["repo_id"])
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}
		}
		if _, ok := m["type"]; ok == true {
			q.Type = m["type"]
		}
		if _, ok := m["sort"]; ok == true {
			q.Sort = m["sort"]
		}
		if _, ok := m["id_set"]; ok == true {
			q.IDSet, err = strconv.Atoi(m["id_set"])
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}
		}
		if _, ok := m["all_ids"]; ok == true {
			q.AllIDs, err = strconv.ParseBool(m["all_ids"])
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}
		}
	*/

	//FIXME: Facets, FilterTerm, SimpleFilter, Exclude... not sure how to form the key/value pairs for GET and POST
	fmt.Printf("DEBUG resolved query submission: %s\n", q)
	return q, nil
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error in POST: %s", err)))
		return
	}

	fmt.Printf("DEBUG r.Form: %v\n", r.Form)
	// Query ArchivesSpace's Solr API or ArchivesSpace's own API
	// Output Results in results template for list or single record as appropriate
	if err != nil {
		log.Printf("API access error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	submission := make(map[string]string)

	// Basic Search results
	if r.Method == "GET" {
		for k, v := range query {
			fmt.Printf("DEBUG k %s v type: %T -> %v\n", k, v, v)
			submission[k] = strings.Join(v, "")
		}
	}

	// Advanced Search results
	if r.Method == "POST" {
		for k, v := range r.Form {
			fmt.Printf("DEBUG k %s v type: %T -> %v\n", k, v, v)
			submission[k] = strings.Join(v, "")
		}
	}

	if _, ok := submission["page"]; ok != true {
		submission["page"] = "1"
	}

	fmt.Printf("DEBUG r.Method: %s\n", r.Method)
	fmt.Printf("DEBUG r.URL.Path: %s\n", r.URL.Path)
	fmt.Printf("DEBUG submission: %v\n", submission)

	//w.Header().Set("Content-Type", "text/html")
	//w.Write([]byte(resultsPage))
	q, err := mapToSearchQuery(submission)
	if err != nil {
		log.Printf("API access error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	fmt.Printf("DEBUG q now: %v\n", q)
	qry := bleve.NewMatchQuery(q.Q)
	search := bleve.NewSearchRequest(qry)
	search.Highlight = bleve.NewHighlightWithStyle("html")
	//search.Explain = true
	searchResults, err := index.Search(search)
	if err != nil {
		log.Printf("Bleve results error %v, %s", qry, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	log.Printf("DEBUG From %d\n", searchResults.Request.From)
	log.Printf("DEBUG PageSize %d\n", searchResults.Request.Size)
	log.Printf("DEBUG Total %d\n", searchResults.Total)
	log.Printf("DEBUG Hits[0].ID %s\n", searchResults.Hits[0].ID)
	log.Printf("DEBUG Hits[0].Fragments %v\n", searchResults.Hits[0].Fragments)
	log.Printf("DEBUG Hits[0].Title %s\n", searchResults.Hits[0].Fragments["title"])
	log.Printf("DEBUG Hits[0].Fragments[content_description] %s\n", searchResults.Hits[0].Fragments["content_description"])

	//content, _ := json.Marshal(searchResults)
	//log.Printf("DEBUG content: %s\n", content)
	/*
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
	*/
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.New("results-search.html")
	tmpl.ParseFiles(path.Join(templatesDir, "results-search.html"))
	err = tmpl.Execute(w, searchResults)
	if err != nil {
		log.Printf("Can't render %s/%s, %s", templatesDir, "results-search.html", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	// If GET with Query String or POST pass to results handler
	// else display Basic Search Form
	query := r.URL.Query()
	if r.Method == "POST" || len(query) > 0 {
		resultsHandler(w, r)
		return
	}
	// Shared form data fields.
	formData := struct {
		URI      string
		Page     int
		PageSize int
	}{
		URI:      "/",
		Page:     1,
		PageSize: 10,
	}

	// Handle the basic or advanced search form requests.
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/search/advanced/" {
		formData.URI = "/search/advanced/"
		tmpl := template.New("advanced-search.html")
		tmpl, err := tmpl.ParseFiles(path.Join(templatesDir, "advanced-search.html"))
		err = tmpl.Execute(w, formData)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%s", err)))
		}
		return
	}

	formData.URI = "/search/basic/"
	tmpl := template.New("basic-search.html")
	tmpl, err := tmpl.ParseFiles(path.Join(templatesDir, "basic-search.html"))
	err = tmpl.Execute(w, formData)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	}
}

func logRequest(r *http.Request) {
	log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG testing call to logger...")
		logRequest(r)
		next.ServeHTTP(w, r)
	})
}

func init() {
	var err error

	uri := os.Getenv("ASPACE_SEARCH_URL")
	indexName = os.Getenv("ASPACE_BLEVE_INDEX")
	htdocsDir = os.Getenv("ASPACE_HTDOCS")
	templatesDir = os.Getenv("ASPACE_TEMPLATES")
	flag.StringVar(&uri, "search", "http://localhost:8501", "The URL to listen on for search requests")
	flag.StringVar(&indexName, "index", "index.bleve", "specify the Bleve index to use")
	flag.StringVar(&htdocsDir, "htdocs", "htdocs", "specify where to write the HTML files to")
	flag.StringVar(&templatesDir, "templates", "templates/default", "The directory path for templates")
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")

	advancedPage, err = ioutil.ReadFile(path.Join(templatesDir, "advanced-search.html"))
	if err != nil {
		log.Fatalf("Can't read templates/advanced.html, %s", err)
	}
	basicPage, err = ioutil.ReadFile(path.Join(templatesDir, "basic-search.html"))
	if err != nil {
		log.Fatalf("Can't read templates/basic.html, %s", err)
	}

	if uri != "" {
		serviceURL, err = url.Parse(uri)
		if err != nil {
			log.Fatalf("Aspace Search URL not valid, %s, %s", uri, err)
		}
	}
}

func main() {
	var err error
	flag.Parse()
	if help == true {
		usage()
	}

	// Wake up our search engine
	index, err = bleve.Open(indexName)
	if err != nil {
		log.Fatalf("Can't open Bleve index %s, %s", indexName, err)
	}
	defer index.Close()

	// Setup static detail pages
	// Wake up our search webserver
	http.HandleFunc("/search/advanced/", searchHandler)
	http.HandleFunc("/search/basic/", searchHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/search/basic/", http.StatusMovedPermanently)
	})
	// Send static file request to the default handler
	http.Handle("/repositories/", logger(http.FileServer(http.Dir(htdocsDir))))

	log.Printf("Listening on %s\n", serviceURL.String())
	err = http.ListenAndServe(serviceURL.Host, nil)
	if err != nil {
		log.Fatal(err)
	}
}
