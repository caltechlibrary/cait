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
	"log"
	"sort"
	"strings"
	"time"
)

//
// Useful view driven data structures and functions
//

// PageView is a simple container for rendering pages
type PageView struct {
	Nav     NavView
	Content []interface{}
}

// NavElementView defined previous, next links used in paging results or browsable record lists
type NavElementView struct {
	ThisLabel string `json:"this_label"`
	ThisURI   string `json:"this_uri"`
	PrevURI   string `json:"prev_uri"`
	PrevLabel string `json:"prev_label"`
	NextURI   string `json:"next_uri"`
	NextLabel string `json:"next_label"`
	Weight    int    `json:"weight"`
}

// NavView is an array of NavelementViews
type NavView []*NavElementView

// NormalizedDigitalObjectView returns a structure suitable for templating public web content.
type NormalizedDigitalObjectView struct {
	URI               string   `json:"uri"`
	Title             string   `json:"title"`
	DigitalObjectType string   `json:"digital_object_type"`
	Publish           bool     `json:"publish"`
	FileURIs          []string `json:"file_uris"`
}

// NormalizedAccessionView returns a structure suitable for templating public web content.
type NormalizedAccessionView struct {
	ID                     string                         `json:"id"`
	URI                    string                         `json:"uri"`
	Title                  string                         `json:"title"`
	Identifier             string                         `json:"identifier"`
	ResourceType           string                         `json:"resource_type"`
	ContentDescription     string                         `json:"content_description"`
	ConditionDescription   string                         `json:"condition_description,omitempty"`
	AccessRestrictions     bool                           `json:"access_restrictions"`
	AccessRestrictionsNote string                         `json:"access_restrictions_notes"`
	UseRestrictions        bool                           `json:"use_restrictions"`
	UseRestrictionsNote    string                         `json:"use_restrictions_notes"`
	Dates                  []*Date                        `json:"dates"`
	DateExpression         string                         `json:"date_expression"`
	Subjects               []string                       `json:"subjects,omitempty"`
	SubjectsFunction       []string                       `json:"subjects_function,omitempty"`
	SubjectsTopical        []string                       `json:"subjects_topical,omitempty"`
	Extents                []string                       `json:"extents"`
	RelatedResources       []string                       `json:"related_resources,omitempty"`
	RelatedAccessions      []string                       `json:"related_accessions,omitempty"`
	DigitalObjects         []*NormalizedDigitalObjectView `json:"digital_objects,omitempty"`
	Deaccessions           string                         `json:"deaccessions,omitempty"`
	LinkedAgentsCreators   []string                       `json:"linked_agents_creators"`
	LinkedAgentsSubjects   []string                       `json:"linked_agents_subjects"`
	LinkedAgentsSources    []string                       `json:"linked_agents_sources,omitempty"`
	AccessionDate          string                         `json:"accession_date"`
	CreatedBy              string                         `json:"created_by"`
	Created                string                         `json:"created"`
	LastModifiedBy         string                         `json:"last_modified_by"`
	LastModified           string                         `json:"last_modified"`
}

// FlattenDates takes an array of Date types, flatten it into a human readable string.
func FlattenDates(dates []*Date) string {
	var out []string
	for _, dt := range dates {
		switch dt.DateType {
		case "single":
			d, _ := time.Parse("2006-01-02", dt.Expression)
			out = append(out, fmt.Sprintf("%s", d.Format("Jan. 2, 2006")))
		case "inclusive":
			start, _ := time.Parse("2006-01-02", dt.Begin)
			end, _ := time.Parse("2006-01-02", dt.End)
			out = append(out, fmt.Sprintf("%s - %s", start.Format("Jan. 2, 2006"), end.Format("Jan. 2, 2006")))
		}
	}
	return strings.Join(out, "; ")
}

// NormalizeView returns a normalized view from an Accession structure and
// an array of subject structures.
func (a *Accession) NormalizeView(agents []*Agent, subjects map[string]*Subject, digitalObjects map[string]*DigitalObject) (*NormalizedAccessionView, error) {
	agentMap := make(map[string]string)
	for _, agent := range agents {
		title := agent.Title
		uri := agent.URI
		agentMap[uri] = title
	}
	v := new(NormalizedAccessionView)
	v.ID = fmt.Sprintf("%d", a.ID)
	v.Title = a.Title
	v.Identifier = strings.Trim(strings.Join([]string{a.ID0, a.ID1, a.ID2, a.ID3}, "-"), "-")
	v.ResourceType = a.ResourceType
	v.URI = a.URI
	v.ContentDescription = a.ContentDescription
	v.ConditionDescription = a.ConditionDescription
	v.AccessRestrictions = a.AccessRestrictions
	v.AccessRestrictionsNote = a.AccessRestrictionsNote
	v.UseRestrictions = a.UseRestrictions
	v.UseRestrictionsNote = a.UseRestrictionsNote
	v.Dates = a.Dates
	v.DateExpression = FlattenDates(a.Dates)
	v.AccessionDate = a.AccessionDate
	v.CreatedBy = a.CreatedBy
	v.Created = a.CreateTime
	v.LastModifiedBy = a.LastModifiedBy
	v.LastModified = a.UserMTime
	for _, extent := range a.Extents {
		v.Extents = append(v.Extents, extent.PhysicalDetails)
	}
	for _, item := range a.Instances {
		//FIXME: assign the URL link to v.Instance
		if m, ok := item["digital_object"]; ok == true {
			kv := map[string]string{}
			src, _ := json.Marshal(m)
			json.Unmarshal(src, &kv)
			if ref, ok := kv["ref"]; ok == true {
				if obj, ok := digitalObjects[ref]; ok == true {
					v.DigitalObjects = append(v.DigitalObjects, obj.NormalizeView())
				}
			}
		}
	}
	for _, item := range a.Subjects {
		if ref, ok := item["ref"]; ok == true {
			rec := subjects[fmt.Sprintf("%s", ref)]
			if rec != nil {
				v.Subjects = append(v.Subjects, rec.Title)
				if len(rec.Terms) > 0 {
					if termType, ok := rec.Terms[0]["term_type"]; ok == true && termType == "function" {
						v.SubjectsFunction = append(v.SubjectsFunction, rec.Title)
					} else {
						v.SubjectsTopical = append(v.SubjectsTopical, rec.Title)
					}
				} else {
					v.SubjectsTopical = append(v.SubjectsTopical, rec.Title)
				}
			}
		}
	}
	//NOTE: Normalized view adds Linked Agents by type creator, subject, sources ...
	for _, item := range a.LinkedAgents {
		if ref, ok := item["ref"].(string); ok == true {
			if title, found := agentMap[ref]; found == true {
				role, _ := item["role"]
				switch role {
				case "creator":
					v.LinkedAgentsCreators = append(v.LinkedAgentsCreators, title)
				case "subject":
					v.LinkedAgentsSubjects = append(v.LinkedAgentsSubjects, title)
				case "source":
					v.LinkedAgentsSources = append(v.LinkedAgentsSources, title)
				}
			}
		}
	}
	return v, nil
}

//FIXME: NormalizeView takes an Agent/People object and returns a normalized view

// NormalizeView takes a digital object and returns a normalized view
func (o *DigitalObject) NormalizeView() *NormalizedDigitalObjectView {
	result := new(NormalizedDigitalObjectView)
	result.URI = o.URI
	result.Title = o.Title
	result.Publish = o.Publish
	result.DigitalObjectType = o.DigitalObjectType
	for _, fv := range o.FileVersions {
		if fv.FileURI != "" {
			result.FileURIs = append(result.FileURIs, fv.FileURI)
		}
	}
	return result
}

// MakeAgentList given a base data directory read in the subject JSON blobs and builds
// a slice or subject data. Takes the path to the subjects directory as a parameter.
func MakeAgentList(dname string) ([]*Agent, error) {
	var agents []*Agent

	keys, err := GetKeys(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read Agents keys from %s, %s", dname, err)
	}
	for _, key := range keys {
		src, err := ReadJSON(dname, key)
		if err != nil {
			return nil, fmt.Errorf("Can't read Agent %s, %s", key, err)
		}
		agent := new(Agent)
		err = json.Unmarshal(src, &agent)
		if err != nil {
			return nil, fmt.Errorf("Can't parse Agent %s, %s", key, err)
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

type subjectList []string

func (s subjectList) Len() int {
	return len(s)
}

func (s subjectList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s subjectList) Less(i, j int) bool {
	return (strings.Compare(s[i], s[j]) == -1)
}

func (s subjectList) HasSubject(a string) bool {
	for _, b := range s {
		if strings.Compare(a, b) == 0 {
			return true
		}
	}
	return false
}

// MakeSubjectList given a base data directory read in the subject JSON blobs and builds
// a slice or subject data. Takes the path to the subjects directory as a parameter.
func MakeSubjectList(dname string) ([]string, error) {
	var subjects subjectList

	keys, err := GetKeys(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read Subjects from %s, %s", dname, err)
	}
	for _, key := range keys {
		src, err := ReadJSON(dname, key)
		if err != nil {
			return nil, fmt.Errorf("Can't read Subjects %s, %s", key, err)
		}
		subject := new(Subject)
		err = json.Unmarshal(src, &subject)
		if err != nil {
			return nil, fmt.Errorf("Can't parse Subject %s, %s", key, err)
		}
		if subject.Publish == true {
			for _, term := range subject.Terms {
				if val, ok := term["term"]; ok == true {
					sval := fmt.Sprintf("%s", val)
					if subjects.HasSubject(sval) == false {
						subjects = append(subjects, sval)
					}
				}
			}
		}
	}

	sort.Sort(subjects)
	return subjects, nil
}

// MakeSubjectMap given a base data directory read in the subject JSON blobs and builds
// a map or subject data. Takes the path to the subjects directory as a parameter.
func MakeSubjectMap(dname string) (map[string]*Subject, error) {
	subjects := make(map[string]*Subject)

	keys, err := GetKeys(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read subject keys from %s, %s", dname, err)
	}
	for _, key := range keys {
		src, err := ReadJSON(dname, key)
		if err != nil {
			return nil, fmt.Errorf("Can't read %s, %s", key, err)
		}
		subject := new(Subject)
		err = json.Unmarshal(src, &subject)
		if err != nil {
			return nil, fmt.Errorf("Can't parse subject %s, %s", key, err)
		}
		subjects[subject.URI] = subject
	}
	return subjects, nil
}

// MakeDigitalObjectMap given a base data directory read in the Digital Object JSON blobs
// and build a map of object data. Takes the path to the subjects directory as a parameter.
func MakeDigitalObjectMap(dname string) (map[string]*DigitalObject, error) {
	digitalObjects := make(map[string]*DigitalObject)

	keys, err := GetKeys(dname)
	if err != nil {
		return nil, fmt.Errorf("Can't read Digital Objects from %s, %s", dname, err)
	}
	for _, key := range keys {
		src, err := ReadJSON(dname, key)
		if err != nil {
			return nil, fmt.Errorf("Can't read Digital Object %s, %s", key, err)
		}
		digitalObject := new(DigitalObject)
		err = json.Unmarshal(src, &digitalObject)
		if err != nil {
			return nil, fmt.Errorf("Can't parse subject %s, %s", key, err)
		}
		digitalObjects[digitalObject.URI] = digitalObject
	}
	return digitalObjects, nil

}

//
// Browsing data
//

// MakeAccessionTitleIndex crawls the path for accession records and generates
// a map of navigation links that can be used in search results or browsing views.
// The parameter dname usually is set to the value of $CAIT_DATASET
// Output is a map of URI pointing at NavElementView for that URI.
func MakeAccessionTitleIndex(dname string) (map[string]*NavElementView, error) {
	// Title index keyed by URI
	titleIndex := make(map[string]*NavElementView)
	titlesWithURI := []string{}
	log.Printf("Making Accession Title Index")
	keys, err := GetKeys(dname)
	if err != nil {
		return nil, fmt.Errorf("MakeAccessionTitleIndex(%q) -> GetKeys(%q), %s", dname, dname, err)
	}
	for _, key := range keys {
		src, err := ReadJSON(dname, key)
		if err != nil {
			log.Printf("Can't read Accession %s, %s", key, err)
		} else {
			accession := new(struct {
				Title     string `json:"title,omitempty"`
				URI       string `json:"uri"`
				JSONModel string `json:"jsonmodel_type"`
			})
			err = json.Unmarshal(src, &accession)
			if err != nil {
				log.Printf("Can't unpack accession info %s, %s", key, err)
			}
			if accession.JSONModel == "accession" {
				//FIXME: Store the info.
				nav := new(NavElementView)
				nav.ThisLabel = accession.Title
				nav.ThisURI = accession.URI
				titleIndex[accession.URI] = nav
				titlesWithURI = append(titlesWithURI, fmt.Sprintf("%s|%s", accession.Title, accession.URI))
			}
			log.Printf("Recorded %s", key)
		}
	}

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
func (nav *NavElementView) String() string {
	var (
		prev string
		this string
		next string
	)
	if nav.PrevURI != "" {
		prev = fmt.Sprintf(`<a class="prev-item" href="%s" title="%s">prev</a>`, nav.PrevURI, nav.PrevLabel)
	}
	if nav.ThisURI != "" {
		this = fmt.Sprintf(`<span class="this-item" data-uri="%s" data-title="%s">%s</span>`, nav.ThisURI, nav.ThisLabel, nav.ThisLabel)
	}
	if nav.NextURI != "" {
		next = fmt.Sprintf(`<a class="next-item" href="%s" title="%s">next</a>`, nav.NextURI, nav.NextLabel)
	}
	return strings.Trim(fmt.Sprintf("%s %s %s", prev, this, next), " ")
}
