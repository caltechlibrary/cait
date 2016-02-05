package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../../cait"
	"github.com/blevesearch/bleve"
)

func openIndex(indexName string, indexMapping *bleve.IndexMapping) (bleve.Index, error) {
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		return bleve.New(indexName, indexMapping)
	}
	return bleve.Open(indexName)
}

func indexSite(index bleve.Index, batchSize int, dataSet map[string]interface{}) error {
	batch := index.NewBatch()
	startT := time.Now()
	i := 0
	for id, data := range dataSet {
		batch.Index(id, data)
		if batch.Size() >= batchSize {
			err := index.Batch(batch)
			if err != nil {
				return err
			}
			i += batch.Size()
			log.Printf("Index %d items in %s", i, time.Now().Sub(startT))
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		err := index.Batch(batch)
		i += batch.Size()
		log.Printf("Index %d items in %s", i, time.Now().Sub(startT))
		return err
	}
	return nil
}

func main() {
	log.Println("Building a subject map...")
	subjectMap, _ := cait.MakeSubjectMap("../data/repositories/2/subjects/")
	log.Println("Setting up index")
	indexMapping := bleve.NewIndexMapping()
	// Add Accession as a specific document map
	accessionMapping := bleve.NewDocumentMapping()

	// Now add specific accession fields
	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = "en"
	titleMapping.Store = true
	accessionMapping.AddFieldMappingsAt("title", titleMapping)

	descriptionMapping := bleve.NewTextFieldMapping()
	descriptionMapping.Analyzer = "en"
	descriptionMapping.Store = true
	accessionMapping.AddFieldMappingsAt("content_description", descriptionMapping)

	conditionMapping := bleve.NewTextFieldMapping()
	conditionMapping.Analyzer = "en"
	conditionMapping.Store = true
	accessionMapping.AddFieldMappingsAt("condition_description", conditionMapping)

	extentsMapping := bleve.NewTextFieldMapping()
	extentsMapping.Analyzer = "en"
	extentsMapping.Store = true
	accessionMapping.AddFieldMappingsAt("extents", extentsMapping)

	createdMapping := bleve.NewDateTimeFieldMapping()
	createdMapping.Store = true
	accessionMapping.AddFieldMappingsAt("created", createdMapping)

	// Add Subjects as a facet

	// Finally add this mapping to the main index mapping
	indexMapping.AddDocumentMapping("accession", accessionMapping)

	index, _ := openIndex("test.bleve", indexMapping)
	log.Println("Start indexing...")
	startT := time.Now()
	indexSite(index, 100, (func() map[string]interface{} {
		i := 0
		m := make(map[string]interface{})
		filepath.Walk("../data/repositories/2/accessions/", func(p string, _ os.FileInfo, _ error) error {
			if strings.HasSuffix(p, ".json") {
				log.Printf("Reading %s\n", p)
				src, _ := ioutil.ReadFile(p)
				data := new(cait.Accession)
				err := json.Unmarshal(src, &data)
				if err == nil {
					i++
					m[data.URI], _ = data.NormalizeView(subjectMap)
				}
				return err
			}
			return nil
		})
		log.Printf("Read %d items in %s", i, time.Now().Sub(startT))
		return m
	})())
	log.Println("Done! %s", time.Now().Sub(startT))
}
