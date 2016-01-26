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
	"strconv"
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
	URL        *url.URL `json:"api_url"`
	AuthToken  string   `json:"token,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	DataSet    string   `json:"aspace_dataset,omitempty"`
	Htdocs     string   `json:"htdocs,omitempty"`
	Templates  string   `json:"templates,omitempty"`
	BleveIndex string   `json:"bleve_index,omitempty"`
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

//
// ArchivesApace Models, below are the structures and functions for working
// with a Go representation of the JSONModel types available through the ArchivesSpaceAPI
// REST API. Port based on http://archivesspace.github.io/archivesspace/api/#schemas
//

// Object JSONModel(:object)
type Object map[string]interface{}

// AbstractAgent JSONModel(:abstract_agent)
type AbstractAgent struct {
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	AgentType                 string                   `json:"agent_type,omitempty"`
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []interface{}            `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents"`
	RightsStatements          []*RightsStatement       `json:"rights_statements"`
	SystemGenerated           bool                     `json:"system_generated,omitempty"`
	Notes                     []*NoteText              `json:"notes,omitmepty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	Publish                   bool                     `json:"publish"`
	LockVersion               int                      `json:"lock_version,omitempty"`
	JSONModelType             string                   `json:"jsonmodel_type"`
	CreatedBy                 string                   `json:"created_by,omitempty"`
	LastModifiedBy            string                   `json:"last_modified_by,omitempty"`
	UserMTime                 string                   `json:"user_mtime,omitempty"`
	SystemMTime               string                   `json:"system_mtime,omitempty"`
	CreateTime                string                   `json:"create_time,omitempty"`
	Repository                map[string]interface{}   `json:"repository,omitmepty"`
}

// AbstractAgentRelationship JSONModel(:abstract_agent_relationship)
type AbstractAgentRelationship struct {
	Description    string                 `json:"description,omitempty"`
	Dates          []*Date                `json:"dates"`
	LockVersion    int                    `json:"lock_version,omitempty"`
	JSONModelType  string                 `json:"jsonmodel_type"`
	CreatedBy      string                 `json:"created_by,omitempty"`
	LastModifiedBy string                 `json:"last_modified_by,omitempty"`
	UserMTime      string                 `json:"user_mtime,omitempty"`
	SystemMTime    string                 `json:"system_mtime,omitempty"`
	CreateTime     string                 `json:"create_time,omitempty"`
	Repository     map[string]interface{} `json:"repository,omitmepty"`
}

// AbstractArchivalObject JSONModel(:abstract_archival_object)
type AbstractArchivalObject struct {
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids"`
	Title             string                   `json:"title,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events,omitmepty"`
	Extents           []*Extent                `json:"extents"`
	Dates             []*Date                  `json:"dates"`
	ExternalDocuments []map[string]interface{} `json:"external_documents"`
	RightsStatements  []*RightsStatement       `json:"rights_statements"`
	LinkedAgents      []*Agent                 `json:"linked_agents"`
	Suppressed        bool                     `json:"suppressed"`
	LockVersion       int                      `json:"lock_version,omitempty"`
	JSONModelType     string                   `json:"jsonmodel_type"`
	CreatedBy         string                   `json:"created_by,omitempty"`
	LastModifiedBy    string                   `json:"last_modified_by,omitempty"`
	UserMTime         string                   `json:"user_mtime,omitempty"`
	SystemMTime       string                   `json:"system_mtime,omitempty"`
	CreateTime        string                   `json:"create_time,omitempty"`
	Repository        map[string]interface{}   `json:"repository,omitmepty"`
}

// AbstractClassification JSONModel(:abstract_classification)
type AbstractClassification struct {
	URI            string                 `json:"uri,omitempty"`
	Identifier     string                 `json:"identifier,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Publish        bool                   `json:"publish"` //NOTE: Default value should be true
	PathFromRoot   map[string]interface{} `json:"path_from_root,omitempty"`
	LinkedRecords  map[string]interface{} `json:"linked_records,omitmepty"`
	Creator        map[string]interface{} `json:"creator,omitmepty"`
	LockVersion    int                    `json:"lock_version,omitempty"`
	JSONModelType  string                 `json:"jsonmodel_type"`
	CreatedBy      string                 `json:"created_by,omitempty"`
	LastModifiedBy string                 `json:"last_modified_by,omitempty"`
	UserMTime      string                 `json:"user_mtime,omitempty"`
	SystemMTime    string                 `json:"system_mtime,omitempty"`
	CreateTime     string                 `json:"create_time,omitempty"`
	Repository     map[string]interface{} `json:"repository,omitmepty"`
}

// AbstractName JSONModel(:abstract_name)
type AbstractName struct {
	AuthorityID          string                 `json:"authority_id,omitmepty"`
	Dates                []*Date                `json:"dates"`
	UsaDates             []*Date                `json:"use_dates"`
	Qualifier            string                 `json:"qualifier,omitmepty"`
	Source               string                 `json:"source,omitempty"`
	Rules                string                 `json:"rules,omitempty"`
	Authorized           bool                   `json:"authorized,omitempty"`
	IsDisplayName        bool                   `json:"is_display_name,omitempty"`
	SortName             string                 `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool                   `json:"sort_name_auto_generate,omitempty"`
	LockVersion          int                    `json:"lock_version,omitempty"`
	JSONModelType        string                 `json:"jsonmodel_type"`
	CreatedBy            string                 `json:"created_by,omitempty"`
	LastModifiedBy       string                 `json:"last_modified_by,omitempty"`
	UserMTime            string                 `json:"user_mtime,omitempty"`
	SystemMTime          string                 `json:"system_mtime,omitempty"`
	CreateTime           string                 `json:"create_time,omitempty"`
	Repository           map[string]interface{} `json:"repository,omitmepty"`
}

// AbstractNote JSONModel(:abstract_note)
type AbstractNote struct {
	Label          string                 `json:"label,omitempty"`
	Publish        bool                   `json:"publish"`
	PersistentID   string                 `json:"persistent_id,omitempty"`
	IngestProblem  string                 `json:"ingest_problem,omitmepty"`
	LockVersion    int                    `json:"lock_version,omitempty"`
	JSONModelType  string                 `json:"jsonmodel_type"`
	CreatedBy      string                 `json:"created_by,omitempty"`
	LastModifiedBy string                 `json:"last_modified_by,omitempty"`
	UserMTime      string                 `json:"user_mtime,omitempty"`
	SystemMTime    string                 `json:"system_mtime,omitempty"`
	CreateTime     string                 `json:"create_time,omitempty"`
	Repository     map[string]interface{} `json:"repository,omitmepty"`
}

// Accession JSONModel(:accession)
type Accession struct {
	ID                     int                      `json:"id,omitempty"`
	URI                    string                   `json:"uri,omitempty"`
	ExternalIDs            []*ExternalID            `json:"external_ids"`
	Title                  string                   `json:"title,omitempty"`
	DisplayString          string                   `json:"display_string,omitempty"`
	ID0                    string                   `json:"id_0,omitempty"`
	ID1                    string                   `json:"id_1,omitempty"`
	ID2                    string                   `json:"id_2,omitempty"`
	ID3                    string                   `json:"id_3,omitempty"`
	ContentDescription     string                   `json:"content_description,omitempty"`
	ConditionDescription   string                   `json:"condition_description,omitempty"`
	Disposition            string                   `json:"disposition,omitempty"`
	Inventory              string                   `json:"inventory,omitempty"`
	Provenance             string                   `json:"provenance,omitempty"`
	RelelatedAccessions    []map[string]interface{} `json:"related_accessions,omitempty"`
	AccessionDate          string                   `json:"accession_date,omitempty"`
	Publish                bool                     `json:"publish"`
	Classifications        []map[string]interface{} `json:"classifications,omitempty"`
	Subjects               []map[string]interface{} `json:"subjects"`
	LinkedEvents           []map[string]interface{} `json:"linked_events"`
	Extents                []*Extent                `json:"extents"`
	Dates                  []*Date                  `json:"dates"`
	ExternalDocuments      []map[string]interface{}/**ExternalDocument */ `json:"external_documents"`
	RightsStatements       []*RightsStatement       `json:"rights_statements"`
	Deaccessions           []*Deaccession           `json:"deaccession,omitempty"`
	CollectionManagement   *CollectionManagement    `json:"collection_management,omitempty"`
	UserDefined            *UserDefined             `json:"user_defined,omitempty"`
	RelelatedResources     []map[string]interface{} `json:"related_resources,omitempty"`
	Suppressed             bool                     `json:"suppressed"`
	AcquisitionType        string                   `json:"acquision_type,omitempty"`
	ResourceType           string                   `json:"resource_type,omitempty"`
	RestrictionsApply      bool                     `json:"restrictions_apply,omitempty"`
	RetentionRule          string                   `json:"retention_rule,omitempty"`
	GeneralNote            string                   `json:"general_note,omitempty"`
	AccessRestrictions     bool                     `json:"access_restrictions,omitempty"`
	AccessRestrictionsNote string                   `json:"access_restrictions_note,omitempty"`
	UseRestrictions        bool                     `json:"use_restrictions,omitempty"`
	UseRestrictionsNote    string                   `json:"use_restrictions_note,omitempty"`
	LinkedAgents           []*Agent                 `json:"linked_agents"`
	Instances              []map[string]interface{} `json:"instances,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// AcccessionPartsRelationship JSONModel(:accession_parts_relationship)
type AcccessionPartsRelationship struct {
	Relator     string                 `json:"relator,omitempty"`
	RelatorType string                 `json:"relator_type,omitmepty"`
	Ref         string                 `json:"ref,omitempty"`
	Resolved    map[string]interface{} `json:"_resolved,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// AccessionSiblingRelationship JSONModel(:accession_sibling_relationship)
type AccessionSiblingRelationship struct {
	Relator     string                 `json:"relator,omitempty"`
	RelatorType string                 `json:"relator_type,omitmepty"`
	Ref         string                 `json:"ref,omitempty"`
	Resolved    map[string]interface{} `json:"_resolved,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ActiveEdits JSONModel(:active_edits)
type ActiveEdits struct {
	URI         string                 `json:"uri,omitempty"`
	ActiveEdits map[string]interface{} `json:"active_edits,omitmepty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// AdvancedQuery JSONModel(:advanced_query)
type AdvancedQuery struct {
	Query map[string]interface{} `json:"query,omitempty"` //FIXME, maybe this should be an interface to boolean_query, field_query, data_field_query,boolean_field_query and Object?

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Agent represents an ArchivesSpace complete agent record from the client point of view
type Agent struct {
	ID                        int                      `json:"id,omitempty"`
	Published                 bool                     `json:"publish"`
	AgentType                 string                   `json:"agent_type,omitempty"`
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	Names                     []*NamePerson            `json:"names,omitempty"`
	DisplayName               *NamePerson              `json:"display_name,omitempty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []interface{}            `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents"`
	RightsStatements          []interface{}            `json:"rights_statements"`
	Notes                     []*NoteBiogHist          `json:"notes"`

	LockVersion    int    `json:"lock_version"`
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty"`
	CreateTime     string `json:"create_time,omitempty"`
}

// AgentContact JSONModel(:agent_contact)
type AgentContact struct {
	Name           string       `json:"name,omitempty"`
	Salutation     string       `json:"salutation,omitemtpy"`
	Address1       string       `json:"address_1,omitemtpy"`
	Address2       string       `json:"address_2,omitemtpy"`
	Address3       string       `json:"address_3,omitemtpy"`
	City           string       `json:"city,omitemtpy"`
	Region         string       `json:"region,omitemtpy"`
	Country        string       `json:"country,omitemtpy"`
	PostCode       string       `json:"post_code,omitemtpy"`
	Telephones     []*Telephone `json:"telephones,omitemtpy"`
	Fax            string       `json:"fax,omitemtpy"`
	EMail          string       `json:"email,omitemtpy"`
	EMailSignature string       `json:"email_signature,omitemtpy"`
	Note           string       `json:"note,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// AgentCorporateEntity JSONModel(:agent_corporate_entity)
type AgentCorporateEntity struct {
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitemtpy"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitemtpy"`
	AgentType                 string                   `json:"agent_type,omitemtpy"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitemtpy"`
	LinkedAgentRoles          []string                 `json:"linked_agent_roles,omitemtpy"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitemtpy"`
	RightsStatements          []*RightsStatement       `json:"rights_statements,omitemtpy"`
	SystemGenerated           bool                     `json:"system_generated,omitemtpy"`
	Notes                     string                   `json:"notes,omitemtpy"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitemtpy"`
	Publish                   bool                     `json:"publish,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NameCorporateEntity `json:"names,omitemtpy"`
	DisplayName   *NameCorporateEntity   `json:"display_name,omitemtpy"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitemtpy"`
}

// AgentFamily JSONModel(:agent_family)
type AgentFamily struct {
	URI                       string              `json:"uri,omitempty"`
	Title                     string              `json:"title,omitemtpy"`
	IsLinkedToPublishedRecord bool                `json:"is_linked_to_published_record,omitemtpy"`
	AgentType                 string              `json:"agent_type,omitemtpy"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact     `json:"agent_contacts,omitemtpy"`
	LinkedAgentRoles          []string            `json:"linked_agent_roles,omitemtpy"`
	ExternalDocuments         []*ExternalDocument `json:"external_documents,omitemtpy"`
	RightsStatements          []*RightsStatement  `json:"rights_statements,omitemtpy"`
	SystemGenerated           bool                `json:"system_generated,omitemtpy"`
	Notes                     string              `json:"notes,omitemtpy"`
	DatesOfExistance          []*Date             `json:"dates_of_existence,omitemtpy"`
	Publish                   bool                `json:"publish,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NameFamily          `json:"names,omitemtpy"`
	DisplayName   *NameFamily            `json:"display_name,omitemtpy"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitemtpy"`
}

// AgentPerson JSONModel(:agent_person)
type AgentPerson struct {
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitemtpy"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitemtpy"`
	AgentType                 string                   `json:"agent_type,omitemtpy"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitemtpy"`
	LinkedAgentRoles          []string                 `json:"linked_agent_roles,omitemtpy"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitemtpy"`
	RightsStatements          []*RightsStatement       `json:"rights_statements,omitemtpy"`
	SystemGenerated           bool                     `json:"system_generated,omitemtpy"`
	Notes                     string                   `json:"notes,omitemtpy"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitemtpy"`
	Publish                   bool                     `json:"publish,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NamePerson          `json:"names,omitemtpy"`
	DisplayName   *NamePerson            `json:"display_name,omitemtpy"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitemtpy"`
}

// AgentRelationshipAssociative JSONModel(:agent_relationship_associative)
type AgentRelationshipAssociative struct {
	Description string  `json:"description,omitempty"`
	Dates       []*Date `json:"dates"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Relator  string                 `json:"relator,omitempty"`
	Ref      map[string]interface{} `json:"ref,omitempty"`
	Resolved map[string]interface{} `json:"_resolved,omitempty"`
}

// AgentRelationshipEarlierlater JSONModel(:agent_relationship_earlierlater)
type AgentRelationshipEarlierlater struct {
	Description string  `json:"description,omitempty"`
	Dates       []*Date `json:"dates"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Relator  string                 `json:"relator,omitempty"`
	Ref      map[string]interface{} `json:"ref,omitempty"`
	Resolved map[string]interface{} `json:"_resolved,omitempty"`
}

// AgentRelationshipParentchild JSONModel(:agent_relationship_parentchild)
type AgentRelationshipParentchild struct {
	Description string  `json:"description,omitempty"`
	Dates       []*Date `json:"dates"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Relator  string                 `json:"relator,omitempty"`
	Ref      map[string]interface{} `json:"ref,omitempty"`
	Resolved map[string]interface{} `json:"_resolved,omitempty"`
}

// AgentRelationshipSubordinatesuperior JSONModel(:agent_relationship_subordinatesuperior)
type AgentRelationshipSubordinatesuperior struct {
	Description string  `json:"description,omitempty"`
	Dates       []*Date `json:"dates"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Relator  string                 `json:"relator,omitempty"`
	Ref      map[string]interface{} `json:"ref,omitempty"`
	Resolved map[string]interface{} `json:"_resolved,omitempty"`
}

// AgentSoftware JSONModel(:agent_software)
type AgentSoftware struct {
	URI                       string                   `json:"uri,omitemtpy"`
	Title                     string                   `json:"title,omitemtpy"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitemtpy"`
	AgentType                 string                   `json:"agent_type,omitemtpy"` // ENUM as: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts"`
	LinkedAgentRoles          string                   `json:"linked_agent_roles,omitemtpy"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitemtpy"`
	RightsStatements          []*RightsStatement       `json:"rights_statements"`
	SystemGenerated           bool                     `json:"system_generated,omitempty"`
	Notes                     []*NoteText              `json:"notes,omitmepty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	Publish                   bool                     `json:"publish"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	DisplayName *NameSoftware   `json:"display_name,omitemtpy"`
	Names       []*NameSoftware `json:"names,omitemtpy"`
}

// ArchivalObject JSONModel(:archival_object)
type ArchivalObject struct {
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids"`
	Title             string                   `json:"title,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events,omitmepty"`
	Extents           []*Extent                `json:"extents"`
	ExternalDocuments []map[string]interface{} `json:"external_documents"`
	RightsStatements  []*RightsStatement       `json:"rights_statements"`
	LinkedAgents      []*Agent                 `json:"linked_agents"`
	Suppressed        bool                     `json:"suppressed"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	RefID                    string                 `json:"ref_id,omitemtpy"`
	ConponentID              string                 `json:"component_id,omitemtpy"`
	Level                    string                 `json:"level,omitemtpy"`
	OtherLevel               string                 `json:"other_level,omitemtpy"`
	DisplayString            string                 `json:"display_string,omitemtpy"`
	RestrictionsApply        bool                   `json:"restrictions_apply,omitemtpy"`
	RepositoryProcessingNote string                 `json:"repository_processing_note,omitemtpy"`
	Parent                   map[string]interface{} `json:"parent,omitemtpy"`
	Resource                 map[string]interface{} `json:"resource,omitemtpy"`
	Series                   map[string]interface{} `json:"series,omitemtpy"`
	Position                 int                    `json:"position,omitemtpy"`
	Instances                []*Instance            `json:"instances,omitemtpy"`
	Notes                    []*NoteText            `json:"notes,omitemtpy"`
	HasUnpublishedAncester   bool                   `json:"has_unpublished_ancestor,omitemtpy"`
}

// ArchivalRecordChildren JSONModel(:archival_record_children)
type ArchivalRecordChildren struct {
	Children []*ArchivalObject `json:"children,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// BooleanFieldQuery JSONModel(:boolean_field_query)
type BooleanFieldQuery struct {
	Field string `json:"field,omitemtpy"`
	Value bool   `json:"value,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// BooleanQuery JSONModel(:boolean_query)
type BooleanQuery struct {
	Op string `json:"op,omitemtpy"` // ENUM as: string AND OR NOT
	//FIXME: this needs to be re-thought, do I use an interface type, a struct?
	Subqueries map[string]interface{} `json:"subqueries,omitemtpy"` // One of 	JSONModel(:boolean_query) object,JSONModel(:field_query) object,JSONModel(:boolean_field_query) object,JSONModel(:date_field_query) object

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Classification JSONModel(:classification)
type Classification struct {
	URI           string                 `json:"uri,omitemtpy"`
	Identifier    string                 `json:"identifier,omitemtpy"`
	Title         string                 `json:"title,omitemtpy"`
	Description   string                 `json:"description,omitemtpy"`
	Publish       bool                   `json:"publish,omitemtpy"` //NOTE: default should true
	PathFromRoot  map[string]interface{} `json:"path_from_root,omitemtpy"`
	LinkedRecords map[string]interface{} `json:"linked_records,omitemtpy"`
	Creator       map[string]interface{} `json:"creator,omitemtpy"`

	LockVersion    int    `json:"lock_version"`
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	CreatedBy      string `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string `json:"create_time,omitempty,omitempty"`
}

// ClassificationTerm JSONModel(:classification_term)
type ClassificationTerm struct {
	URI           string                 `json:"uri,omitemtpy"`
	Identifier    string                 `json:"identifier,omitemtpy"`
	Title         string                 `json:"title,omitemtpy"`
	Description   string                 `json:"description,omitemtpy"`
	Publish       bool                   `json:"publish,omitemtpy"` //NOTE: default should true
	PathFromRoot  map[string]interface{} `json:"path_from_root,omitemtpy"`
	LinkedRecords map[string]interface{} `json:"linked_records,omitemtpy"`
	Creator       map[string]interface{} `json:"creator,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Position       int                    `json:"position,omitemtpy"`
	Parent         map[string]interface{} `json:"parent,omitemtpy"`
	Classification map[string]interface{} `json:"classification,omitemtpy"`
}

// RecordTree JSONModel(:record_tree)
type RecordTree struct {
	URI         string `json:"uri,omitemtpy"`
	ID          int    `json:"id,omitemtpy"`
	RecordURI   string `json:"record_uri,omitemtpy"`
	Title       string `json:"title,omitemtpy"`
	Suppressed  bool   `json:"suppressed,omitemtpy"`
	Publish     bool   `json:"publish,omitemtpy"`
	HasChildren bool   `json:"has_children,omitemtpy"`
	NodeType    string `json:"node_type,omitemtpy"`

	LockVersion    int    `json:"lock_version"`
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	CreatedBy      string `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string `json:"create_time,omitempty,omitempty"`
}

// ClassificationTree JSONModel(:classification_tree)
type ClassificationTree struct {
	URI         string `json:"uri,omitemtpy"`
	ID          int    `json:"id,omitemtpy"`
	RecordURI   string `json:"record_uri,omitemtpy"`
	Title       string `json:"title,omitemtpy"`
	Suppressed  bool   `json:"suppressed,omitemtpy"`
	Publish     bool   `json:"publish,omitemtpy"`
	HasChildren bool   `json:"has_children,omitemtpy"`
	NodeType    string `json:"node_type,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Identifier string                `json:"identifier,omitemtpy"`
	Children   []*ClassificationTree `json:"children,omitemtpy"`
}

// CollectionManagement JSONModel(:collection_management)
type CollectionManagement struct {
	URI                            string        `json:"uri,omitempty"`
	ExternalIDs                    []*ExternalID `json:"external_ids"`
	ProcessingHoursPerFootEstimate string        `json:"processing_hours_per_foot_estimate,omitempty"`
	ProcessingTotalExtent          string        `json:"processing_total_extent,omitempty"`
	ProcessingTotalExtentType      string        `json:"processing_total_extent_type,omitempty"`
	ProcessingHoursTotal           string        `json:"processing_hours_total,omitempty"`
	ProcessingPlan                 string        `json:"processing_plan,omitempty"`
	ProcessingPriority             string        `json:"processing_priority,omitempty"`
	ProcessingFundingSource        string        `json:"processing_funding_source,omitempty"`
	Processors                     string        `json:"processors,omitempty"`
	RightsDetermined               bool          `json:"rights_determined,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Container JSONModel(:container)
type Container struct {
	ContainerProfileKey string               `json:"container_profile_key,omitemtpy"`
	Type1               string               `json:"type_1,omitemtpy"`
	Indicator1          string               `json:"indicator_1,omitemtpy"`
	Barcode1            string               `json:"Barcode_1,omitemtpy"`
	Type2               string               `json:"type_2,omitemtpy"`
	Indicator2          string               `json:"indicator_2,omitemtpy"`
	Type3               string               `json:"type_3,omitemtpy"`
	Indicator3          string               `json:"indicator_3"`
	ContainerExtent     string               `json:"container_extent,omitemtpy"`
	ContainerExtentType string               `json:"container_extent_type,omitemtpy"`
	ContainerLocations  []*ContainerLocation `json:"container_locations,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ContainerLocation JSONModel(:container_location)
type ContainerLocation struct {
	Status    string                 `json:"status,omitemtpy"`
	StartDate *Date                  `json:"start_date,omitemtpy"`
	EndDate   *Date                  `json:"end_date,omitemtpy"`
	Note      string                 `json:"note,omitemtpy"`
	Ref       string                 `json:"location,omitemtpy"`
	Resolved  map[string]interface{} `json:"_resolved,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ContainerProfile JSONModel(:container_profile)
type ContainerProfile struct {
	URI             string `json:"uri,omitemtpy"`
	Name            string `json:"name,omitemtpy"`
	URL             string `json:"url,omitemtpy"`
	DimensionUnits  string `json:"dimension_units,omitemtpy"`
	ExtentDimension string `json:"extent_dimension,omitemtpy" ` //ENUM as: height width depth
	Height          string `json:"height,omitemtpy"`
	Width           string `json:"width,omitemtpy"`
	Depth           string `json:"width,omitemtpy"`
	DisplayString   string `json:"display_string,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Date  JSONModel(:date)
type Date struct {
	DateType   string `json:"date_type,omitemtpy"`
	Label      string `json:"label,omitemtpy"`
	Certainty  string `json:"certainty,omitemtpy"`
	Expression string `json:"expression,omitempty"`
	Begin      string `json:"begin,omitempty"`
	End        string `json:"end,omitempty"`
	Era        string `json:"era,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// DateFieldQuery JSONModel(:date_field_query)
type DateFieldQuery struct {
	Comparator string `json:"comparator,omitemtpy"` // ENUM as: greater_than lesser_than equal
	Field      string `json:"field,omitemtpy"`
	Value      *Date  `json:"value,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Deaccession JSONModel(:deaccession)
type Deaccession struct {
	Scope        string    `json:"scope,omitemtpy"`
	Description  string    `json:"description,omitemtpy"`
	Reason       string    `json:"reason,omitemtpy"`
	Disposition  string    `json:"disposition,omitemtpy"`
	Notification bool      `json:"notification,omitemtpy"`
	Date         *Date     `json:"date,omitemtpy"`
	Extents      []*Extent `json:"extents,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// DefaultValues JSONModel(:default_values)
type DefaultValues struct {
	URI        string                 `json:"uri,omitemtpy"`
	RecordType string                 `json:"record_type,omitemtpy"` //ENUM of : archival_object digital_object_component resource accession subject digital_object agent_person agent_family agent_software agent_corporate_entity event location classification classification_term
	Defaults   map[string]interface{} `json:"defaults,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Defaults JSONModel(:defaults)
type Defaults struct {
	ShowSuppressed             bool   `json:"show_suppressed,omitemtpy"`
	Publish                    bool   `json:"publish,omitemtpy"`
	AccessionBrowseColumn1     string `json:"accession_browse_column_1,omitemtpy"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn2     string `json:"accession_browse_column_2,omitemtpy"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn3     string `json:"accession_browse_column_3,omitemtpy"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn4     string `json:"accession_browse_column_4,omitemtpy"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn5     string `json:"accession_browse_column_5,omitemtpy"`      //  enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	ResourceBrowseColumn1      string `json:"resource_browse_column_1,omitemtpy"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn2      string `json:"resource_browse_column_2,omitemtpy"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn3      string `json:"resource_browse_column_3,omitemtpy"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn4      string `json:"resource_browse_column_4,omitemtpy"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn5      string `json:"resource_browse_column_5,omitemtpy"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	DigitalObjectBrowseColumn1 string `json:"digital_object_browse_column_1,omitemtpy"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn2 string `json:"digital_object_browse_column_2,omitemtpy"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn3 string `json:"digital_object_browse_column_3,omitemtpy"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn4 string `json:"digital_object_browse_column_4,omitemtpy"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn5 string `json:"digital_object_browse_column_5,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DefaultValues              bool   `json:"default_values,omitemtpy"`
	NoteOrder                  string `json:"note_order,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// DigitalObject represents a digital object that will eventually become a EAD at COA
type DigitalObject struct {
	ID                int                      `json:"id,omitemtpy"`
	URI               string                   `json:"uri,omitmepty"`
	ExternalIDs       []string                 `json:"external_ids"`
	Title             string                   `json:"title,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events"`
	Extents           []*Extent                `json:"extents"`
	Dates             []*Date                  `json:"dates"`
	ExternalDocuments []map[string]interface{} `json:"external_documents"`
	RightsStatements  []*RightsStatement       `json:"rights_statements"`
	LinkedAgents      []*Agent                 `json:"linked_agents"`
	Suppressed        bool                     `json:"suppressed,omitmepty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	DigitalObjectID      string                   `json:"digital_object_id,omitmepty"`
	Level                string                   `json:"level,omitmepty"`
	DigitalObjectType    string                   `json:"digital_object_type"`
	FileVersions         []*FileVersion           `json:"file_versions,omitemtpy"`
	Restrictions         bool                     `json:"restrictions,omitmepty"`
	Tree                 map[string]interface{}   `json:"tree,omitmepty"`
	Notes                []*NoteText              `json:"notes,omitmepty"`
	CollectionManagement *CollectionManagement    `json:"collection_management,omitempty"`
	UserDefined          []map[string]interface{} `json:"user_defined,omitmepty"`
	LinkedInstances      []map[string]interface{} `json:"linked_instances,omitemtpy"`
}

// DigitalObjectComponent JSONModel(:digital_object_component)
type DigitalObjectComponent struct {
	URI               string                   `json:"uri,omitemtpy"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitemtpy"`
	Title             string                   `json:"title,omitemtpy"`
	Language          string                   `json:"language,omitemtpy"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events,omitemtpy"`
	Extents           []*Extent                `json:"extents,omitemtpy"`
	Dates             []*Date                  `json:"dates,omitemtpy"`
	ExternalDocuments []map[string]interface{} `json:"external_documents,omitemtpy"`
	RightsStatements  []*RightsStatement       `json:"rights_statements,omitemtpy"`
	LinkedAgents      []*Agent                 `json:"linked_agents,omitemtpy"`
	Suppressed        bool                     `json:"suppressed,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	ComponentID            string                 `json:"component_id,omitemtpy"`
	Label                  string                 `json:"label,omitemtpy"`
	DisplayString          string                 `json:"display_string,omitemtpy"`
	FileVersions           []*FileVersion         `json:"file_versions,omitemtpy"`
	Parent                 map[string]interface{} `json:"parent,omitemtpy"`
	DigitalObject          *DigitalObject         `json:"digital_object,omitemtpy"`
	Position               int                    `json:"position,omitemtpy"`
	Notes                  []*NoteText            `json:"notes,omitemtpy"`
	HasUnpublishedAncestor bool                   `json:"has_unpublished_ancestor,omitemtpy"`
}

// DigitalObjectTree JSONModel(:digital_object_tree)
type DigitalObjectTree struct {
	URI         string `json:"uri,omitemtpy"`
	ID          int    `json:"id,omitemtpy"`
	RecordURI   string `json:"record_uri,omitemtpy"`
	Title       string `json:"title,omitemtpy"`
	Suppressed  bool   `json:"suppressed,omitemtpy"`
	Publish     bool   `json:"publish"`
	HasChildren bool   `json:"has_children,omitemtpy"`
	NodeType    string `json:"node_type,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Level             string               `json:"level,omitemtpy"`
	DigitalObjectType string               `json:"digital_object_type,omitemtpy"`
	FileVersions      []*FileVersion       `json:"file_versions,omitemtpy"`
	Children          []*DigitalObjectTree `json:"children,omitemtpy"`
}

// DigitalRecordChildren JSONModel(:digital_record_children)
type DigitalRecordChildren struct {
	Children []*DigitalObjectComponent `json:"children,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Enumeration JSONModel(:enumeration)
type Enumeration struct {
	URI               string              `json:"uri,omitemtpy"`
	Name              string              `json:"name,omitemtpy"`
	DefaultValue      string              `json:"default_value,omitemtpy"`
	Editable          bool                `json:"editable,omitemtpy"`
	Relationships     []string            `json:"relationships,omitemtpy"`
	EnumerationValues []*EnumerationValue `json:"enumeration_values,omitemtpy"`
	Values            []string            `json:"values,omitemtpy"`
	ReadonlyValues    []string            `json:"readonly_values,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// EnumerationMigration JSONModel(:enumeration_migration)
type EnumerationMigration struct {
	URI     string       `json:"uri,omitemtpy"`
	EnumURI *Enumeration `json:"enum_uri,omitemtpy"`
	From    string       `json:"from,omitemtpy"`
	To      string       `json:"to,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// EnumerationValue JSONModel(:enumeration_value)
type EnumerationValue struct {
	URI        string `json:"uri,omitemtpy"`
	Value      string `json:"value,omitemtpy"`
	Position   int    `json:"position,omitemtpy"`
	Suppressed bool   `json:"suppressed,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Event JSONModel(:event)
type Event struct {
	URI               string                   `json:"uri,omitemtpy"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitemtpy"`
	ExternalDocuments []map[string]interface{} `json:"external_documents,omitemtpy"`
	EventType         string                   `json:"event_type,omitemtpy"`
	Date              *Date                    `json:"date,omitemtpy"`
	Timestamp         string                   `json:"timestamp,omitemtpy"`
	Outcome           string                   `json:"outcome,omitemtpy"`
	OutcomeNote       string                   `json:"outcome_note,omitemtpy"`
	Suppressed        bool                     `json:"suppressed,omitemtpy"`
	LinkedAgents      []*Agent                 `json:"linked_agents,omitemtpy"`
	LinkedRecords     map[string]interface{}   `json:"linked_records,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Extent represents an extends json model found in Accession records
type Extent struct {
	Portion          string `json:"portion,omitempty"`
	Number           string `json:"number,omitempty"`
	ExtendType       string `json:"extent_type,omitempty"`
	ContainerSummary string `json:"container_summary,omitempty"`
	PhysicalDetails  string `json:"physical_details,omitempty"`
	Dimensions       string `json:"dimensions,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ExternalDocument a pointer to external documents
type ExternalDocument struct {
	Title    string `json:"title,omitempty"`
	Location string `json:"location,omitempty"`
	Publish  bool   `json:"publish"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ExternalID represents an external ID as found in Accession records
type ExternalID struct {
	ExternalID string `json:"external_id,omitempty"`
	Source     string `json:"source,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// FieldQuery JSONModel(:field_query)
type FieldQuery struct {
	Negated bool   `json:"negated,omitemtpy"`
	Field   string `json:"field,omitemtpy"`
	Value   string `json:"value,omitemtpy"`
	Literal bool   `json:"literal,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// FileVersion JSONModel(:file_version)
type FileVersion struct {
	Identifier            string `json:"identifier,omitemtpy"`
	FileURI               string `json:"file_uri,omitemtpy"`
	Publish               bool   `json:"publish"`
	UseStatement          string `json:"use_statement,omitemtpy"`
	XLinkActuateAttribute string `json:"xlink_actuate_attribute,omitemtpy"`
	XLinkShowAttribute    string `json:"xlink_show_attribute,omitemtpy"`
	FileFormatName        string `json:"file_format_name,omitemtpy"`
	FileFormatVersion     string `json:"file_format_version,omitemtpy"`
	FileSizeBytes         int    `json:"file_size_bytes,omitemtpy"`
	Checksum              string `json:"checksum,omitemtpy"`
	ChecksumMethod        string `json:"checksum_method,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// FindAndReplaceJob JSONModel(:find_and_replace_job)
type FindAndReplaceJob struct {
	Find          string `json:"find,omitemtpy"`
	Replace       string `json:"replace,omitemtpy"`
	RecordType    string `json:"record_type,omitemtpy"`
	Property      string `json:"property,omitemtpy"`
	BaseRecordURI string `json:"base_record_uri,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Group JSONModel(:group)
type Group struct {
	URI               string   `json:"uri,omitemtpy"`
	GroupCode         string   `json:"group_code,omitemtpy"`
	Description       string   `json:"description,omitemtpy"`
	MemberUsernames   []string `json:"member_usernames,omitemtpy"`
	GrantsPermissions []string `json:"grants_permissions,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ImportJob JSONModel(:import_job)
type ImportJob struct {
	Filenames  []string `json:"filenames,omitemtpy"`
	ImportType string   `json:"import_type,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Instance JSONModel(:instance)
type Instance struct {
	InstanceType  string                 `json:"instance_type,omitemtpy"`
	Container     *Container             `json:"container,omitemtpy"`
	SubContainer  *SubContainer          `json:"sub_container,omitemtpy"`
	DigitalObject map[string]interface{} `json:"digital_object,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Job JSONModel(:job)
type Job struct {
	URI           string                 `json:"uri,omitemtpy"`
	JobType       string                 `json:"job_type,omitemtpy"`
	Job           map[string]interface{} `json:"job,omitemtpy"`
	JobParams     string                 `json:"job_params,omitemtpy"`
	TimeSubmitted string                 `json:"time_submitted,omitemtpy"`
	TimeStarted   string                 `json:"time_started,omitemtpy"`
	TimeFinished  string                 `json:"time_finished,omitemtpy"`
	Owner         string                 `json:"owner"`
	Status        string                 `json:"status"` // enum string running completed canceled queued failed default queued
	QueuePosition int                    `json:"queue_position,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Location JSONModel(:location)
type Location struct {
	ID                   int           `json:"id,omitempty"`
	URI                  string        `json:"uri,omitemtpy"`
	Title                string        `json:"title,omitemtpy"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitemtpy"`
	Building             string        `json:"building,omitemtpy"`
	Floor                string        `json:"Floor,omitemtpy"`
	Room                 string        `json:"Room,omitemtpy"`
	Area                 string        `json:"area,omitemtpy"`
	Barcode              string        `json:"barcode,omitemtpy"`
	Classification       `json:"string,omitemtpy"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitemtpy"`
	Coordinate2Label     string `json:"coordinate_2_label,omitemtpy"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitemtpy"`
	Coordinate3Label     string `json:"coordinate_3_label,omitemtpy"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitemtpy"`
	Temporary            string `json:"temporary,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// LocationBatch JSONModel(:location_batch)
type LocationBatch struct {
	URI                  string        `json:"uri,omitemtpy"`
	Title                string        `json:"title,omitemtpy"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitemtpy"`
	Building             string        `json:"building,omitemtpy"`
	Floor                string        `json:"Floor,omitemtpy"`
	Room                 string        `json:"Room,omitemtpy"`
	Area                 string        `json:"area,omitemtpy"`
	Barcode              string        `json:"barcode,omitemtpy"`
	Classification       `json:"string,omitemtpy"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitemtpy"`
	Coordinate2Label     string `json:"coordinate_2_label,omitemtpy"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitemtpy"`
	Coordinate3Label     string `json:"coordinate_3_label,omitemtpy"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitemtpy"`
	Temporary            string `json:"temporary,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Locations        []*Location            `json:"locations,omitemtpy"`
	Coordinate1Range map[string]interface{} `json:"coordinate_1_range,omitemtpy"`
	Coordinate2Range map[string]interface{} `json:"coordinate_2_range,omitemtpy"`
	Coordinate3Range map[string]interface{} `json:"coordinate_3_range,omitemtpy"`
}

// LocationBatchUpdate JSONModel(:location_batch_update)
type LocationBatchUpdate struct {
	URI                  string        `json:"uri,omitemtpy"`
	Title                string        `json:"title,omitemtpy"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitemtpy"`
	Building             string        `json:"building,omitemtpy"`
	Floor                string        `json:"Floor,omitemtpy"`
	Room                 string        `json:"Room,omitemtpy"`
	Area                 string        `json:"area,omitemtpy"`
	Barcode              string        `json:"barcode,omitemtpy"`
	Classification       `json:"string,omitemtpy"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitemtpy"`
	Coordinate2Label     string `json:"coordinate_2_label,omitemtpy"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitemtpy"`
	Coordinate3Label     string `json:"coordinate_3_label,omitemtpy"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitemtpy"`
	Temporary            string `json:"temporary,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	RecordURIs []*Location `json:"record_uris,omitemtpy"`
}

// MergeRequest JSONModel(:merge_request)
type MergeRequest struct {
	URI     string                 `json:"uri,omitemtpy"`
	Target  map[string]interface{} `json:"target,omitemtpy"`
	Victims map[string]interface{} `json:"victims,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NameCorporateEntity JSONModel(:name_corporate_entity)
type NameCorporateEntity struct {
	AuthorityID          string  `json:"authority_id,omitemtpy"`
	Dates                string  `json:"dates,omitemtpy"`
	UseDates             []*Date `json:"use_dates,omitemtpy"`
	Qualifier            string  `json:"qualifier,omitemtpy"`
	Source               string  `json:"source,omitemtpy"`
	Rules                string  `json:"rules,omitemtpy"`
	Authorized           bool    `json:"authorized,omitemtpy"`
	IsDisplayName        bool    `json:"is_display_name,omitemtpy"`
	SortName             string  `json:"sort_name,omitemtpy"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	PrimaryName      string `json:"primary_name,omitemtpy"`
	SubordinateName1 string `json:"subordinate_name_1,omitemtpy"`
	SubordinateName2 string `json:"subordinate_name_2,omitemtpy"`
	Number           string `json:"number,omitemtpy"`
}

// NameFamily JSONModel(:name_family)
type NameFamily struct {
	AuthorityID          string  `json:"authority_id,omitemtpy"`
	Dates                string  `json:"dates,omitemtpy"`
	UseDates             []*Date `json:"use_dates,omitemtpy"`
	Qualifier            string  `json:"qualifier,omitemtpy"`
	Source               string  `json:"source,omitemtpy"`
	Rules                string  `json:"rules,omitemtpy"`
	Authorized           bool    `json:"authorized,omitemtpy"`
	IsDisplayName        bool    `json:"is_display_name,omitemtpy"`
	SortName             string  `json:"sort_name,omitemtpy"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	FamilyName string `json:"family_name,omitemtpy"`
	Prefix     string `json:"prefix,omitemtpy"`
}

// NameForm JSONModel(:name_form)
type NameForm struct {
	URI      string `json:"uri,omitemtpy"`
	Kind     string `json:"kind,omitemtpy"`
	SortName string `json:"sort_name,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NamePerson JSONModel(:name_person)
type NamePerson struct {
	AuthorityID          string  `json:"authority_id,omitemtpy"`
	Dates                string  `json:"dates,omitemtpy"`
	UseDates             []*Date `json:"use_dates,omitemtpy"`
	Qualifier            string  `json:"qualifier,omitemtpy"`
	Source               string  `json:"source,omitemtpy"`
	Rules                string  `json:"rules,omitemtpy"`
	Authorized           bool    `json:"authorized,omitemtpy"`
	IsDisplayName        bool    `json:"is_display_name,omitemtpy"`
	SortName             string  `json:"sort_name,omitemtpy"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitemtpy"` //NOTE: default should be true

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	PrimaryName string `json:"primary_name,omitemtpy"`
	Title       string `json:"title,omitemtpy"`
	NameOrder   string `json:"name_order,omitemtpy"`
	Prefix      string `json:"prefix,omitemtpy"`
	RestOfName  string `json:"rest_of_name,omitemtpy"`
	Suffix      string `json:"suffix,omitemtpy"`
	FullerForm  string `json:"fuller_form,omitemtpy"`
	Number      string `json:"number,omitemtpy"`
}

// NameSoftware JSONModel(:name_software)
type NameSoftware struct {
	AuthorityID          string  `json:"authority_id,omitemtpy"`
	Dates                string  `json:"dates,omitemtpy"`
	UseDates             []*Date `json:"use_dates,omitemtpy"`
	Qualifier            string  `json:"qualifier,omitemtpy"`
	Source               string  `json:"source,omitemtpy"`
	Rules                string  `json:"rules,omitemtpy"`
	Authorized           bool    `json:"authorized,omitemtpy"`
	IsDisplayName        bool    `json:"is_display_name,omitemtpy"`
	SortName             string  `json:"sort_name,omitemtpy"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitemtpy"` //NOTE: default should be true

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	SoftwareName string `json:"software_name,omitemtpy"`
	Version      string `json:"version,omitemtpy"`
	Manufacturer string `json:"manufacturer,omitemtpy"`
}

// NoteAbstract JSONModel(:note_abstract)
type NoteAbstract struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitemtpy"`
}

// NoteBibliography JSONModel(:note_bibliography)
type NoteBibliography struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitemtpy"`
	Type    string   `json:"type,omitemtpy"`
	Items   []string `json:"items,omitemtpy"`
}

// NoteBiogHist JSONModel(:note_bioghist)
type NoteBiogHist struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	SubNotes []*NoteText `json:"subnotes"`
}

// NoteChronology JSONModel(:note_chronology)
type NoteChronology struct {
	Title   string   `json:"title,omitemtpy"`
	Publish bool     `json:"publish"`
	Items   []string `json:"items,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteCitation JSONModel(:note_citation)
type NoteCitation struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string               `json:"content,omitemtpy"`
	XLink   map[string]interface{} `json:"xlink,omitemtpy"`
}

// NoteDefinedlist JSONModel(:note_definedlist)
type NoteDefinedlist struct {
	Title   string   `json:"title,omitemtpy"`
	Publish bool     `json:"publish"`
	Items   []string `json:"items,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteDigitalObject JSONModel(:note_digital_object)
type NoteDigitalObject struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitemtpy"`
	Type    string   `json:"type,omitemtpy"`
}

// NoteIndex JSONModel(:note_index)
type NoteIndex struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string               `json:"content,omitemtpy"`
	Type    string                 `json:"type,omitemtpy"`
	Items   map[string]interface{} `json:"items,omitemtpy"`
}

// NoteIndexItem JSONModel(:note_index_item)
type NoteIndexItem struct {
	Value         string                 `json:"value,omitemtpy"`
	Type          string                 `json:"type,omitemtpy"`
	Reference     string                 `json:"reference,omitemtpy"`
	ReferenceText string                 `json:"reference_text,omitemtpy"`
	ReferenceRef  map[string]interface{} `json:"reference_ref,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteMultipart JSONModel(:note_multipart)
type NoteMultipart struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Type              string             `json:"type,omitemtpy"`
	RightsRestriction *RightsRestriction `json:"rights_restriction,omitemtpy"`
	Subnotes          map[string]interface{}
}

// NoteOrderedlist JSONModel(:note_orderedlist)
type NoteOrderedlist struct {
	Title       string   `json:"title,omitemtpy"`
	Publish     bool     `json:"publish"`
	Enumeration string   `json:"enumeration,omitemtpy"`
	Items       []string `json:"items,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteOutline JSONModel(:note_outline)
type NoteOutline struct {
	Publish bool                `json:"publish"`
	Levels  []*NoteOutlineLevel `json:"levels,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteOutlineLevel JSONModel(:note_outline_level)
type NoteOutlineLevel struct {
	Items map[string]interface{} `json:"items,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// NoteSinglepart JSONModel(:note_singlepart)
type NoteSinglepart struct {
	Label         string `json:"label,omitemtpy"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitemtpy"`
	IngestProblem string `json:"ingest_problem,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitemtpy"`
	Type    string   `json:"type,omitemtpy"`
}

// NoteText JSONModel(:note_text)
type NoteText struct {
	Content string `json:"content,omitempty"`
	Publish bool   `json:"publish"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Permission JSONModel(:permission)
type Permission struct {
	URI            string `json:"uri,omitempty"`
	PermissionCode string `json:"permission_code,omitemtpy"`
	Description    string `json:"description,omitemtpy"`
	Level          string `json:"level,omitemtpy"` // enum string repository global

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Preference JSONModel(:preference)
type Preference struct {
	URI      string    `json:"uri,omitemtpy"`
	UserID   int       `json:"user_id,omitemtpy"`
	Defaults *Defaults `json:"defaults,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// PrintToPDFJob JSONModel(:print_to_pdf_job)
type PrintToPDFJob struct {
	Source string `json:"source,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// RdeTemplate JSONModel(:rde_template)
type RdeTemplate struct {
	URI        string                 `json:"uri,omitemtpy"`
	Name       string                 `json:"name,omitemtpy"`
	RecordType string                 `json:"record_type,omitemtpy"` // enum string archival_object digital_object_component
	Order      []string               `json:"order,omitemtpy"`
	Visible    []string               `json:"visible,omitemtpy"`
	Defaults   map[string]interface{} `json:"defaults,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// ReportJob JSONModel(:report_job)
type ReportJob struct {
	ReportType string `json:"report_type,omitemtpy"`
	Format     string `json:"format,omitemtpy"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Repository represents an ArchivesSpace repository from the client point of view
type Repository struct {
	ID int `json:"id,omitempty"`

	URI                   string                 `json:"uri,omitempty"`
	RepoCode              string                 `json:"repo_code"`
	Name                  string                 `json:"name"`
	OrgCode               string                 `json:"org_code,omitempty"`
	Country               string                 `json:"country,omitempty"`
	ParentInstitutionName string                 `json:"parent_institution_name,omitempty"`
	URL                   string                 `json:"url,omitempty"`
	ImageURL              string                 `json:"image_url,omitempty"`
	ContactPersons        string                 `json:"contact_persons,omitemtpy"`
	AgentRepresentation   map[string]interface{} `json:"agent_representation,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// RepositoryWithAgent JSONModel(:repository_with_agent)
type RepositoryWithAgent struct {
	URI                 string                 `json:"uri,omitempty"`
	Repository          map[string]interface{} `json:"repository,omitempty"`
	AgentRepresentation *AgentCorporateEntity  `json:"agent_representation,omitemtpy"`

	LockVersion    int    `json:"lock_version"`
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	CreatedBy      string `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string `json:"create_time,omitempty,omitempty"`
}

// Resource JSONModel(:resource)
type Resource struct {
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitemtpy"`
	Title             string                   `json:"title,omitemtpy"`
	Language          string                   `json:"language,omitemtpy"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events"`
	Extents           []*Extent                `json:"extents"`
	Dates             *Date                    `json:"dates"`
	ExternalDocuments []map[string]interface{} `json:"external_documents"`
	RightsStatements  *RightsStatement         `json:"rights_statement"`
	LinkedAgents      map[string]interface{}   `json:"linked_agents"`
	Suppressed        bool                     `json:"suppressed"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	ID0                         string                 `json:"id_0,omitempty"`
	ID1                         string                 `json:"id_1,omitempty"`
	ID2                         string                 `json:"id_2,omitempty"`
	ID3                         string                 `json:"id_3,omitempty"`
	Level                       string                 `json:"level,omitempty"`
	OtherLevel                  string                 `json:"other_level,omitempty"`
	ResourceType                string                 `json:"resource_type,omitempty"`
	Tree                        map[string]interface{} `json:"tree,omitempty"`
	Restrictions                bool                   `json:"restrictioons,omitempty"`
	RepositoryProcessingNote    string                 `json:"repository_processing_note,omitempty"`
	EADID                       string                 `json:"ead_id,omitempty"`
	EADLocation                 string                 `json:"ead_location,omitempty"`
	FindingAidTitle             string                 `json:"finding_aid_title,omitempty"`
	FindingAidSubtitle          string                 `json:"finding_aid_subtitle,omitempty"`
	FindingAidFileTitle         string                 `json:"find_aid_filing_title,omitempty"`
	FindingAidDate              string                 `json:"finding_aid_date,omitempty"`
	FindingAidAuthor            string                 `json:"finding_aid_author,omitempty"`
	FindingAidDescriptionRultes string                 `json:"finding_aid_decription_rules,omitempty"`
	FindingAidLanguage          string                 `json:"finding_aid_language,omitempty"`
	FindingAidSponsor           string                 `json:"finding_aid_spansor,omitempty"`
	FindingAidEditionStatement  string                 `json:"finding_aid_edition_statement,omitempty"`
	FindingAidSeriesStatement   string                 `json:"finding_aid_series_statement,omitempty"`
	FindingAidStatus            string                 `json:"finging_aid_status,omitempty"`
	FindingAidNote              string                 `json:"finding_aid_note,omitempty"`
	RevisionStatements          []*RevisionStatement   `json:"revision_statements,omitempty"`
	Instances                   []*Instance            `json:"instances,omitempty"`
	Deaccessions                []*Deaccession         `json:"deaccession,omitempty"`
	CollectionManagement        *CollectionManagement  `json:"collection_management"`
	UserDefined                 *UserDefined           `json:"user_defined,omitempty"`
	ReleatedAccessions          map[string]interface{} `json:"related_accessions,omitempty"`
	Classification              map[string]interface{} `json:"classifications,omitempty"`
	Notes                       map[string]interface{} `json:"notes"`
}

// ResourceTree JSONModel(:resource_tree)
type ResourceTree struct {
	URI         string `json:"uri,omitempty"`
	ID          int    `json:"id,omitempty"`
	RecordURI   string `json:"record_uri,omitempty"`
	Title       string `json:"title,omitempty"`
	Suppressed  bool   `json:"suppressed"`
	Publish     bool   `json:"publish"`
	HasChildren bool   `json:"has_children,omitempty"`
	NodeType    string `json:"node_type,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	FindingAidFilingTitle string                 `json:"finding_aid_filing_title,omitempty"`
	Level                 string                 `json:"level,omitempty"`
	InstanceTypes         []string               `json:"instance_types,omitempty"`
	Containers            map[string]interface{} `json:"containers,omitempty"`
	Children              []*ResourceTree        `json:"children,omitempty"`
}

// RevisionStatement JSONModel(:revision_statement)
type RevisionStatement struct {
	URI         string `json:"uri,omitempty"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// RightsRestriction JSONModel(:rights_restriction)
type RightsRestriction struct {
	Begin                      string                 `json:"begin,omitempty"`
	End                        string                 `json:"end,omitempty"`
	LocalAccessRestrictionType []string               `json:"local_access_restriction_type,omitempty"`
	LinkedRecords              map[string]interface{} `json:"linked_records,omitempty"`
	RestrictionNoteType        string                 `json:"restriction_note_type,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// RightsStatement JSONModel(:rights_statement)
type RightsStatement struct {
	RightsType             string                   `json:"rights_type,omitempty"`
	Identifier             string                   `json:"identifier,omitempty"`
	Active                 bool                     `json:"active,omitempty"`
	Materials              string                   `json:"materials,omitempty"`
	IPStatus               string                   `json:"ip_status,omitempty"`
	IPExpirationDate       *Date                    `json:"ip_expieration_date,omitempty"`
	LicenseIdentifierTerms string                   `json:"license_identifier_terms,omitempty"`
	StatuteCitation        string                   `json:"statute_citation,omitemtpy"`
	Jurisdiction           string                   `json:"jurisdiction,omitempty"`
	TypeNote               string                   `json:"type_note,omitempty"`
	Permissions            string                   `json:"permissions,omitempty"`
	Restrictions           string                   `json:"restrictions"`
	RestrictionStartDate   *Date                    `json:"restrictions_start_date,omitempty"`
	RestrictionEndDate     *Date                    `json:"restriction_end_date,omitempty"`
	GrantedNote            string                   `json:"granted_note,omitempty"`
	ExternalDocuments      []map[string]interface{} `json:"external_documents"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// SubContainer JSONModel(:sub_container)
type SubContainer struct {
	TopContainer  map[string]interface{} `json:"top_container,omitempty"`
	Type2         string                 `json:"type_2,omitempty"`
	Indicator2    string                 `json:"indicator_2,omitempty"`
	Type3         string                 `json:"type_3,omitempty"`
	Indicator3    string                 `json:"indicator_3,omitempty"`
	DisplayString string                 `json:"display_string,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Subject JSONModel(:subject)
type Subject struct {
	ID          int           `json:"id,omitempty"`
	URI         string        `json:"uri,omitempty"`
	Title       string        `json:"title,omitempty"`
	ExternalIDs []*ExternalID `json:"external_ids"`

	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	Publish                   bool                     `json:"publish"`
	Source                    string                   `json:"source,omitempty"`
	ScopeNote                 string                   `json:"scope_note,omitempty"`
	Terms                     []*Term                  `json:"terms,omitempty"` // uri_or_object
	Vocabulary                []map[string]interface{} `json:"vocabularly,omitempty"`
	AuthorityID               string                   `json:"authority_id,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Telephone JSONModel(:telephone)
type Telephone struct {
	URI        string `json:"uri,omitempty"`
	Number     string `json:"number,omitempty"`
	Ext        string `json:"ext,omitempty"`
	NumberType string `json:"number_type"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Term JSONModel(:term)
type Term struct {
	ID         int    `json:"id,omitempty"`
	URI        string `json:"uri,omitempty"`
	Term       string `json:"term,omitempty"`
	TermType   string `json:"term_type,omitempty"`
	Vocabulary string `json:"vocabulary,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// TopContainer JSONModel(:top_container)
type TopContainer struct {
	URI                string                 `json:"uri,omitempty"`
	Indicator          string                 `json:"indicator,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Barcode            string                 `json:"barcode,omitempty"`
	DisplayString      string                 `json:"display_string,omitempty"`
	LongDisplayString  string                 `json:"long_display_string,omitempty"`
	ILSHoldingID       string                 `json:"ils_holding_id,omitempty"`
	ILSItemID          string                 `json:"ils_item_id,omitempty"`
	ExportedToILS      string                 `json:"exported_to_ils,omitempty"`
	Restricted         bool                   `json:"restricted,omitempty"`
	ActiveRestrictions map[string]interface{} `json:"active_restrictions,omitempty"`
	ContainerLocations map[string]interface{} `json:"container_locations,omitempty"`
	ContainerProfile   map[string]interface{} `json:"container_profile,omitempty"`
	Series             map[string]interface{} `json:"series,omitempty"`
	Collection         map[string]interface{} `json:"collection,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

//

// User is a JSONModel used to administer ArchivesSpace
type User struct {
	URI          string                 `json:"uri,omitempty"`
	Username     string                 `json:"username,omitempty"`
	Name         string                 `json:"name,omitempty"`
	IsSystemUser bool                   `json:"is_system_user,omitempty"`
	Permissions  map[string]string      `json:"permissions,omitempty"`
	Groups       map[string]interface{} `json:"groups,omitempty"`
	EMail        string                 `json:"email,omitempty"`
	FirstName    string                 `json:"first_name,omitempty"`
	LastName     string                 `json:"last_name,omitempty"`
	Telephone    string                 `json:"telephone,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Department   string                 `json:"department,omitempty"`
	AgentRecord  map[string]interface{} `json:"agent_record,omitempty"`
	IsAdmin      bool                   `json:"is_admin,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// UserDefined JSONModel(:user_defined)
type UserDefined struct {
	Boolean1 bool   `json:"boolean_1,omitempty"`
	Boolean2 bool   `json:"boolean_2,omitempty"`
	Boolean3 bool   `json:"boolean_3,omitempty"`
	Integer1 string `json:"integer_1,omitempty"`
	Integer2 string `json:"integer_2,omitempty"`
	Integer3 string `json:"integer_3,omitempty"`
	Real1    string `json:"real_1,omitempty"`
	Real2    string `json:"real_2,omitempty"`
	Real3    string `json:"real_3,omitempty"`
	String1  string `json:"string_1,omitempty"`
	String2  string `json:"string_2,omitempty"`
	String3  string `json:"string_3,omitempty"`
	String4  string `json:"string_4,omitempty"`
	Text1    string `json:"text_1,omitempty"`
	Text2    string `json:"text_2,omitempty"`
	Text3    string `json:"text_3,omitempty"`
	Text4    string `json:"text_4,omitempty"`
	Text5    string `json:"text_5,omitempty"`
	Date1    *Date  `json:"date_1,omitempty"`
	Date2    *Date  `json:"date_2,omitempty"`
	Date3    *Date  `json:"date_3,omitempty"`
	Enum1    string `json:"enum_1,omitempty"`
	Enum2    string `json:"enum_2,omitempty"`
	Enum3    string `json:"enum_3,omitempty"`
	Enum4    string `json:"enum_4,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

// Vocabulary JSONModel(:vocabulary)
type Vocabulary struct {
	ID    int     `json:"id,omitempty"`
	URI   string  `json:"uri,omitempty"`
	RefID string  `json:"ref_id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Terms []*Term `json:"terms,omitempty"`

	LockVersion    int               `json:"lock_version"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
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

// String return a Vocabulary
func (vocabulary *Vocabulary) String() string {
	return stringify(vocabulary)
}

// String return a Term
func (term *Term) String() string {
	return stringify(term)
}

// String return a Term
func (obj *DigitalObject) String() string {
	return stringify(obj)
}

// URIToID return an id integer value from a URI for given type.
func URIToID(uri string) int {
	p := strings.LastIndex(uri, "/") + 1
	i, err := strconv.Atoi(uri[p:])
	if err != nil {
		return 0
	}
	return i
}

// URIToRepoID return the repository ID from a URI
func URIToRepoID(uri string) int {
	p := strings.Split(uri, "/")
	if len(p) < 3 || p[1] != "repositories" {
		return 0
	}
	i, err := strconv.Atoi(p[2])
	if err != nil {
		return 0
	}
	return i
}

// URIToVocabularyID return the vocabulary ID from a URI
func URIToVocabularyID(uri string) int {
	p := strings.Split(uri, "/")
	if len(p) < 3 || p[1] != "vocabularies" {
		return 0
	}
	i, err := strconv.Atoi(p[2])
	if err != nil {
		return 0
	}
	return i

}

// IntListToString String from an array of instances
func IntListToString(intList []int, sep string) string {
	stringList := make([]string, len(intList))
	for i := range intList {
		s := fmt.Sprintf("%d", intList[i])
		stringList[i] = s
	}
	return strings.Join(stringList, sep)
}
