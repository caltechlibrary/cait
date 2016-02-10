//
// cmds/caitjs/caitjs.go - A command line JavaScript interpreter making the full cait API
// scriptable in JavaScript.
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
	"io/ioutil"
	"log"
	"os"

	"../../../cait"
)

var (
	showHelp    bool
	showVersion bool
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "display this message")
	flag.BoolVar(&showVersion, "v", false, "display version information")
}

func main() {
	flag.Parse()

	if showHelp == true {
		fmt.Println(`
 USAGE: caitjs [OPTIONS] JAVASCRIPT_FILENAME [OPTIONS_PASSED_TO_JAVASCRIPT_FILE]

 OPTIONS

`)
		flag.PrintDefaults()
		fmt.Printf("\nVersion %s\n", cait.Version)
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Printf("Version %s\n", cait.Version)
		os.Exit(0)
	}

	caitAPIURL := os.Getenv("CAIT_API_URL")
	caitUsername := os.Getenv("CAIT_USERNAME")
	caitPassword := os.Getenv("CAIT_PASSWORD")
	if caitAPIURL == "" {
		log.Fatalf("You need to setup your environment vairables to use caitjs.")
	}

	jsArgs := flag.Args()
	jsFilename, jsArgs := jsArgs[0], jsArgs[1:]

	jsSrc, err := ioutil.ReadFile(jsFilename)
	if err != nil {
		log.Fatalf("Can't ready %s, %s", jsFilename, err)
	}
	api := cait.New(caitAPIURL, caitUsername, caitPassword)
	vm := cait.NewJavaScript(api, jsArgs)

	script, err := vm.Compile(jsFilename, jsSrc)
	if err != nil {
		log.Fatalf("Compile error, %s", err)
	}
	_, err := vm.Run(script)
	if err != nil {
		log.Fatal("Runtime error, %s", err)
	}
}
