//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2015
// Caltech Library
//
package aspace

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

//
// models.go - these are the models implemented in the ArchivesSpace by aspace.go.
// They are a subset of those available from the ArchivesSpace API. It also includes
// simple methods to stringify the models so it is easy to verify visual their contents.
//

// ArchivesSpaceAPI is a struct holding the essentials for communicating
// with the ArchicesSpace REST API
type ArchivesSpaceAPI struct {
	URL       *url.URL `json:"api_url"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	AuthToken string   `json:"token,omitempty"`
}

// ResponseMsg is a structure to hold the JSON portion of a response from the ArchivesSpaceAPI
type ResponseMsg struct {
	Status      string      `json:"status,omitempty"`
	ID          int         `json:"id,omitempty"`
	LockVersion int         `json:"lock_version,omitempty"`
	Stale       interface{} `json:"stale,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Warnings    []string    `json:"warnings,omitempty"`
	Error       interface{} `json:"error,omitempty"`
}

// Repository represents an ArchivesSpace repository from the client point of view
type Repository struct {
	JSONModelType         string                 `json:"json_model_type,omitempty"`
	ID                    int                    `json:"id,omitempty"`
	RepoCode              string                 `json:"repo_code"`
	Name                  string                 `json:"name"`
	URI                   string                 `json:"uri,omitempty"`
	URL                   string                 `json:"url,omitempty"`
	AgentRepresentation   map[string]interface{} `json:"agent_representation,omitempty"`
	Country               string                 `json:"country,omitempty"`
	ImageURL              string                 `json:"image_url,omitempty"`
	OrgCode               string                 `json:"org_code,omitempty"`
	ParentInstitutionName string                 `json:"parent_institution_name,omitempty"`
	LockVersion           int                    `json:"lock_version"`
	CreatedBy             string                 `json:"created_by,omitempty"`
	CreateTime            string                 `json:"create_time,omitempty"`
	SystemMTime           string                 `json:"system_mtime,omitempty"`
	UserMTime             string                 `json:"user_mtime,omitempty"`
}

// Date an ArchivesSpace Date structure
type Date struct {
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	LockVersion    int    `json:"lock_version"`
	Expression     string `json:"expression,omitempty"`
	Begin          string `json:"begin,omitempty"`
	End            string `json:"end,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	CreateTime     string `json:"create_time,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	DateType       string `json:"date_type,omitempty"`
	Label          string `json:"label,omitempty"`
}

// NoteText is the content type of subnotes
type NoteText struct {
	JSONModelType string `json:"jsonmodel_type,omitempty"`
	Content       string `json:"content,omitempty"`
	Publish       bool   `json:"publish,omitempty"`
}

// NoteBiogHist - Notes Biographical Historical
type NoteBiogHist struct {
	JSONModelType string      `json:"jsonmodel_type,omitempty"`
	Label         string      `json:"label,omitempty"`
	PersistentID  string      `json:"persistent_id,omitempty"`
	SubNotes      []*NoteText `json:"subnotes,omitempty"`
	Publish       bool        `json:"publish,omitempty"`
}

// NamePerson a single agent name structure
type NamePerson struct {
	JSONModelType        string  `json:"json_model_type,omitempty"`
	LockVersion          int     `json:"lock_version"`
	PrimaryName          string  `json:"primary_name,omitempty"`
	RestOfName           string  `json:"rest_of_name,omitempty"`
	SortName             string  `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitempty"`
	CreatedBy            string  `json:"created_by,omitempty"`
	CreateTime           string  `json:"create_time,omitempty"`
	SystemMTime          string  `json:"system_mtime,omitempty"`
	UserMTime            string  `json:"user_mtime,omitempty"`
	LastModifiedBy       string  `json:"last_modified_by,omitempty"`
	Authorized           bool    `json:"authorized,omitempty"`
	IsDisplayName        bool    `json:"is_display_name,omitempty"`
	Source               string  `json:"source,omitempty"`
	Rules                string  `json:"rules,omitempty"`
	NameOrder            string  `json:"name_order,omitempty"`
	UseDates             []*Date `json:"use_dates,omitempty"`
}

// User is a JSONModel used to administer ArchivesSpace
type User struct {
	JSONModelType  string                   `json:"json_model_type,omitempty"`
	LockVersion    int                      `json:"lock_version"`
	AgentRecord    map[string]interface{}   `json:"agent_record,omitempty"`
	CreatedBy      string                   `json:"created_by,omitempty"`
	CreateTime     string                   `json:"create_time,omitempty"`
	SystemMTime    string                   `json:"system_mtime,omitempty"`
	UserMTime      string                   `json:"user_mtime,omitempty"`
	LastModifiedBy string                   `json:"last_modified_by,omitempty"`
	Department     string                   `json:"department,omitempty"`
	EMail          string                   `json:"email,omitempty"`
	Name           string                   `json:"name,omitempty"`
	FirstName      string                   `json:"first_name,omitempty"`
	LastName       string                   `json:"last_name,omitempty"`
	Groups         []map[string]interface{} `json:"groups,omitempty"`
	IsAdmin        bool                     `json:"is_admin,omitempty"`
	IsSystemUser   bool                     `json:"is_system_user,omitempty"`
	Permissions    map[string]string        `json:"permissions,omitempty"`
	Telephone      string                   `json:"telephone,omitempty"`
	Title          string                   `json:"title,omitempty"`
	URI            string                   `json:"uri,omitempty"`
}

// AgentContact is a JSONModel for the AgentContacts array/map
type AgentContact struct {
	JSONModelType  string   `json:"json_model_type,omitempty"`
	LockVersion    int      `json:"lock_version"`
	Name           string   `json:"name,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	CreateTime     string   `json:"create_time,omitempty"`
	SystemMTime    string   `json:"system_mtime,omitempty"`
	UserMTime      string   `json:"user_mtime,omitempty"`
	LastModifiedBy string   `json:"last_modified_by,omitempty"`
	Telephones     []string `json:"telephones,omitempty"`
}

// Agent represents an ArchivesSpace complete agent record from the client point of view
type Agent struct {
	JSONModelType             string          `json:"json_model_type,omitempty"`
	LockVersion               int             `json:"lock_version"`
	ID                        int             `json:"id,omitempty"`
	Published                 bool            `json:"publish,omitempty"`
	CreatedBy                 string          `json:"created_by,omitempty"`
	CreateTime                string          `json:"create_time,omitempty"`
	SystemMTime               string          `json:"system_mtime,omitempty"`
	UserMTime                 string          `json:"user_mtime,omitempty"`
	LastModifiedBy            string          `json:"last_modified_by,omitempty"`
	AgentType                 string          `json:"agent_type,omitempty"`
	URI                       string          `json:"uri,omitempty"`
	Title                     string          `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool            `json:"is_linked_to_published_record,omitempty"`
	Names                     []*NamePerson   `json:"names,omitempty"`
	DisplayName               *NamePerson     `json:"display_name,omitempty"`
	RelatedAgents             []interface{}   `json:"related_agents,omitempty"`
	DatesOfExistance          []*Date         `json:"dates_of_existence,omitempty"`
	AgentContacts             []*AgentContact `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []interface{}   `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []interface{}   `json:"external_documents,omitempty"`
	RightsStatements          []interface{}   `json:"rights_statements,omitempty"`
	Notes                     []*NoteBiogHist `json:"notes,omitempty"`
}

// ExternalID represents an external ID as found in Accession records
type ExternalID struct {
	JSONModelType  string `json:"json_model_type,omitempty"`
	ID             string `json:"external_id,omitempty"`
	Source         string `json:"source,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	CreateTime     string `json:"create_time,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
}

// Extent represents an extends json model found in Accession records
type Extent struct {
	JSONModelType    string `json:"json_model_type,omitempty"`
	LockVersion      int    `json:"lock_version"`
	CreatedBy        string `json:"created_by,omitempty"`
	CreateTime       string `json:"create_time,omitempty"`
	SystemMTime      string `json:"system_mtime,omitempty"`
	UserMTime        string `json:"user_mtime,omitempty"`
	LastModifiedBy   string `json:"last_modified_by,omitempty"`
	Number           string `json:"number,omitempty"`
	PhysicalDetails  string `json:"physical_details,omitempty"`
	Portion          string `json:"portion,omitempty"`
	ExtendType       string `json:"extent_type,omitempty"`
	ContainerSummary string `json:"container_summary,omitempty"`
	Dimensions       string `json:"dimensions,omitempty"`
}

// UserDefined struct used in accession records for holding user defined data.
type UserDefined struct {
	JSONModelType  string            `json:"json_model_type,omitempty"`
	LockVersion    int               `json:"lock_version"`
	CreatedBy      string            `json:"created_by,omitempty"`
	CreateTime     string            `json:"create_time,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	Boolean1       bool              `json:"boolean_1,omitempty"`
	Boolean2       bool              `json:"boolean_2,omitempty"`
	Boolean3       bool              `json:"boolean_3,omitempty"`
	Text1          string            `json:"text_1,omitempty"`
	Text2          string            `json:"test_2,omitempty"`
	Text3          string            `json:"text_3,omitempty"`
	Text4          string            `json:"text_4,omitempty"`
	Text5          string            `json:"text_5,omitempty"`
	Integer1       string            `json:"integer_1,omitempty"`
	Integer2       string            `json:"integer_2,omitempty"`
	Integer3       string            `json:"integer_3,omitempty"`
	String1        string            `json:"string_1,omitempty"`
	String2        string            `json:"string_2,omitempty"`
	String3        string            `json:"string_3,omitempty"`
	String4        string            `json:"string_4,omitempty"`
	Enum1          string            `json:"enum_1,omitempty"`
	Enum2          string            `json:"enum_1,omitempty"`
	Enum3          string            `json:"enum_1,omitempty"`
	Enum4          string            `json:"enum_1,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ExternalDocument a pointer to external documents
type ExternalDocument struct {
	Title          string            `json:"title,omitempty"`
	Location       string            `json:"location,omitempty"`
	Publish        bool              `json:"publish,omitempty"`
	Integer        int               `json:"integer,omitempty"`
	JSONModelType  string            `json:"json_model_type,omitempty"`
	LockVersion    int               `json:"lock_version"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// RightsStatement records an Accession Rights' statement in a data structure
type RightsStatement struct {
	JSONModelType     string              `json:"json_model_type,omitempty"`
	LockVersion       int                 `json:"lock_version"`
	Active            bool                `json:"active,omitemtpy"`
	CreatedBy         string              `json:"created_by,omitempty,omitempty"`
	CreateTime        string              `json:"create_time,omitempty,omitempty"`
	SystemMTime       string              `json:"system_mtime,omitempty,omitempty"`
	UserMTime         string              `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy    string              `json:"last_modified_by,omitempty"`
	ExternalDocuments []*ExternalDocument `json:"external_documents,omitempty"`
	Identifier        string              `json:"identifier,omitempty"`
	Restrictions      string              `json:"restrictions,omitempty"`
	RightsType        string              `json:"rights_type,omitempty"`
}

// Deaccession records for Accession
type Deaccession struct {
	JSONModelType  string            `json:"json_model_type,omitempty"`
	LockVersion    int               `json:"lock_version"`
	Active         bool              `json:"active,omitemtpy"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
	Scope          string            `json:"scope,omitempty"`
	Description    string            `json:"description,omitempty"`
	Reason         string            `json:"reason,omitempty"`
	Disposition    string            `json:"disposition,omitempty"`
	Notification   bool              `json:"notification,omitempty"`
	Date           string            `json:"date,omitempty"`
	Extents        []*Extent         `json:"extends,omitempty"`
}

// CollectionManagement records
type CollectionManagement struct {
	URI                            string            `json:"uri,omitempty"`
	ExternalIDs                    []*ExternalID     `json:"external_ids,omitempty"`
	ProcessingHoursPerFootEstimate string            `json:"processing_hours_per_foot_estimate,omitempty"`
	ProcessingTotalExtent          string            `json:"processing_total_extent,omitempty"`
	ProcessingTotalExtentType      string            `json:"processing_total_extent_type,omitempty"`
	ProcessingHoursTotal           string            `json:"processing_hours_total,omitempty"`
	ProcessingPlan                 string            `json:"processing_plan,omitempty"`
	ProcessingPriority             string            `json:"processing_priority,omitempty"`
	ProcessingFundingSource        string            `json:"processing_funding_source,omitempty"`
	Processors                     string            `json:"processors,omitempty"`
	RightsDetermined               bool              `json:"rights_determined,omitempty"`
	JSONModelType                  string            `json:"json_model_type,omitempty"`
	LockVersion                    int               `json:"lock_version"`
	CreatedBy                      string            `json:"created_by,omitempty,omitempty"`
	CreateTime                     string            `json:"create_time,omitempty,omitempty"`
	SystemMTime                    string            `json:"system_mtime,omitempty,omitempty"`
	UserMTime                      string            `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy                 string            `json:"last_modified_by,omitempty"`
	Repository                     map[string]string `json:"repository,omitempty"`
}

// Accession represents an accession record in ArchivesSpace from the client point of view
type Accession struct {
	JSONModelType          string                   `json:"json_model_type,omitempty"`
	LockVersion            int                      `json:"lock_version"`
	CreatedBy              string                   `json:"created_by,omitempty,omitempty"`
	CreateTime             string                   `json:"create_time,omitempty,omitempty"`
	SystemMTime            string                   `json:"system_mtime,omitempty,omitempty"`
	UserMTime              string                   `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy         string                   `json:"last_modified_by,omitempty"`
	ID                     int                      `json:"id,omitempty"`
	Suppressed             bool                     `json:"suppressed"`
	Title                  string                   `json:"title,omitempty"`
	DisplayString          string                   `json:"display_string,omitempty"`
	Publish                bool                     `json:"publish,omitempty"`
	ContentDescription     string                   `json:"content_description,omitempty"`
	ConditionDescription   string                   `json:"condition_description,omitempty"`
	Disposition            string                   `json:"disposition,omitempty"`
	Inventory              string                   `json:"inventory,omitempty"`
	Provenance             string                   `json:"provenance,omitempty"`
	AccessionDate          string                   `json:"accession_date,omitempty"`
	RestrictionsApply      bool                     `json:"restrictions_apply,omitempty"`
	UseRestrictions        bool                     `json:"use_restrictions,omitempty"`
	UseRestrictionsNote    string                   `json:"use_restrictions_note,omitempty"`
	ID0                    string                   `json:"id_0,omitempty"`
	ID1                    string                   `json:"id_1,omitempty"`
	ID2                    string                   `json:"id_2,omitempty"`
	ID3                    string                   `json:"id_3,omitempty"`
	ExternalIDs            []*ExternalID            `json:"external_ids,omitempty"`
	RelelatedAccessions    []map[string]interface{} `json:"related_accessions,omitempty"`
	Classifications        []map[string]interface{} `json:"classifications,omitempty"`
	Subjects               []map[string]interface{} `json:"subjects,omitempty"`
	LinkedEvents           []map[string]interface{} `json:"linked_events,omitempty"`
	Extents                []*Extent                `json:"extents,omitempty"`
	Dates                  []*Date                  `json:"dates,omitempty"`
	ExternalDocuments      []*ExternalDocument      `json:"external_documents,omitempty"`
	RightsStatements       []*RightsStatement       `json:"rights_statements,omitempty"`
	RelelatedResources     []map[string]interface{} `json:"related_resources,omitempty"`
	LinkedAgents           []*Agent                 `json:"linked_agents,omitempty"`
	Instances              []map[string]interface{} `json:"instances,omitempty"`
	URI                    string                   `json:"uri,omitempty"`
	Repository             map[string]string        `json:"repository,omitempty"`
	UserDefined            map[string]interface{}   `json:"user_defined,omitempty"`
	Deaccessions           []*Deaccession           `json:"deaccession,omitempty"`
	CollectionManagement   *CollectionManagement    `json:"collection_management,omitempty"`
	AcquisitionType        string                   `json:"acquision_type,omitempty"`
	ResourceType           string                   `json:"resource_type,omitempty"`
	RetentionRule          string                   `json:"retention_rule,omitempty"`
	GeneralNote            string                   `json:"general_note,omitempty"`
	AccessRestrictions     bool                     `json:"access_restrictions,omitempty"`
	AccessRestrictionsNote string                   `json:"access_restrictions_note,omitempty"`
}

// Vocabulary defines a structure used in both Term and Subject
type Vocabulary struct {
	JSONModelType  string  `json:"json_model_type,omitempty"`
	LockVersion    int     `json:"lock_version"`
	ID             int     `json:"id,omitempty"`
	CreatedBy      string  `json:"created_by,omitempty,omitempty"`
	CreateTime     string  `json:"create_time,omitempty,omitempty"`
	SystemMTime    string  `json:"system_mtime,omitempty,omitempty"`
	UserMTime      string  `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy string  `json:"last_modified_by,omitempty"`
	Name           string  `json:"name,omitempty"`
	RefID          string  `json:"ref_id,omitempty"`
	Terms          []*Term `json:"terms,omitempty"`
	URI            string  `json:"uri,omitempty"`
}

// Term is used in defining a Subject
type Term struct {
	JSONModelType  string      `json:"json_model_type,omitempty"`
	LockVersion    int         `json:"lock_version"`
	ID             int         `json:"id,omitempty"`
	CreatedBy      string      `json:"created_by,omitempty,omitempty"`
	CreateTime     string      `json:"create_time,omitempty,omitempty"`
	SystemMTime    string      `json:"system_mtime,omitempty,omitempty"`
	UserMTime      string      `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy string      `json:"last_modified_by,omitempty"`
	Term           string      `json:"term,omitempty"`
	TermType       string      `json:"term_type,omitempty"`
	URI            string      `json:"uri,omitempty"`
	Vocabulary     *Vocabulary `json:"vocabulary,omitempty"`
}

// Subject represents a subject that can be associated with an accession in a repository
type Subject struct {
	JSONModelType             string                   `json:"json_model_type,omitempty"`
	LockVersion               int                      `json:"lock_version"`
	ID                        int                      `json:"id,omitempty"`
	CreatedBy                 string                   `json:"created_by,omitempty,omitempty"`
	CreateTime                string                   `json:"create_time,omitempty,omitempty"`
	SystemMTime               string                   `json:"system_mtime,omitempty,omitempty"`
	UserMTime                 string                   `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy            string                   `json:"last_modified_by,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitempty"`
	ExternalIDs               []*ExternalID            `json:"external_ids,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record"`
	Publish                   bool                     `json:"publish,omitempty"`
	Source                    string                   `json:"source,omitempty"`
	Terms                     []*Term                  `json:"terms,omitempty"`
	Title                     string                   `json:"title,omitempty"`
	URI                       string                   `json:"uri,omitempty"`
	Vocabulary                *Vocabulary              `json:"vocabulary,omitempty"`
}

// Location represents a item location possible in the archive
type Location struct {
	JSONModelType        string        `json:"json_model_type,omitempty"`
	LockVersion          int           `json:"lock_version"`
	ID                   int           `json:"id,omitempty"`
	CreatedBy            string        `json:"created_by,omitempty,omitempty"`
	CreateTime           string        `json:"create_time,omitempty,omitempty"`
	SystemMTime          string        `json:"system_mtime,omitempty,omitempty"`
	UserMTime            string        `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy       string        `json:"last_modified_by,omitempty"`
	URI                  string        `json:"uri,omitempty"`
	Area                 string        `json:"area"`
	Barcode              string        `json:"barcode"`
	Building             string        `json:"building,omitempty"`
	Classification       string        `json:"classification,omitempty"`
	Coordinate1Indicator string        `json:"coordinate_1_indicator,omitempty"`
	Coordinate1Label     string        `json:"coordinate_1_label,omitempty"`
	Coordinate2Indicator string        `json:"coordinate_2_indicator,omitempty"`
	Coordinate2Label     string        `json:"coordinate_2_label,omitempty"`
	Coordinate3Indicator string        `json:"coordinate_3_indicator,omitempty"`
	Coordinate3Label     string        `json:"coordinate_3_label,omitempty"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitempty"`
	Floor                string        `json:"floor,omitempty"`
	Room                 string        `json:"room,omitempty"`
	Title                string        `json:"title,omitempty"`
}

// SearchQuery represents the query options supported by search
type SearchQuery struct {
	JSONModelType string `json:"json_model_type,omitempty"`
	URI           string `json:"uri,omitempty"`
	Q             string `json:"q,omitempty"`
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"page_size,omitempty"`

	//FIXME: some of these I don't understand what their data structure actually are, RSD 2016-01-11
	//RepoID        int    `json:"repo_id,omitempty"`
	//Type          string `json:"type,omitempty"` //NOTE: empty string means search all record types
	//Sort          string `json:"sort,omitempty"`
	//Facet         SearchFacets      `json:"facet,omitemtpy"`
	//FilterTerm   map[string]string `json:"filter_term,omitempty"`
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
	ID               int   `json:"id,omitempty"`
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
	linkedAgentRoles []string `json:"linked_agent_roles,omitempty"`
	URI              string   `json:"uri,omitempty"`
	JSONModelType    string   `json:"json_model_type,omitempty"`
}

//
// String functions for aspace public structures
//
func stringify(o interface{}) string {
	src, _ := json.Marshal(o)
	return string(src)
}

// String convert NoteText struct as a JSON formatted string
func (aspace *NoteText) String() string {
	return stringify(aspace)
}

// String convert an ArchicesSpaceAPI struct as a JSON formatted string
func (aspace *ArchivesSpaceAPI) String() string {
	return stringify(aspace)
}

// String return a Repository as a JSON formatted string
func (repository *Repository) String() string {
	return stringify(repository)
}

// String return an Agent as a JSON formatted string
func (agent *Agent) String() string {
	return stringify(agent)
}

// String return a ResponseMsg
func (responseMsg *ResponseMsg) String() string {
	return stringify(responseMsg)
}

// String return a UserDefined
func (userDefined *UserDefined) String() string {
	return stringify(userDefined)
}

// String return a ExternalID
func (externalID *ExternalID) String() string {
	return stringify(externalID)
}

// String return an Extent
func (extent *Extent) String() string {
	return stringify(extent)
}

// String return an Accession
func (accession *Accession) String() string {
	return stringify(accession)
}

//String return a Subject
func (subject *Subject) String() string {
	return stringify(subject)
}

//String return a Vocabulary
func (vocabulary *Vocabulary) String() string {
	return stringify(vocabulary)
}

//String return a Term
func (term *Term) String() string {
	return stringify(term)
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

// String from an array of instances
func IntListToString(intList []int, sep string) string {
	stringList := make([]string, len(intList))
	for i := range intList {
		s := fmt.Sprintf("%s", intList[i])
		stringList[i] = s
	}
	return strings.Join(stringList, sep)
}
