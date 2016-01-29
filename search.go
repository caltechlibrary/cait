//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2016
// Caltech Library
//
package aspace

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
