//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
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
package cait

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

//
// Useful view driven data structures and functions
//

// NavRecord defined previous, next links used in paging results or browsable record lists
type NavRecord struct {
	ThisLabel string `json:"this_label,omitempty"`
	ThisURI   string `json:"this_uri, omitempty"`
	PrevURI   string `json:"prev_uri,omitempty"`
	PrevLabel string `json:"prev_label,omitempty"`
	NextURI   string `json:"next_uri,omitempty"`
	NextLabel string `json:"next_label,omitempty"`
}

// NormalizedAccessionView returns a structure suitable for templating public web content.
type NormalizedAccessionView struct {
	URI                  string     `json:"uri"`
	Title                string     `json:"title"`
	ContentDescription   string     `json:"content_description,omitempty"`
	ConditionDescription string     `json:"content_description,omitempty"`
	Subjects             []string   `json:"subjects,omitempty"`
	Extents              []*Extent  `json:"extents,omitempty"`
	RelatedResources     []string   `json:"related_resources,omitempty"`
	RelatedAccessions    []string   `json:"related_accessions,omitempty"`
	Instances            []string   `json:"instances,omitempty"`
	LinkedAgents         []string   `json:"linked_agents,omitempty"`
	CreatedBy            string     `json:"created_by"`
	Created              string     `json:"created"`
	LastModifiedBy       string     `json:"last_modified_by,omitempty"`
	LastModified         string     `json:"last_modified"`
	Nav                  *NavRecord `json:"nav,omitempty"`
}

// MakeSubjectList given a base data directory read in the subject JSON blobs and builds
// a slice or subject data. Takes the path to the subjects directory as a parameter.
func MakeSubjectList(dname string) ([]*Subject, error) {
	var subjects []*Subject

	dir, err := ioutil.ReadDir(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read subjects from %s, %s", dname, err)
	}
	for _, finfo := range dir {
		fname := path.Join(dname, finfo.Name())
		src, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, fmt.Errorf("Can't read %s, %s", fname, err)
		}
		subject := new(Subject)
		err = json.Unmarshal(src, &subject)
		if err != nil {
			return nil, fmt.Errorf("Can't parse subject %s, %s", fname, err)
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

// MakeSubjectMap given a base data directory read in the subject JSON blobs and builds
// a slice or subject data. Takes the path to the subjects directory as a parameter.
func MakeSubjectMap(dname string) (map[string]*Subject, error) {
	subjects := make(map[string]*Subject)

	dir, err := ioutil.ReadDir(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read subjects from %s, %s", dname, err)
	}
	for _, finfo := range dir {
		fname := path.Join(dname, finfo.Name())
		src, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, fmt.Errorf("Can't read %s, %s", fname, err)
		}
		subject := new(Subject)
		err = json.Unmarshal(src, &subject)
		if err != nil {
			return nil, fmt.Errorf("Can't parse subject %s, %s", fname, err)
		}
		subjects[subject.URI] = subject
	}
	return subjects, nil
}

// NormalizeView returns a normalized view from an Accession structure and
// an array of subject structures.
func (a *Accession) NormalizeView(subjects map[string]*Subject, nav *NavRecord) (*NormalizedAccessionView, error) {
	var (
		subjectLabels []string
		instanceLinks []string
	)

	v := new(NormalizedAccessionView)
	v.Title = a.Title
	v.URI = a.URI
	v.ContentDescription = a.ContentDescription
	v.ConditionDescription = a.ConditionDescription
	v.CreatedBy = a.CreatedBy
	v.Created = a.CreateTime
	v.LastModifiedBy = a.LastModifiedBy
	v.LastModified = a.UserMTime
	v.Extents = a.Extents
	//FIXME: need to figure out how to handle these...
	//v.RelatedResources = a.RelatedResources
	//v.RelatedAccessions = a.RelatedAccessions
	//v.LinkedAgents = a.LinkedAgents
	for _, item := range a.Instances {
		fmt.Printf("DEBUG instance item: %+v\n", item)
		//FIXME: assign the URL link to v.Instance
		instanceLinks = append(instanceLinks, fmt.Sprintf("%+v", item))
	}
	for _, item := range a.Subjects {
		ref, ok := item["ref"]
		if ok == true {
			rec := subjects[fmt.Sprintf("%s", ref)]
			if rec != nil {
				subjectLabels = append(subjectLabels, rec.Title)
			}
		}
	}
	v.Subjects = subjectLabels
	if nav != nil {
		v.Nav = nav
	}
	return v, nil
}

//
// Browsing data
//

// MakeAccessionTitleIndex crawls the path for accession records and generates
// a map of navigation links that can be used in search results or browsing views.
// The parameter dname usually is set to the value of $CAIT_DATASETS
// Output is a map of URI pointing at NavRecord for that URI.
func MakeAccessionTitleIndex(dname string) (map[string]*NavRecord, error) {
	// Title index keyed by URI
	titleIndex := make(map[string]*NavRecord)
	titlesWithURI := []string{}
	log.Printf("Making Accession Title Index")
	filepath.Walk(dname, func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".json") {
			src, err := ioutil.ReadFile(p)
			if err != nil {
				log.Printf("Can't read %s, %s", p, err)
				return nil
			}
			accession := new(struct {
				Title     string `json:"title,omitempty"`
				URI       string `json:"uri"`
				JSONModel string `json:"jsonmodel_type"`
			})
			err = json.Unmarshal(src, &accession)
			if err != nil {
				log.Printf("Can't unpack accession info %s, %s", p, err)
			}
			if accession.JSONModel == "accession" {
				//FIXME: Store the info.
				nav := new(NavRecord)
				nav.ThisLabel = accession.Title
				nav.ThisURI = accession.URI
				titleIndex[accession.URI] = nav
				titlesWithURI = append(titlesWithURI, fmt.Sprintf("%s|%s", accession.Title, accession.URI))
			}
			log.Printf("Recorded %s", p)
		}
		return nil
	})

	if len(titlesWithURI) == 0 {
		return nil, fmt.Errorf("No titles found")
	}
	if len(titleIndex) == 0 {
		return nil, fmt.Errorf("title index empty")
	}

	// make a uri extraction func
	extractURI := func(s string) string {
		pos := strings.LastIndex(s, "|")
		pos++
		return s[pos:]
	}

	// Sort the titles
	log.Printf("Sorting %d titles", len(titlesWithURI))
	sort.Strings(titlesWithURI)
	// go through sorted titles and populate Next and Prev appropriately
	log.Printf("Linked %d titles", len(titleIndex))
	lastI := len(titlesWithURI) - 1
	for i, val := range titlesWithURI {
		uri := extractURI(val)
		log.Printf("DEBUG val: %s -> uri: %s", val, uri)
		_, thisOk := titleIndex[uri]
		if thisOk == true {
			if i > 0 {
				prevURI := extractURI(titlesWithURI[i-1])
				prev, prevOK := titleIndex[prevURI]
				if prevOK == true {
					titleIndex[uri].PrevLabel = prev.ThisLabel
					titleIndex[uri].PrevURI = prev.ThisURI
				}
			}

			if i < lastI {
				nextURI := extractURI(titlesWithURI[i+1])
				next, nextOK := titleIndex[nextURI]
				if nextOK == true {
					titleIndex[uri].NextLabel = next.ThisLabel
					titleIndex[uri].NextURI = next.ThisURI
				}
			}
		}
		log.Printf("%s, nav: %s\n", uri, titleIndex[uri])
	}
	return titleIndex, nil
}

//
// String() implementations
//
func (nav *NavRecord) String() string {
	var (
		prev string
		this string
		next string
	)
	if nav.PrevURI != "" {
		prev = fmt.Sprintf(`<a class="prev-item" href="%s" title="%s">prev</a>`, nav.PrevURI, nav.PrevLabel)
	}
	if nav.ThisURI != "" {
		this = fmt.Sprintf(`<span class="this-item" data-uri="%s" data-title="%s"></span>`, nav.ThisURI, nav.ThisLabel)
	}
	if nav.NextURI != "" {
		next = fmt.Sprintf(`<a class="next-item" href="%s" title="%s">next</a>`, nav.NextURI, nav.NextLabel)
	}
	return strings.Trim(fmt.Sprintf("%s %s %s", prev, this, next), " ")
}
