package aspace

import (
	"fmt"
)

//
// Public view data structures
//

// PublicAccessionView returns a structure suitable for templating public web content.
type PublicAccessionView struct {
	URI                string   `json:"uri"`
	Title              string   `json:"title"`
	ContentDescription string   `json:"content_description,omitempty"`
	Subjects           []string `json:"subjects,omitempty"`
}

//
// Public view functions
//
func (a *Accession) PublicAccessionView(subjects map[string]*Subject) (*PublicAccessionView, error) {
	return nil, fmt.Errorf("PublicAccessionView() not implemented yet")
}
