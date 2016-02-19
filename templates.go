//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
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
	"io/ioutil"
	"strings"
	"text/template"
)

var (
	tmplFuncs = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"arraylength": func(a []string) int {
			return len(a)
		},
		"mapsize": func(m map[string]string) int {
			return len(m)
		},
		"prevPage": func(from, size, max int) int {
			next := from - size
			if next < 0 {
				return 0
			}
			return next
		},
		"nextPage": func(from, size, max int) int {
			next := from + size
			if next > max {
				return from
			}
			return next
		},
		"getType": func(t interface{}) string {
			switch tp := t.(type) {
			default:
				return fmt.Sprintf("%T", tp)
			}
		},
		"asList": func(li []interface{}, sep string) string {
			var l []string
			for _, item := range li {
				l = append(l, fmt.Sprintf("%s", item))
			}
			return strings.Join(l, sep)
		},
		"digitalObjectLink": func(m map[string]interface{}) string {
			var (
				title string
				href  string
			)
			if _, ok := m["digital_objects.title"]; ok == false {
				return ""
			}
			if _, ok := m["digital_objects.file_uris"]; ok == false {
				return ""
			}
			title = fmt.Sprintf("%s", m["digital_objects.title"])
			href = fmt.Sprintf("%s", m["digital_objects.file_uris"])
			return fmt.Sprintf(`<a href="%s">%s</a>`, href, title)
		},
	}
)

// AssembleTemplate support a very simple template setup of an outer HTML file with a content include
// used by caitpage and caitserver
func AssembleTemplate(htmlFilename, includeFilename string) (*template.Template, error) {
	htmlTmpl, err := ioutil.ReadFile(htmlFilename)
	if err != nil {
		return nil, fmt.Errorf("Can't read html template %s, %s", htmlFilename, err)
	}
	includeTmpl, err := ioutil.ReadFile(includeFilename)
	if err != nil {
		return nil, fmt.Errorf("Can't read included template %s, %s", includeFilename, err)
	}
	return template.New(includeFilename).Funcs(tmplFuncs).Parse(fmt.Sprintf(`{{ define "content" }}%s{{ end }}%s`, htmlTmpl, includeTmpl))
}

// Template generate a template struct with functions attach.
func Template(filename string) (*template.Template, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Can't read template %s, %s", filename, err)
	}
	return template.New(filename).Funcs(tmplFuncs).Parse(string(src))
}
