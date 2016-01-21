package aspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

//
// Useful view driven data structures
//

// NormalizedAccessionView returns a structure suitable for templating public web content.
type NormalizedAccessionView struct {
	URI                string   `json:"uri"`
	Title              string   `json:"title"`
	ContentDescription string   `json:"content_description,omitempty"`
	Subjects           []string `json:"subjects,omitempty"`
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
