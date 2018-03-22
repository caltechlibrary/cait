//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
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
package cait

import (
	"encoding/json"
	"fmt"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
)

func CreateCollection(api *ArchivesSpaceAPI, dname string) (*dataset.Collection, error) {
	// dname is appended to api.Dataset to form the correct collection name
	// e.g. s3://test-dataset.library.example.edu/Archives/repositories/2/accessions
	if api == nil {
		return nil, fmt.Errorf("api not initialized")
	}
	if api.Dataset == "" {
		return nil, fmt.Errorf("api not initialized")
	}
	return dataset.InitCollection(fmt.Sprintf("%s/%s", api.Dataset, dname))
}

func OpenCollection(api *ArchivesSpaceAPI, dname string) (*dataset.Collection, error) {
	// dname is appended to api.Dataset to form the correct collection name
	// e.g. s3://test-dataset.library.example.edu/Archives/repositories/2/accessions
	if api == nil {
		return nil, fmt.Errorf("api not initialized")
	}
	if api.Dataset == "" {
		return nil, fmt.Errorf("api not initialized")
	}
	return dataset.Open(fmt.Sprintf("%s/%s", api.Dataset, dname))
}

func GetKeys(c *dataset.Collection) []string {
	// dir is the name of the collection
	// fname is the key in the collection
	return c.Keys()
}

// ReadJSON read saved JSON file from a dataset collection
func ReadJSON(c *dataset.Collection, fname string) ([]byte, error) {
	// dir is the name of the collection
	// fname is the key in the collection
	return c.ReadJSON(fname)
}

// WriteJSON write out an ArchivesSpace data structure as a JSON file.
func WriteJSON(c *dataset.Collection, fname string, data interface{}) error {
	// dir is the name of the collection
	// fname is the key in the collection
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON(c, %q, data) -> JSON encode, %s", fname, err)
	}
	err = c.CreateJSON(fname, src)
	if err != nil {
		return fmt.Errorf("Could not write JSON data, %s, %s", fname, err)
	}
	return nil
}
