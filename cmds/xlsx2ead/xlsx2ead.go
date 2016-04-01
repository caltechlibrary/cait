//
// xlsx2ead.go is a command line utility design to finding aids described in an Excel workbook
// and turn them into suitable EADs for importation into ArchivesSpace via the Create -> Background Jobs -> Import Data -> EAD.
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
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"

	// 3rd Party packages
	"github.com/robertkrimen/otto"
	"github.com/tealeg/xlsx"

	// Caltech Library maintained packages
	"github.com/caltechlibrary/cait"
	"github.com/caltechlibrary/ostdlib"
)

var (
	showHelp      bool
	showVersion   bool
	jsInteractive bool
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "display this help information")
	flag.BoolVar(&showVersion, "v", false, "display version information")
	flag.BoolVar(&jsInteractive, "i", false, "run a JavaScript Repl after loading spreadsheet")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if showHelp == true {
		fmt.Printf(`
 USAGE: xlsx2ead [OPTIONS] EXCEL_FILENAME

 OVERVIEW

 Read an .xlsx file of finding aid information and accession URI and generate
 an EAD suitable for import into ArchivesSpace via the "Background Jobs".

 OPTIONS

  -h          display help information
  -v          display version information
  -i          run an interactive JavaScript shell

 EXAMPLES

    xlsx2ead myFindingAid.xlsx
    xlsx2ead -i myFindingAid.xlsx
    xlsx2ead -h
    xlsx2ead -v

 Version %s
`, cait.Version)
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Printf("Version %s\n", cait.Version)
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Fatalf("Missing excel filename")
	}

	api := cait.New(os.Getenv("CAIT_URL"), os.Getenv("CAIT_USERNAME"), os.Getenv("CAIT_PASSWORD"))

	vm := otto.New()
	js := ostdlib.New(vm)
	js.AddExtensions()
	cait.AddExtensions(api, js)

	var resources []*cait.Resource
	if jsInteractive == true {
		js.AddHelp()
		cait.AddHelp(api, js)
		js.PrintDefaultWelcome()
		// We need to adjust i by 1 since Humans tend to count from one rather than zero
		js.VM.Eval(fmt.Sprintf(`
function MakeWorkbook() {
	return {
		__fname: "",
		__sheets: {},
		accessionInfo: function () {
			var accession_info = {},
				repo_id = 0,
				accession_id = 0;

			if (this.__sheets.Accession === undefined) {
				return {};
			}
			repo_id = parseInt(this.__sheets.Accession[0][1], 10);
			accession_id = parseInt(this.__sheets.Accession[1][1], 10);
			return {
				repo_id: repo_id,
				accession_id: accession_id
			};
		},
		getSheet: function(name) {
			if (this.__sheets[name] === undefined) {
				return {};
			}
			return this.__sheets[name];
		},
		read: function(fname) {
			this.__fname = fname;
			return (this.__sheets = xlsx.read(fname));
		},
		write: function(fname) {
			return xlsx.write(fname, this.__sheets);
		},
		sheetNames: function () {
			return Object.keys(this.__sheets);
		}
	};
};
var Workbook = MakeWorkbook();
Workbook.read(%q);
console.log("Available spreadsheets in 'Workbook' object by name");
console.log("\n  " + Workbook.sheetNames().join("\n  "));
console.log("\n");
`, args[0]))
		js.Repl()
	}

	var (
		xlFile *xlsx.File
	)
	// Read from the given file path
	xlFile, err := xlsx.OpenFile(args[0])
	if err != nil {
		log.Fatalf("Can't open %s, %s", args[0], err)
	}
	fmt.Println(" Available spreadsheets")
	for i, sheet := range xlFile.Sheets {
		fmt.Printf(" %0.2d %s\n", i, sheet.Name)
	}
	for _, record := range resources {
		src, err := xml.Marshal(record)
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Println(src)
	}
}
