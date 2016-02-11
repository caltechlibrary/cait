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
	//"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/blevesearch/bleve"
)

// SearchQuery represents the query options supported by search
type SearchQuery struct {
	// Bleve specific properties
	Explain    bool              `json:"explain,omitempty"`
	FilterTerm map[string]string `json:"filter_term,omitempty"`
	Type       string            `json:"type,omitempty"`

	// Unified search form properties, works for both Basic and Advanced search
	Method   string `json:"method"`
	Action   string `json:"action"`
	AllIDs   bool   `json:"all_ids"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
	// Simple Search
	Q string `json:"q,omitempty"`
	// Advanced Search
	QRequired string `json:"q_required"`
	QExact    string `json:"q_exact"`
	QExcluded string `json:"q_excluded"`

	// Subjects can be a comma delimited list of subjects (e.g. Manuscript Collection, Image Archive)
	Subjects string `json:"q_subjects"`

	// These fields are where we carry search results and request for nav usage
	Total           int    `json:"total"`
	DetailsBaseURI  string `json:"details_base_uri"`
	QueryURLEncoded string
	DetailedResult  NormalizedAccessionView
	Request         *bleve.SearchRequest
	Results         *bleve.SearchResult
}

func uInt64ToInt(u uint64) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d", u))
}

// AttachSearchResults sets the value os the SearchResults field in SearchQuery structs.
func (sq *SearchQuery) AttachSearchResults(sr *bleve.SearchResult) {
	sq.Results = sr
	sq.Total, _ = uInt64ToInt(sr.Total)
	sq.Request = sr.Request

	v := url.Values{}
	if sq.AllIDs == true {
		v.Add("all_ids", "true")
	}
	v.Add("page_size", fmt.Sprintf("%d", sq.PageSize))
	v.Add("page", fmt.Sprintf("%d", sq.Page))
	v.Add("total", fmt.Sprintf("%d", sq.Total))
	v.Add("q", sq.Q)
	v.Add("q_required", sq.QRequired)
	v.Add("q_exact", sq.QExact)
	v.Add("q_excluded", sq.QExcluded)
	sq.QueryURLEncoded = v.Encode()
}

//String return a SearchQuery
func (sq *SearchQuery) String() string {
	return stringify(sq)
}
