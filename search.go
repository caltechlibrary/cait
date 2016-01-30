//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, California Institute of Technology
// All rights reserved.
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

// SearchQuery represents the query options supported by search
type SearchQuery struct {
	JSONModelType string            `json:"json_model_type,omitempty"`
	URI           string            `json:"uri,omitempty"`
	Q             string            `json:"q,omitempty"`
	Page          int               `json:"page,omitempty"`
	PageSize      int               `json:"page_size,omitempty"`
	Explain       bool              `json:"explain,omitempty"`
	FilterTerm    map[string]string `json:"filter_term,omitempty"`
	Type          string            `json:"type,omitempty"`

	//FIXME: some of these I don't understand what their data structure actually are, RSD 2016-01-11
	//RepoID        int    `json:"repo_id,omitempty"`
	//Type          string `json:"type,omitempty"` //NOTE: empty string means search all record types
	//Sort          string `json:"sort,omitempty"`
	//Facet         SearchFacets      `json:"facet,omitemtpy"`
	//SimpleFilter string `json:"simple_filter,omitempty"`
	//Exclude      []int  `json:"exclude,omitempty"`
	//RESTHelpers  bool   `json:"RESTHelpers,omitempty"`
	//RootRecord   string `json:"root_record,omitempty"`
	//IDSet    []int `json:"id_set,omitempty"`
	//AllIDs bool `json:"all_ids,omitempty"`
}

// SearchFacets presents the facets requested in a search request
type SearchFacets map[string]map[string]string

// SearchResults represents the paged results return from a search request
type SearchResults struct {
	FirstPage   int `json:"first_page,omitempty"`
	LastPage    int `json:"last_page,omitempty"`
	ThisPage    int `json:"this_page,omitempty"`
	OffsetFirst int `json:"offset_first,omitempty"`
	OffsetLast  int `json:"offset_last,omitempty"`
	TotalHits   int `json:"total_hits,omitempty"`
	Results     []SearchResult
	Facets      SearchFacets
}

// SearchResult represents individual reult from the dataset returned by a search request
type SearchResult struct {
	ID               int      `json:"id,omitempty"`
	Title            string   `json:"title,omitempty"`
	Types            []string `json:"types,omitempty"`
	JSON             string   `json:"json,omitempty"`
	Suppressed       bool     `json:"suppressed,omitmepty"`
	SystemGenerated  bool     `json:"system_generated,omitempty"`
	Repository       string   `json:"repository,omitempty"`
	SourceEnumS      []string `json:"source_enum_s,omitemtpy"`
	RulesEnumS       []string `json:"rules_enum_s,omitempty"`
	NameOrderEnumS   []string `json:"name_order_enum_s,omitempty"`
	CreatedBy        string   `json:"created_by,omitempty,omitempty"`
	CreateTime       string   `json:"create_time,omitempty,omitempty"`
	SystemMTime      string   `json:"system_mtime,omitempty,omitempty"`
	UserMTime        string   `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy   string   `json:"last_modified_by,omitempty"`
	Source           string   `json:"source,omitempty"`
	Rules            string   `json:"rules,omitempty"`
	LinkedAgentRoles []string `json:"linked_agent_roles,omitempty"`
	URI              string   `json:"uri,omitempty"`
	JSONModelType    string   `json:"json_model_type,omitempty"`
}

//String return a SearchQuery
func (searchquery *SearchQuery) String() string {
	return stringify(searchquery)
}

//String return a SearchFacets
func (searchfacets *SearchFacets) String() string {
	return stringify(searchfacets)
}

//String return a SearchResults
func (searchresults *SearchResults) String() string {
	return stringify(searchresults)
}

//String return a SearchResult
func (searchresult *SearchResult) String() string {
	return stringify(searchresult)
}
