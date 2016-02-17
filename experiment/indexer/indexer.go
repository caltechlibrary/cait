package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"../../../cait"
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
			log.Printf("Index %d accessions", i)
			batch = index.NewBatch()
		}
	}
	if batch.Size() > 0 {
		err := index.Batch(batch)
		i += batch.Size()
		log.Printf("Index %d accessions in %s", i, time.Now().Sub(startT))
		return err
	}
	return nil
}

func main() {
	datasets := os.Getenv("CAIT_DATASETS")
	indexName := os.Getenv("CAIT_BLEVE_INDEX")

	log.Println("Building subject map...")
	subjectMap, _ := cait.MakeSubjectMap(path.Join(datasets, "repositories/2/subjects"))
	log.Println("Building digital object map...")
	digitalObjectMap, _ := cait.MakeDigitalObjectMap(path.Join(datasets, "repositories/2/digital_objects"))

	log.Println("Setting up index...")
	indexMapping := bleve.NewIndexMapping()
	// Add Accession as a specific document map
	accessionMapping := bleve.NewDocumentMapping()

	// Now add specific accession fields
	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = "en"
	titleMapping.Store = true
	titleMapping.Index = true
	accessionMapping.AddFieldMappingsAt("title", titleMapping)

	descriptionMapping := bleve.NewTextFieldMapping()
	descriptionMapping.Analyzer = "en"
	descriptionMapping.Store = true
	descriptionMapping.Index = true
	accessionMapping.AddFieldMappingsAt("content_description", descriptionMapping)

	objectsMapping := bleve.NewTextFieldMapping()
	objectsMapping.Analyzer = "en"
	objectsMapping.Store = true
	objectsMapping.Index = true
	accessionMapping.AddFieldMappingsAt("digital_objects.title", objectsMapping)

	extentsMapping := bleve.NewTextFieldMapping()
	extentsMapping.Analyzer = "en"
	extentsMapping.Store = true
	extentsMapping.Index = true
	accessionMapping.AddFieldMappingsAt("extents", extentsMapping)

	//FIXME: seems like this could benefit from a better analyzer, e.g. pulling out the terms
	// Do I need to iterate over each item in the subject? Do I need a custom analyzer?
	subjectsMapping := bleve.NewTextFieldMapping()
	subjectsMapping.Analyzer = "en"
	subjectsMapping.Store = true
	subjectsMapping.Index = true
	subjectsMapping.IncludeTermVectors = true
	accessionMapping.AddFieldMappingsAt("subjects", subjectsMapping)

	createdMapping := bleve.NewDateTimeFieldMapping()
	accessionMapping.AddFieldMappingsAt("created", createdMapping)

	// Finally add this mapping to the main index mapping
	indexMapping.AddDocumentMapping("accession", accessionMapping)

	index, _ := openIndex(indexName, indexMapping)
	log.Println("Start indexing...")
	startT := time.Now()
	indexSite(index, 1000, (func() map[string]interface{} {
		i := 0
		m := make(map[string]interface{})
		filepath.Walk(path.Join(datasets, "repositories/2/accessions"), func(p string, _ os.FileInfo, _ error) error {
			if strings.HasSuffix(p, ".json") {
				if (i % 100) == 0 {
					log.Printf("Read %d accessions", i)
				}
				src, _ := ioutil.ReadFile(p)
				data := new(cait.Accession)
				err := json.Unmarshal(src, &data)
				if err == nil {
					i++
					m[data.URI], _ = data.NormalizeView(subjectMap, digitalObjectMap)
				}
				return err
			}
			return nil
		})
		log.Printf("Read %d accessions in %s", i, time.Now().Sub(startT))
		return m
	})())
	log.Printf("Done! %s", time.Now().Sub(startT))
}
