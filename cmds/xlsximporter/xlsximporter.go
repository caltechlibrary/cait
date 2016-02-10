//
// xlsximporter.go is a command line utility design to make it easyer to
// import content from Excel spreadsheets to ArchivesSpace via
// ArchivesSpace REST API. By default it reads an Excel file turning each
// row into a JSON blob.  If you provide a JavaScript file and callback function
// It will use that callback function to generate the resulting JSON blob.
// The JavaScript environment include Getenv() for accessing environment variables
// (e.g. CAIT_API_URL) as well as httpGET(), httpPOST() for interacting with
// the ArchivesSpace REST API directly.
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
	"github.com/rsdoiel/otto"
	"github.com/tealeg/xlsx"
)

var (
	help          bool
	asArray       bool
	inputFilename *string
	jsFilename    *string
	jsCallback    *string
	sheetNo       int
)

type jsResponse struct {
	Path   string                 `json:"path,omitempty"`
	Object map[string]interface{} `json:"object,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

func usage() {
	fmt.Println(`

 USAGE: xlsximporter [OPTIONS] EXCEL_FILENAME

 OVERVIEW

 Read a .xlsx file and return each row as a JSON object (or array of objects).
 If a JavaScript file and callback name are provided then that will be used to
 generate the resulting JSON object per row.

 JAVASCRIPT

 The callback function in JavaScript should return an object that looks like

     {"path": ..., "object": ..., "error": ...}

 The "path" property should contain the desired filename to use for storing
 the JSON blob. If it is empty the output will only be displayed to standard out.

 The "object" property should be the final version of the object. It is what
 will be transformed into a JSON blob.

 The "error" property is a string and if the string is not empty it will be
 used as an error message and cause the processing to stop.

 A simple JavaScript Examples:

    // Counter i is used to name the JSON output files.
    var i = 0;

    // callback is the default name looked for when processing.
    // the command line option -callback lets you used a different name.
    function callback(row) {
        i += 1;
        if (i > 10) {
            // Stop if processing more than 10 rows.
            return {"error": "too many rows..."}
        }
        return {
            "path": "data/" + i + ".json",
            "object": row,
            "error": ""
        }
    }

 Three functions are available to enhance your import. They are

 + os.Getenv() this get an environment variable as a string.
 + http.get() this performs a HTTP GET operation returning content retrieved.
 + http.post() this perform a HTTP POST operation returning content retrieved.
 + api.*() functions are wrappers for the cait API

 os.getEnv() takes one argument, a string, matching the environment variable
 you will to retreive. E.g. Getenv("CAIT_API_URL")

 http.get() accepts a URL (including an parameters), optional headers and
 returns the response body. E.g.

	 // content is the text handed back unaltered from the API call
     content = http.get("http://localhost:8089/repositories?all_ids=true",
		 [{'X-ArchivesSpace-Session':apiToken}]);


 OPTIONS
`)
	flag.PrintDefaults()
	fmt.Println(`

 Examples

    xlsximporter -as-array myfile.xlsx

    xlsximporter -js row2obj.js -callback row2obj myfile.xlsx

`)
	os.Exit(0)
}

func processSheet(sheet *xlsx.Sheet, asArray, jsMap bool, vm *otto.Otto) {
	columnNames := []string{}
	if asArray == true {
		fmt.Println("[")
	}
	for rowNo, row := range sheet.Rows {
		if asArray == true && rowNo > 1 {
			fmt.Printf(", ")
		}
		jsonBlob := make(map[string]string)
		for colNo, cell := range row.Cells {
			if rowNo == 0 {
				columnNames = append(columnNames, cell.String())
			} else {
				// Build a map and render it out
				if colNo >= len(columnNames) {
					k := fmt.Sprintf("column_%d", colNo+1)
					columnNames = append(columnNames, k)
				}
				s := cell.String()
				jsonBlob[columnNames[colNo]] = s
			}
		}
		if rowNo > 0 {
			src, err := json.Marshal(jsonBlob)
			if err != nil {
				log.Fatalf("Can't render JSON blob, %s", err)
			}
			if jsMap == true {
				// We need to eval the callback from inside a closure to be safer
				js := fmt.Sprintf("(function(){ return %s(%s);}())", *jsCallback, src)
				jsValue, err := vm.Eval(js)
				if err != nil {
					log.Fatalf("row: %d, Can't run %s, %s", rowNo, *jsFilename, err)
				}
				response := new(jsResponse)
				err = jsValue.ToStruct(&response)
				if err != nil {
					log.Fatalf("row: %d, do not understand response %s, %s", rowNo, src, err)
				}
				if response.Error != "" {
					log.Fatalf("row: %d, %s", rowNo, response.Error)
				}
				if response.Object == nil {
					log.Fatalf("row: %d, response.object missing, %s", rowNo, src)
				}
				// Now re-package response.Object into a JSON blob
				src, err = json.Marshal(response.Object)
				if err != nil {
					log.Fatalf("row: %d, %s", rowNo, err)
				}
				if response.Path != "" {
					d := path.Dir(response.Path)
					if d != "." {
						os.MkdirAll(d, 0775)
					}
					ioutil.WriteFile(response.Path, src, 0664)
				}
			}
			fmt.Printf("%s\n", src)
		}
	}
	if asArray == true {
		fmt.Println("]")
	}
}

func init() {
	sheetNo = 0
	flag.BoolVar(&help, "h", false, "display this help message")
	flag.BoolVar(&help, "help", false, "display this help message")
	flag.BoolVar(&asArray, "as-array", false, "Write the JSON blobs output as an array")
	flag.IntVar(&sheetNo, "sheet", 0, "Process a specific sheet number, index starts at 1, zero means process all sheets")
	inputFilename = flag.String("i", "", "Read the Excel file from this name")
	jsFilename = flag.String("js", "", "The name of the JavaScript file containing callback function")
	jsCallback = flag.String("callback", "callback", "The name of the JavaScript function to use as a callback")
}

func main() {
	var (
		xlFile   *xlsx.File
		jsSource []byte
	)
	flag.Parse()

	if help == true {
		usage()
	}

	args := flag.Args()
	if len(args) > 0 {
		*inputFilename = args[0]
	}
	if *inputFilename == "" {
		// Read Excel file from standard
		log.Fatalf("Need to provide an xlsx file for input, -i")
	}
	// Read from the given file path
	xlFile, err := xlsx.OpenFile(*inputFilename)
	if err != nil {
		log.Fatalf("Can't open %s, %s", *inputFilename, err)
	}
	api := cait.New(os.Getenv("CAIT_URL"), os.Getenv("CAIT_USERNAME"), os.Getenv("CAIT_PASSWORD"))
	vm := cait.NewJavaScript(api, args)
	jsMap := false
	if *jsFilename != "" {
		jsMap = true
		fname := fmt.Sprintf("%s", *jsFilename)
		jsSource, err = ioutil.ReadFile(fname)
		if err != nil {
			log.Fatalf("Can't read JavaScript file %s, %s", fname, err)
		}
		script, err := vm.Compile(fname, jsSource)
		if err != nil {
			log.Fatalf("JavaScript compile error %s, %s", fname, err)
		}

		// Define any functions, will evaluate each row with vm.Eval()
		_, err = vm.Run(script)
		if err != nil {
			log.Fatalf("Error %s", err)
		}
	}

	// We need to adjust i by 1 since Humans tend to count from 1 rather than zero
	sheetNo = sheetNo - 1
	for i, sheet := range xlFile.Sheets {
		if sheetNo < 0 || sheetNo == i {
			processSheet(sheet, asArray, jsMap, vm)
		}
	}
}
