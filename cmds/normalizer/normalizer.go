//
// cmds/normalizer/normalizer.go - utility to display a normalized view of the accessions in the dataset directory.
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
