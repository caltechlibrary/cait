package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"../../../cait"
)

var (
	help bool
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.BoolVar(&help, "h", false, "display this message")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if help == true || len(args) == 0 {
		fmt.Println("USAGE: normalizer PATH_TO_ACCESSION_JSON")
		flag.PrintDefaults()
		os.Exit(0)
	}
	dataset := os.Getenv("CAIT_DATASET")
	subjectMap, err := cait.MakeSubjectMap(path.Join(dataset, "subjects"))
	check(err)
	digitalObjectMap, err := cait.MakeDigitalObjectMap(path.Join(dataset, "repositories/2/digital_objects"))
	check(err)
	for _, arg := range args {
		src, err := ioutil.ReadFile(arg)
		check(err)
		accession := new(cait.Accession)
		err = json.Unmarshal(src, &accession)
		check(err)
		view, err := accession.NormalizeView(subjectMap, digitalObjectMap)
		check(err)
		src, err = json.Marshal(view)
		check(err)
		fmt.Printf("%s", src)
	}
}
