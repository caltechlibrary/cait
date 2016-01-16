/**
 * aspaceindexer.go - A search indexer for [Bleve Search](https://github.com/blevesearch/bleve)
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/blevesearch/bleve"

	"../../../aspace"
)

var (
	help               bool
	aspaceDataSet      string
	aspaceSitePrefix   string
	aspaceBleveIndex   string
	aspaceBleveMapping string
)

const usage = `
 USAGE: aspaceindexer [-h|--help]

 SYNOPSIS

 aspaceindexer is a command line utility to index content fetched from
 an ArchivesSpace instance via the ArchivesSpace REST API (e.g. with
 aspace tool). It indexes content for the Bleve search library.
 Configuration is done through environmental variables.

 OPTIONS
    -h, --help  Display this help page

 CONFIGURATION

 aspaceindexer relies on the following environment variables for
 configuration when overriding the defaults:

    ASPACE_DATASET	(default: data)
                    This should be the path to the directory tree containings
                    the JSON files to be index. E.g. data-export vs. the default
					data.

    ASPACE_BLEVE_INDEX	(default: index.bleve)
                    This is the directory that will contain all the Bleve
                    indexes.

`

func init() {
	flag.BoolVar(&help, "h", false, usage)
	flag.BoolVar(&help, "help", false, usage)
	aspaceDataSet = os.Getenv("ASPACE_DATASET")
	aspaceBleveIndex = os.Getenv("ASPACE_BLEVE_INDEX")

	if aspaceDataSet == "" {
		aspaceDataSet = "./data"
	}
	if aspaceBleveIndex == "" {
		aspaceBleveIndex = "index.bleve"
	}
}

func indexAgents(index bleve.Index, dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fileCount := len(files)

	/*Batch ingest implementation */
	batchSize := 100
	startTime := time.Now()
	batchNo := 0
	batchI := 0
	batch := index.NewBatch()
	for i, fp := range files {
		fname := path.Join(dirname, fp.Name())
		//log.Printf("%d/%d: batching %s\n", i, fileCount, fname)
		src, err := ioutil.ReadFile(fname)
		if err != nil {
			return fmt.Errorf("Can't read %s, %s", fname, err)
		}
		//FIXME: This is the only change in the process, just the structure we're rendering...
		var data *aspace.Agent
		err = json.Unmarshal(src, &data)
		if err != nil {
			return fmt.Errorf("Can't parse %s, %s", fname, err)
		}
		batch.Index(data.URI, data)
		batchI++
		if batchI >= batchSize {
			err = index.Batch(batch)
			if err != nil {
				return fmt.Errorf("Error processing batch %d, %s", batchNo, err)
			}
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(i)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", i, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
			batchNo++
			batchI = 0
			batch.Reset()
		}
	}

	// Run any remaining batch
	if batchI < fileCount {
		err = index.Batch(batch)
		if err != nil {
			return fmt.Errorf("Error processing batch %d, %s", batchNo, err)
		}
		batch.Reset()
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(fileCount)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", fileCount, indexDurationSeconds, timePerDoc/float64(time.Millisecond))

	return nil
}

func indexAccessions(index bleve.Index, dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fileCount := len(files)

	/*Batch ingest implementation */
	batchSize := 100
	startTime := time.Now()
	batchNo := 0
	batchI := 0
	batch := index.NewBatch()
	for i, fp := range files {
		fname := path.Join(dirname, fp.Name())
		//log.Printf("%d/%d: batching %s\n", i, fileCount, fname)
		src, err := ioutil.ReadFile(fname)
		if err != nil {
			return fmt.Errorf("Can't read %s, %s", fname, err)
		}

		//FIXME: This is the only change in the process, just the structure we're rendering...
		var data *aspace.Accession
		err = json.Unmarshal(src, &data)
		if err != nil {
			return fmt.Errorf("Can't parse %s, %s", fname, err)
		}
		batch.Index(data.URI, data)
		batchI++
		if i > 0 && batchI >= batchSize {
			err = index.Batch(batch)
			if err != nil {
				return fmt.Errorf("Error processing batch %d, %s", batchNo, err)
			}
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(i)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", i, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
			batchNo++
			batchI = 0
			batch.Reset()
		}
	}

	// Run any remaining batch
	if batchI < fileCount {
		err = index.Batch(batch)
		if err != nil {
			return fmt.Errorf("Error processing batch %d, %s", batchNo, err)
		}
		batch.Reset()
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(fileCount)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", fileCount, indexDurationSeconds, timePerDoc/float64(time.Millisecond))

	return nil
}

func main() {
	var (
		index bleve.Index
		err   error
	)

	flag.Parse()
	if help == true {
		fmt.Println(usage)
		os.Exit(0)
	}

	if _, err = os.Stat(aspaceBleveIndex); os.IsNotExist(err) {
		log.Printf("Creating Bleve index at %s\n", aspaceBleveIndex)
		mapping := bleve.NewIndexMapping()
		mapping.DefaultAnalyzer = "en"
		index, err = bleve.New(aspaceBleveIndex, mapping)
		if err != nil {
			log.Printf("Can't create new bleve index %s, %s", aspaceBleveIndex, err)
		}
	} else {
		log.Printf("Opening Bleve index at %s\n", aspaceBleveIndex)
		index, err = bleve.Open(aspaceBleveIndex)
		if err != nil {
			log.Printf("Can't open bleve index %s, %s", aspaceBleveIndex, err)
		}
	}
	defer index.Close()

	// Walk our data import tree and index things
	log.Printf("Start indexing of %s in %s\n", aspaceDataSet, aspaceBleveIndex)
	wholeProcessStartTime := time.Now()
	dirCount := 0

	dataSet := path.Join(aspaceDataSet, "repositories", "2", "accessions")
	log.Printf("Indexing %s\n", dataSet)
	err = indexAccessions(index, dataSet)
	if err != nil {
		log.Printf("Can't properly index %s, %s\n", dataSet, err)
	}
	dirCount++

	for _, folder := range []string{"corporate_entities", "people", "families", "software"} {
		dataSet := path.Join(aspaceDataSet, "agents", folder)
		log.Printf("Indexing %s\n", dataSet)
		err := indexAgents(index, dataSet)
		if err != nil {
			log.Printf("Can't properly index %s, %s\n", dataSet, err)
		}
		dirCount++
	}

	indexDuration := time.Since(wholeProcessStartTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	log.Printf("Finished, Indexed %d directories, in %.2fs for %s in %s\n",
		dirCount, indexDurationSeconds, aspaceDataSet, aspaceBleveIndex)
}
