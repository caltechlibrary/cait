// views_test.go contains test methods for the views.go and models.go test module.
// the tests assume a successful export of data from ArchivesSpace has been performed.
package cait

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestViews(t *testing.T) {
	datasets := os.Getenv("CAIT_DATASETS")
	apiURL := os.Getenv("CAIT_API_URL")
	username := os.Getenv("CAIT_USERNAME")
	password := os.Getenv("CAIT_PASSWORD")
	if datasets == "" {
		log.Println("$CAIT_DATASETS not set, skipping TestViews()")
		t.SkipNow()
	}
	_, err := os.Stat(datasets)
	if os.IsNotExist(err) == true {
		if apiURL == "" {
			log.Println("$CAIT_API_URL not set, skipping TestViews()")
			t.SkipNow()
		}
		api := New(apiURL, username, password)
		err := api.Login()
		if err != nil {
			t.Errorf("Login %s, error %s", api.URL, err)
			t.FailNow()
		}
		err = api.ExportArchivesSpace()
		if err != nil {
			t.Errorf("Failed to export ArchivesSpace, %s", err)
			t.FailNow()
		}
	}

	subjectDir := path.Join(datasets, "subjects")
	SubjectList, err := MakeSubjectList(subjectDir)
	if err != nil {
		t.Errorf("Failed to make subject list from %s, %s", subjectDir, err)
	}
	SubjectMap, err := MakeSubjectMap(subjectDir)
	if err != nil {
		t.Errorf("Failed to make subject map from %s, %s", subjectDir, err)
	}
	if len(SubjectList) == 0 {
		t.Errorf("SubjectList should not be empty, length %d", len(SubjectList))
	}
	if len(SubjectMap) == 0 {
		t.Errorf("SubjectMap should not be empty, length %d", len(SubjectMap))
	}

	titleIndex, err := MakeAccessionTitleIndex(datasets)
	if err != nil {
		t.Errorf("Should be able to make an accessions title index from %s, %s", datasets, err)
	}
	filepath.Walk(datasets, func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".json") && strings.Contains(p, "/accessions/") {
			log.Printf("Testing %s\n", p)
			//FIXME: test normalized view here...
			src, err := ioutil.ReadFile(p)
			if err != nil {
				t.Errorf("Can't read %s, %s", p, err)
				return err
			}
			accession := new(Accession)
			err = json.Unmarshal(src, &accession)
			if err != nil {
				t.Errorf("Can't unmarshal accession %s, %s", src, err)
				return err
			}
			accessionView, err := accession.NormalizeView(SubjectMap, titleIndex[accession.URI])
			if err != nil {
				t.Errorf("Can't make a normalized view for %s, %s", p, err)
				return err
			}
			nav := accessionView.Nav
			if strings.Compare(nav.NextURI, nav.ThisURI) == 0 || strings.Compare(nav.PrevURI, nav.ThisURI) == 0 {
				t.Errorf("Nav problem for accession %s, %s", p, nav)
				t.FailNow()
			}
			if accessionView.Title != accession.Title {
				t.Errorf("Title does not match %s, [%s] != [%s]", p, accessionView.Title, accession.Title)
			}
		}
		return nil
	})
}