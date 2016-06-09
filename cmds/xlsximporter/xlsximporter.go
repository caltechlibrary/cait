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
	"flag"
	"fmt"
	"os"
	"strings"

	// 3rd Party packages
	"github.com/robertkrimen/otto"

	// Caltech Library maintained packages
	"github.com/caltechlibrary/cait"
	"github.com/caltechlibrary/ostdlib"
)

var (
	showHelp      bool
	showVersion   bool
	jsInteractive bool
	jsCallback    string
)

type jsResponse struct {
	Path   string                 `json:"path,omitempty"`
	Object map[string]interface{} `json:"object,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

func usage(exitCode int) {
	fmt.Printf(`

 USAGE: xlsximporter [OPTIONS] [JAVA_SCRIPT_FILENAME|EXCEL_FILE]

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

 In additional to the standard REPL objects and methods a cait api object
 is available with create, update, delete, get and list methods for various
 schema and resources.

 OPTIONS

`)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t(defaults to %s) %s\n", f.Name, f.DefValue, f.Usage)
	})
	fmt.Printf(`

 Examples

 Return a workbook as JONS object of 2d array (rows/cols) with each property
 a name of a spreadsheet

    xlsximporter myfile.xlsx

Read in a JavaScript file, then workbook and then run second JavaScript file.

	xlsximporter preprocess.js myfile.xlsx format-output.js

Read in a JavaScript, then a workbook and drop into a JavaScript REPL

	xlsximporter -i preprocess.js myfile.xlsx

Version %s, repl %s
`, cait.Version, ostdlib.Version)
	os.Exit(exitCode)
}

func init() {
	flag.BoolVar(&showHelp, "h", false, "display this help information")
	flag.BoolVar(&showVersion, "v", false, "display version information")
	flag.BoolVar(&jsInteractive, "i", false, "run a JavaScript Repl after loading spreadsheet")
	flag.StringVar(&jsCallback, "callback", "", "Use this callback name to process spreadsheet")
}

func main() {
	flag.Parse()

	if showHelp == true {
		usage(0)
	}
	if showVersion == true {
		fmt.Printf("Version %s, repl %s\n", cait.Version, ostdlib.Version)
		os.Exit(0)
	}
	if len(os.Args) == 0 && jsInteractive == false {
		usage(1)
	}

	// Initialize the API and setup the JavaScript Environment
	api := cait.New(os.Getenv("CAIT_URL"), os.Getenv("CAIT_USERNAME"), os.Getenv("CAIT_PASSWORD"))
	vm := otto.New()
	//cait.NewJavaScript(api, args)
	js := ostdlib.New(vm)
	// Add basic extensions
	js.AddExtensions()
	// Add general cait extensions
	api.AddExtensions(js)
	// 	// Now add this commands specific extensions
	// 	// jsExtension adds additional JavaScript functionality to process the container
	// 	// Workbook and support pushing it into a resource
	// 	jsExtension := `
	// var Container = {
	// 	isAuth: false,
	// 	config: {},
	// 	Accession: {},
	// 	Resource: {}
	// };
	// Container.getConfiguration = function (workbook) {
	// 	var sheet = workbook.getSheet("Configuration");
	// 	if (!sheet || sheet[1] === undefined) {
	// 		return {
	// 			repoID: 0,
	// 			accessionID: 0,
	// 			resourceID: 0
	// 		};
	// 	}
	// 	this.config = {
	// 		repoID: parseInt(sheet[1][0]) || 0,
	// 		accessionID: parseInt(sheet[1][1]) || 0,
	// 		resourceID: parseInt(sheet[1][2]) || 0
	// 	};
	// 	return this.config;
	// };
	// Container.getAccession = function (repoID, accessionID) {
	// 	if (this.isAuth === undefined || this.isAuth === false) {
	// 		res = api.login();
	// 		this.isAuth = res.isAuth || false;
	// 	}
	// 	if (this.isAuth === true) {
	// 		res = api.getAccession(repoID, accessionID);
	// 		if (res.error !== undefined) {
	// 			return false;
	// 		}
	// 		this.Accession = res;
	// 		return this.Accession;
	// 	}
	// 	return false;
	// };
	// Container.getDigitalObject = function (repoID, digitalObject) {
	// 	if (this.isAuth === undefined || this.isAuth === false) {
	// 		res = api.login();
	// 		this.isAuth = res.isAuth || false;
	// 	}
	// 	if (this.isAuth === true) {
	// 		res = api.createDigitalObject(repoID, digitalObject);
	// 		if (res.error !== undefined) {
	// 			return false;
	// 		}
	// 		this.DigitalObject = res;
	// 		return this.DigitalObject;
	// 	}
	// 	return false;
	// };
	// Container.getDigitalObject = function (repoID, objectID) {
	// 	if (this.isAuth === undefined || this.isAuth === false) {
	// 		res = api.login();
	// 		this.isAuth = res.isAuth || false;
	// 	}
	// 	if (this.isAuth === true) {
	// 		res = api.getDigitalObject(repoID, objectID);
	// 		if (res.error !== undefined) {
	// 			return false;
	// 		}
	// 		this.DigitalObject = res;
	// 		return this.DigitalObject;
	// 	}
	// 	return false;
	// };
	// Container.createResource = function (repoID, resource) {
	// 	if (this.isAuth === undefined || this.isAuth === false) {
	// 		res = api.login();
	// 		this.isAuth = res.isAuth || false;
	// 	}
	// 	if (this.isAuth === true) {
	// 		//FIXME: requires title, digital_object_id as minimal fields
	// 		res = api.createResource(repoID, resource);
	// 		if (res.error !== undefined) {
	// 			return false;
	// 		}
	// 		this.Resource = res;
	// 		return this.Resource;
	// 	}
	// 	return false;
	// };
	// Container.getResource = function (repoID, resourceID) {
	// 	if (this.isAuth === undefined || this.isAuth === false) {
	// 		res = api.login();
	// 		this.isAuth = res.isAuth || false;
	// 	}
	// 	if (this.isAuth === true) {
	// 		res = api.getResource(repoID, resourceID);
	// 		if (res.error !== undefined) {
	// 			return false;
	// 		}
	// 		this.Resource = res;
	// 		return this.Resource;
	// 	}
	// 	return false;
	// };
	// `
	// 	js.Eval(jsExtension)
	// 	js.SetHelp("Container", "getConfiguration", []string{"Workbook object"}, "Read the contents for the 'Configuration' worksheet and return it or an empty object")
	// 	js.SetHelp("Container", "getAccession", []string{"configuration object"}, "With 'Configuration' and fetch the accession from ArchivesSpace")

	// Read in each file listed on the command line then apply to the the
	// JavaScript VM.
	args := flag.Args()
	for _, fname := range args {
		if strings.HasSuffix(fname, ".xlsx") == true {
			js.Eval(fmt.Sprintf("Workbook.read(%q);", fname))
			if jsCallback != "" {
				js.Eval(fmt.Sprintf("(%s(Workbook);)", jsCallback))
			}
		}
		if strings.HasSuffix(fname, ".js") == true {
			js.Run(fname)
		}
	}

	if jsInteractive == true {
		js.AddHelp()
		api.AddHelp(js)
		js.AddAutoComplete()
		js.PrintDefaultWelcome()
		js.Repl()
	} else {
		js.Eval("console.log(JSON.stringify(Workbook));")
	}
}
