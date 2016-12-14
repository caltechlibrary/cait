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
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

var (
	// TmplMap adds functions for working specifically with ArchivesSpace objects.
	TmplMap = template.FuncMap{
		"digitalObjectLink": func(m map[string]interface{}) string {
			if _, ok := m["digital_objects.title"]; ok == false {
				return ""
			}
			if _, ok := m["digital_objects.file_uris"]; ok == false {
				return ""
			}
			var links []string
			titles := m["digital_objects.title"]
			hrefs := m["digital_objects.file_uris"]
			switch reflect.TypeOf(titles).Kind() {
			case reflect.Slice:
				t := reflect.ValueOf(titles)
				h := reflect.ValueOf(hrefs)
				// Now merge everything into links
				for i := 0; i < t.Len() && i < h.Len(); i++ {
					url := fmt.Sprintf("%s", h.Index(i))
					if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
						anchor := fmt.Sprintf(`<a href="%s">%s</a>`, h.Index(i), t.Index(i))
						if i == 0 {
							links = append(links, anchor)
						}
						if strings.Compare(anchor, links[0]) != 0 {
							links = append(links, anchor)
						}
					}
				}
				return strings.Join(links, " ")
			default:
				url := fmt.Sprintf("%s", m["digital_objects.file_uris"])
				if strings.Contains(url, "://") {
					return fmt.Sprintf(`<a href="%s">%s</a>`, url, m["digital_objects.title"])
				}
			}
			return ""
		},
	}
)
