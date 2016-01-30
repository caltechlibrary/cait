//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, California Institute of Technology
// All rights reserved.
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
	"path"
)

//
// Useful view driven data structures and functions
//

// NormalizedAccessionView returns a structure suitable for templating public web content.
type NormalizedAccessionView struct {
	URI                string   `json:"uri"`
	Title              string   `json:"title"`
	ContentDescription string   `json:"content_description,omitempty"`
	ConditionDescription string `json:"content_description,omitempty"`
	Subjects           []string `json:"subjects,omitempty"`
	Extents []string `json:"extents,omitempty"`
	RelatedResource []string `json:"related_resources,omitempty"`
	Instances []string `json:"instances,omitempty"`
	LinkedAgents []string `json:"linked_agents,omitempty"`
	DigitalObjects []map[string]interface{} `json:"digital_objects,omitempty"`
}

// MakeSubjectList given a base data directory read in the subject JSON blobs and builds
// a slice or subject data.
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
// a slice or subject data.
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
func (a *Accession) NormalizeView(subjects map[string]*Subject) (*NormalizedAccessionView, error) {
	var subjectLabels []string
	v := new(NormalizedAccessionView)
	v.Title = a.Title
	v.URI = a.URI
	v.ContentDescription = a.ContentDescription
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
	return v, nil
}
