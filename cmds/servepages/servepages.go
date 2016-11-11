//
// This is a static file web server and search service.
// for interacting with ArchivesSpace's REST API
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
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	//"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	// 3rd Party packages
	"github.com/blevesearch/bleve"

	// Caltech Library packages
	"github.com/caltechlibrary/cait"
)

var (
	description = `
 USAGE: %s [OPTIONS]

 OVERVIEW

	%s provides search services defined by CAIT_SITE_URL for the
	website content defined by CAIT_HTDOCS using the index defined
	by CAIT_HTDOCS_INDEX. Additionally a webhook call can be defined
	to trigger an action such as pulling new site content.

 OPTIONS
`
	configuration = `
 CONFIGURATION

 %s can be configured through environment variables. The following
 variables are supported-

   CAIT_SITE_URL

   CAIT_HTDOCS_INDEX

   CAIT_TEMPLATES

   CAIT_WEBHOOK_PATH
   
   CAIT_WEBHOOK_SECRET
   
   CAIT_WEBHOOK_COMMAND

`
	showHelp    bool
	showVersion bool

	indexName      string
	htdocsDir      string
	templatesDir   string
	serviceURL     *url.URL
	webhookPath    string
	webhookSecret  string
	webhookCommand string

	advancedPage []byte
	basicPage    []byte

	index bleve.Index
)

func usage(appName, version string) {
	fmt.Printf(description, appName, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(configuration, appName)
	fmt.Printf("%s %s\n", appName, version)
	os.Exit(0)
}

func mapToSearchQuery(m map[string]interface{}) (*cait.SearchQuery, error) {
	var err error

	// raw is a tempory data structure to sanitize the
	// form request submitted via the query.
	raw := &struct {
		Q         string `json:"q"`
		QExact    string `json:"q_exact"`
		QExcluded string `json:"q_excluded"`
		QRequired string `json:"q_required"`
		Size      int    `json:"size"`
		From      int    `json:"from"`
		AllIDs    bool   `json:"all_ids"`
	}{}

	isQuery := false

	q := new(cait.SearchQuery)
	src, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Can't marshal %+v, %s", m, err)
	}
	err = json.Unmarshal(src, &raw)
	if err != nil {
		return nil, fmt.Errorf("Can't unmarshal %s, %s", src, err)
	}
	if len(raw.Q) > 0 {
		q.Q = raw.Q
		isQuery = true
	}
	if len(raw.QExact) > 0 {
		q.QExact = raw.QExact
		isQuery = true
	}
	if len(raw.QExcluded) > 0 {
		q.QExcluded = q.QExact
	}
	if len(raw.QRequired) > 0 {
		q.QRequired = raw.QRequired
		isQuery = true
	}

	if isQuery == false {
		return nil, fmt.Errorf("Missing query value fields")
	}

	if raw.AllIDs == true {
		q.AllIDs = true
	}

	//Note: if q.Size is not set by the query request pick a nice default value
	if raw.Size <= 1 {
		q.Size = 10
	} else {
		q.Size = raw.Size
	}
	if raw.From < 0 {
		q.From = 0
	} else {
		q.From = raw.From
	}
	return q, nil
}

func urlToRepoAccessionIDs(uri string) (int, int, error) {
	var err error
	repoID := 0
	accessionID := 0

	parts := strings.SplitN(uri, "/", 7)
	if len(parts) > 4 {
		repoID, err = strconv.Atoi(parts[4])
		if err != nil {
			return 0, 0, fmt.Errorf("Cannot parse repository id %s, %s", uri, err)
		}
	}
	if len(parts) >= 6 {
		id := filepath.Base(uri)
		accessionID, err = strconv.Atoi(id)
		if err != nil {
			return repoID, 0, fmt.Errorf("Cannot parse accession id %s, %s", uri, err)
		}
	}
	return repoID, accessionID, nil
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		pageHTML    = "results-search.html"
		pageInclude = "results-search.include"
	)

	urlQuery := r.URL.Query()
	err := r.ParseForm()
	if err != nil {
		responseLogger(r, http.StatusBadRequest, err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error in POST: %s", err)))
		return
	}

	// Collect the submissions fields.
	submission := make(map[string]interface{})
	// Basic Search results
	if r.Method == "GET" {
		for k, v := range urlQuery {
			if k == "all_ids" {
				if b, err := strconv.ParseBool(strings.Join(v, "")); err == nil {
					submission[k] = b
				}
			} else if k == "from" || k == "size" || k == "total" {
				if i, err := strconv.Atoi(strings.Join(v, "")); err == nil {
					submission[k] = i
				}
			} else if k == "q" || k == "q_exact" || k == "q_excluded" || k == "q_required" {
				submission[k] = strings.Join(v, "")
			}
		}
	}

	// Advanced Search results
	if r.Method == "POST" {
		for k, v := range r.Form {
			if k == "all_ids" {
				if b, err := strconv.ParseBool(strings.Join(v, "")); err == nil {
					submission[k] = b
				}
			} else if k == "from" || k == "size" || k == "total" {
				if i, err := strconv.Atoi(strings.Join(v, "")); err == nil {
					submission[k] = i
				}
			} else if k == "q" || k == "q_exact" || k == "q_excluded" || k == "q_required" {
				submission[k] = strings.Join(v, "")
			}
		}
	}

	q, err := mapToSearchQuery(submission)
	if err != nil {
		responseLogger(r, http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	//
	// Note: Add logic to handle basic and advanced search...
	//
	// q           NewQueryStringQuery
	// q_required  NewQueryStringQuery with a + prefix for each strings.Fields(q_required) value
	// q_exact     NewMatchPhraseQuery
	// q_excluded NewQueryStringQuery with a - prefix for each strings.Feilds(q_excluded) value
	//
	var conQry []bleve.Query

	if q.Q != "" {
		conQry = append(conQry, bleve.NewQueryStringQuery(q.Q))
	}
	if q.QExact != "" {
		conQry = append(conQry, bleve.NewMatchPhraseQuery(q.QExact))
	}
	var terms []string
	for _, s := range strings.Fields(q.QRequired) {
		terms = append(terms, fmt.Sprintf("+%s", strings.TrimSpace(s)))
	}
	for _, s := range strings.Fields(q.QExcluded) {
		terms = append(terms, fmt.Sprintf("-%s", strings.TrimSpace(s)))
	}
	if len(terms) > 0 {
		qString := strings.Join(terms, " ")
		conQry = append(conQry, bleve.NewQueryStringQuery(qString))
	}

	qry := bleve.NewConjunctionQuery(conQry)
	if q.Size == 0 {
		q.Size = 10
	}
	searchRequest := bleve.NewSearchRequestOptions(qry, q.Size, q.From, q.Explain)
	if searchRequest == nil {
		responseLogger(r, http.StatusBadRequest, fmt.Errorf("Can't build new search request options %+v, %s", qry, err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("content_description")
	searchRequest.Highlight.AddField("subjects")
	searchRequest.Highlight.AddField("subjects_function")
	searchRequest.Highlight.AddField("subjects_topical")
	searchRequest.Highlight.AddField("extents")

	subjectFacet := bleve.NewFacetRequest("subjects", 3)
	searchRequest.AddFacet("subjects", subjectFacet)

	subjectTopicalFacet := bleve.NewFacetRequest("subjects_topical", 3)
	searchRequest.AddFacet("subjects_topical", subjectTopicalFacet)

	subjectFunctionFacet := bleve.NewFacetRequest("subjects_function", 3)
	searchRequest.AddFacet("subjects_function", subjectFunctionFacet)

	// Return all fields
	searchRequest.Fields = []string{
		"title",
		"identifier",
		"content_description",
		"content_condition",
		"resource_type",
		"access_restrictions",
		"access_restrictions_note",
		"use_restrictins",
		"use_restrictons_note",
		"dates",
		"date_expression",
		"extents",
		"subjects",
		"subjects_function",
		"subjects_topical",
		"linked_agents_creators",
		"linked_agents_subjects",
		"link_agents_sources",
		"digital_objects.title",
		"digital_objects.file_uris",
		"related_resources",
		"deaccessions",
		"accession_date",
		"created",
	}

	searchResults, err := index.Search(searchRequest)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Bleve results error %v, %s", qry, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	// q (ciat.SearchQuery) performs double duty as both the structure for query submission as well
	// as carring the results to support paging and other types of navigation through
	// the query set. Results are a query with the bleve.SearchReults merged
	q.AttachSearchResults(searchResults)
	pageHTML = "results-search.html"
	pageInclude = "results-search.include"

	// Load my templates and setup to execute them
	tmpl, err := cait.AssembleTemplate(path.Join(templatesDir, pageHTML), path.Join(templatesDir, pageInclude))
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Template Errors: %s, %s, %s\n", pageHTML, pageInclude, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Template errors: %s", err)))
		return
	}
	// Render the page
	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, q)
	//err = tmpl.Execute(w, q)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Can't render %s, %s/%s, %s", templatesDir, pageHTML, pageInclude, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Template error"))
		return
	}
	//NOTE: This bit of ugliness is here because I need to allow <mark> elements and ellipis in the results fragments
	w.Write(bytes.Replace(bytes.Replace(bytes.Replace(buf.Bytes(), []byte("&lt;mark&gt;"), []byte("<mark>"), -1), []byte("&lt;/mark&gt;"), []byte("</mark>"), -1), []byte(`…`), []byte(`&hellip;`), -1))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//logRequest(r)
	// If GET with Query String or POST pass to results handler
	// else display Basic Search Form
	query := r.URL.Query()
	if r.Method == "POST" || len(query) > 0 {
		resultsHandler(w, r)
		return
	}

	// Shared form data fields for a New Search.
	formData := struct {
		URI string
	}{
		URI: "/",
	}

	// Handle the basic or advanced search form requests.
	var (
		tmpl *template.Template
		err  error
	)
	w.Header().Set("Content-Type", "text/html")
	if strings.HasPrefix(r.URL.Path, "/search/advanced") == true {
		formData.URI = "/search/advanced/"
		tmpl, err = cait.AssembleTemplate(path.Join(templatesDir, "advanced-search.html"), path.Join(templatesDir, "advanced-search.include"))
		if err != nil {
			fmt.Printf("Can't read advanced-search templates, %s", err)
			return
		}
	} else {
		formData.URI = "/search/basic/"
		tmpl, err = cait.AssembleTemplate(path.Join(templatesDir, "basic-search.html"), path.Join(templatesDir, "basic-search.include"))
		if err != nil {
			log.Printf("Can't read basic-search templates, %s\n", err)
			return
		}
	}

	err = tmpl.Execute(w, formData)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//FIXME: add the response status returned.
		q := r.URL.Query()
		if len(q) > 0 {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q)
		} else {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

func responseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("Response: %s Path: %s RemoteAddr: %s Query: %+v StatusCode: %d Error: %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, err)
	} else {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s StatusCode: %d Error: %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, err)
	}
}

// isMultiViewPath checks to see if the path requested behaves like an Apache MultiView request
func isMultiViewPath(p string) bool {
	// check to see if p plus .html extension exists
	fname := fmt.Sprintf("%s.html", p)
	if _, err := os.Stat(path.Join(htdocsDir, fname)); err == nil {
		return true
	}
	return false
}

func multiViewPath(p string) string {
	return fmt.Sprintf("%s.html", p)
}

func customRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle webhook route
		if webhookPath != "" && strings.HasPrefix(r.URL.Path, webhookPath) == true {
			webhookHandler(w, r)
			return
		}

		// Handler are searches and results
		if strings.HasPrefix(r.URL.Path, "/search/results/") == true {
			resultsHandler(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/search/") == true {
			searchHandler(w, r)
			return
		}

		// If this is a MultiViews style request (i.e. missing .html) then update r.URL.Path
		if isMultiViewPath(r.URL.Path) == true {
			p := multiViewPath(r.URL.Path)
			r.URL.Path = p
		}
		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}

func getenv(envvar, defaultValue string) string {
	tmp := os.Getenv(envvar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func init() {
	var err error

	uri := getenv("CAIT_SITE_URL", "http://localhost:8501")
	htdocsDir = getenv("CAIT_HTDOCS", "htdocs")
	indexName = getenv("CAIT_HTDOCS_INDEX", "htdocs.bleve")
	templatesDir = getenv("CAIT_TEMPLATES", "templates/default")
	webhookPath = getenv("CAIT_WEBHOOK_PATH", "")
	webhookSecret = getenv("CAIT_WEBHOOK_SECRET", "")
	webhookCommand = getenv("CAIT_WEBHOOK_COMMAND", "")

	flag.StringVar(&uri, "search", uri, "The URL to listen on for search requests")
	flag.StringVar(&indexName, "index", indexName, "specify the Bleve index to use")
	flag.StringVar(&htdocsDir, "htdocs", htdocsDir, "specify where to write the HTML files to")
	flag.StringVar(&templatesDir, "templates", templatesDir, "The directory path for templates")
	flag.BoolVar(&showHelp, "h", false, "display this help message")
	flag.BoolVar(&showHelp, "help", false, "display this help message")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&showVersion, "version", false, "display version info")

	flag.StringVar(&webhookPath, "webhook-path", webhookPath, "the webhook path, e.g. /my-webhook/something")
	flag.StringVar(&webhookSecret, "webhook-secret", webhookSecret, "the secret to validate before executing command")
	flag.StringVar(&webhookCommand, "webhook-command", webhookCommand, "the command to execute if webhook validates")

	templateName := path.Join(templatesDir, "advanced-search.html")
	advancedPage, err = ioutil.ReadFile(templateName)
	if err != nil {
		log.Fatalf("Can't read %s, %s", templateName, err)
	}
	templateName = path.Join(templatesDir, "basic-search.html")
	basicPage, err = ioutil.ReadFile(templateName)
	if err != nil {
		log.Fatalf("Can't read %s, %s", templateName, err)
	}

	if uri != "" {
		serviceURL, err = url.Parse(uri)
		if err != nil {
			log.Fatalf("Aspace Search URL not valid, %s, %s", uri, err)
		}
	}
}

// signBody and verifySignature based on Gist https://gist.github.com/rjz/b51dc03061dbcff1c521
func verifySignature(secret []byte, signature string, body []byte) bool {
	signBody := func(secret, body []byte) []byte {
		computed := hmac.New(sha1.New, secret)
		computed.Write(body)
		return []byte(computed.Sum(nil))
	}

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Always reeturn text/plain OK with a 200 to obscure that this actually is.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")

	log.Printf("Webhook Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
	if strings.Compare(r.Method, "POST") == 0 {
		header := r.Header
		contentType := header.Get("Content-Type")
		xGithubSignature := header.Get("X-Hub-Signature")
		if strings.Compare(contentType, "application/json") == 0 && xGithubSignature != "" {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil && verifySignature([]byte(webhookSecret), xGithubSignature, body) == true {
				log.Printf("Webhook validated, running %q", webhookCommand)
				out, err := exec.Command(webhookCommand).Output()
				if err != nil {
					log.Printf("Webhook error: %s", err)
					return
				}
				log.Printf("Webhook out: %s", out)
				return
			}
		}
	}
	log.Printf("Webhook invalid request method.")
	return
}

func main() {
	var err error

	appName := path.Base(os.Args[0])

	flag.Parse()
	if showHelp == true {
		usage(appName, cait.Version)
	}
	if showVersion == true {
		fmt.Printf("%s %s\n", appName, cait.Version)
		os.Exit(0)
	}

	// Wake up our search engine
	index, err = bleve.Open(indexName)
	if err != nil {
		log.Fatalf("Can't open Bleve index %s, %s", indexName, err)
	}
	defer index.Close()

	// Send static file request to the default handler,
	// search routes are handled by middleware customRoutes()
	http.Handle("/", http.FileServer(http.Dir(htdocsDir)))

	log.Printf("%s %s\n", appName, cait.Version)
	log.Printf("Listening on %s\n", serviceURL.String())
	err = http.ListenAndServe(serviceURL.Host, requestLogger(customRoutes(http.DefaultServeMux)))
	if err != nil {
		log.Fatal(err)
	}
}
