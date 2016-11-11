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
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

//
// models.go - these are the models implemented in the ArchivesSpace by cait.go.
// They are a subset of those available from the ArchivesSpace API. It also includes
// simple methods to stringify the models so it is easy to verify visual their contents.
//

// ArchivesSpaceAPI is a struct holding the essentials for communicating
// with the ArchicesSpace REST API
type ArchivesSpaceAPI struct {
	BaseURL      *url.URL `json:"api_url"`
	CallURL      *url.URL `json:"api_call_url"`
	AuthToken    string   `json:"token,omitempty"`
	Username     string   `json:"username,omitempty"`
	Password     string   `json:"password,omitempty"`
	Dataset      string   `json:"cait_dataset,omitempty"`
	DatasetIndex string   `json:"cait_dataset_index,omitempty"`
	Htdocs       string   `json:"htdocs,omitempty"`
	HtdocsIndex  string   `json:"htdocs_index,omitempty"`
	Templates    string   `json:"templates,omitempty"`
}

// ResponseMsg is a structure to hold the JSON portion of a response from the ArchivesSpaceAPI
type ResponseMsg struct {
	Status      string      `json:"status,omitempty"`
	ID          int         `json:"id,omitempty"`
	LockVersion json.Number `json:"lock_version,Number"`
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
	LockVersion               json.Number              `json:"lock_version,Number"`
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
	LockVersion    json.Number            `json:"lock_version,Number"`
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
	LockVersion       json.Number              `json:"lock_version,Number"`
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
	LockVersion    json.Number            `json:"lock_version,Number"`
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
	LockVersion          json.Number            `json:"lock_version,Number"`
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
	LockVersion    json.Number            `json:"lock_version,Number"`
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
	ID                     int                      `json:"id"`
	URI                    string                   `json:"uri,omitempty"`
	ExternalIDs            []*ExternalID            `json:"external_ids"`
	Title                  string                   `json:"title"`
	DisplayString          string                   `json:"display_string"`
	ID0                    string                   `json:"id_0,omitempty"`
	ID1                    string                   `json:"id_1,omitempty"`
	ID2                    string                   `json:"id_2,omitempty"`
	ID3                    string                   `json:"id_3,omitempty"`
	ContentDescription     string                   `json:"content_description"`
	ConditionDescription   string                   `json:"condition_description"`
	Disposition            string                   `json:"disposition"`
	Inventory              string                   `json:"inventory"`
	Provenance             string                   `json:"provenance"`
	RelatedAccessions      []map[string]interface{} `json:"related_accessions"`
	AccessionDate          string                   `json:"accession_date"`
	Publish                bool                     `json:"publish"`
	Classifications        []map[string]interface{} `json:"classifications"`
	Subjects               []map[string]interface{} `json:"subjects"`
	LinkedEvents           []map[string]interface{} `json:"linked_events"`
	Extents                []*Extent                `json:"extents"`
	Dates                  []*Date                  `json:"dates"`
	ExternalDocuments      []map[string]interface{}/**ExternalDocument */ `json:"external_documents"`
	RightsStatements       []*RightsStatement       `json:"rights_statements"`
	Deaccessions           []*Deaccession           `json:"deaccession,omitempty"`
	CollectionManagement   *CollectionManagement    `json:"collection_management,omitempty"`
	UserDefined            *UserDefined             `json:"user_defined,omitempty"`
	RelatedResources       []map[string]interface{} `json:"related_resources,omitempty"`
	Suppressed             bool                     `json:"suppressed"`
	AcquisitionType        string                   `json:"acquision_type,omitempty"`
	ResourceType           string                   `json:"resource_type"`
	RestrictionsApply      bool                     `json:"restrictions_apply"`
	RetentionRule          string                   `json:"retention_rule,omitempty"`
	GeneralNote            string                   `json:"general_note"`
	AccessRestrictions     bool                     `json:"access_restrictions"`
	AccessRestrictionsNote string                   `json:"access_restrictions_note"`
	UseRestrictions        bool                     `json:"use_restrictions"`
	UseRestrictionsNote    string                   `json:"use_restrictions_note"`

	//	LinkedAgents           []*Agent                 `json:"linked_agents"`

	LinkedAgents []map[string]interface{} `json:"linked_agents"`
	Instances    []map[string]interface{} `json:"instances"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	//	RightsStatements          []*RightsStatement       `json:"rights_statements"`
	RightsStatements []interface{}   `json:"rights_statements"`
	Notes            []*NoteBiogHist `json:"notes"`

	LockVersion    json.Number `json:"lock_version,Number"`
	JSONModelType  string      `json:"jsonmodel_type,omitempty"`
	CreatedBy      string      `json:"created_by,omitempty"`
	LastModifiedBy string      `json:"last_modified_by,omitempty"`
	UserMTime      string      `json:"user_mtime,omitempty"`
	SystemMTime    string      `json:"system_mtime,omitempty"`
	CreateTime     string      `json:"create_time,omitempty"`
}

// AgentContact JSONModel(:agent_contact)
type AgentContact struct {
	Name           string       `json:"name,omitempty"`
	Salutation     string       `json:"salutation,omitempty"`
	Address1       string       `json:"address_1,omitempty"`
	Address2       string       `json:"address_2,omitempty"`
	Address3       string       `json:"address_3,omitempty"`
	City           string       `json:"city,omitempty"`
	Region         string       `json:"region,omitempty"`
	Country        string       `json:"country,omitempty"`
	PostCode       string       `json:"post_code,omitempty"`
	Telephones     []*Telephone `json:"telephones,omitempty"`
	Fax            string       `json:"fax,omitempty"`
	EMail          string       `json:"email,omitempty"`
	EMailSignature string       `json:"email_signature,omitempty"`
	Note           string       `json:"note,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Title                     string                   `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	AgentType                 string                   `json:"agent_type,omitempty"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []string                 `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitempty"`
	RightsStatements          []*RightsStatement       `json:"rights_statements,omitempty"`
	SystemGenerated           bool                     `json:"system_generated,omitempty"`
	Notes                     string                   `json:"notes,omitempty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	Publish                   bool                     `json:"publish,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NameCorporateEntity `json:"names,omitempty"`
	DisplayName   *NameCorporateEntity   `json:"display_name,omitempty"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitempty"`
}

// AgentFamily JSONModel(:agent_family)
type AgentFamily struct {
	URI                       string              `json:"uri,omitempty"`
	Title                     string              `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                `json:"is_linked_to_published_record,omitempty"`
	AgentType                 string              `json:"agent_type,omitempty"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact     `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []string            `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []*ExternalDocument `json:"external_documents,omitempty"`
	RightsStatements          []*RightsStatement  `json:"rights_statements,omitempty"`
	SystemGenerated           bool                `json:"system_generated,omitempty"`
	Notes                     string              `json:"notes,omitempty"`
	DatesOfExistance          []*Date             `json:"dates_of_existence,omitempty"`
	Publish                   bool                `json:"publish,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NameFamily          `json:"names,omitempty"`
	DisplayName   *NameFamily            `json:"display_name,omitempty"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitempty"`
}

// AgentPerson JSONModel(:agent_person)
type AgentPerson struct {
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	AgentType                 string                   `json:"agent_type,omitempty"` //Enum: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []string                 `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitempty"`
	RightsStatements          []*RightsStatement       `json:"rights_statements,omitempty"`
	SystemGenerated           bool                     `json:"system_generated,omitempty"`
	Notes                     string                   `json:"notes,omitempty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	Publish                   bool                     `json:"publish,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Names         []*NamePerson          `json:"names,omitempty"`
	DisplayName   *NamePerson            `json:"display_name,omitempty"`
	RelatedAgents map[string]interface{} `json:"related_agents,omitempty"`
}

// AgentRelationshipAssociative JSONModel(:agent_relationship_associative)
type AgentRelationshipAssociative struct {
	Description string  `json:"description,omitempty"`
	Dates       []*Date `json:"dates"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI                       string                   `json:"uri,omitempty"`
	Title                     string                   `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool                     `json:"is_linked_to_published_record,omitempty"`
	AgentType                 string                   `json:"agent_type,omitempty"` // ENUM as: agent_person agent_corporate_entity agent_software agent_family user
	AgentContacts             []*AgentContact          `json:"agent_contacts"`
	LinkedAgentRoles          string                   `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents,omitempty"`
	RightsStatements          []*RightsStatement       `json:"rights_statements"`
	SystemGenerated           bool                     `json:"system_generated,omitempty"`
	Notes                     []*NoteText              `json:"notes,omitmepty"`
	DatesOfExistance          []*Date                  `json:"dates_of_existence,omitempty"`
	Publish                   bool                     `json:"publish"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	DisplayName *NameSoftware   `json:"display_name,omitempty"`
	Names       []*NameSoftware `json:"names,omitempty"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	RefID                    string                 `json:"ref_id,omitempty"`
	ConponentID              string                 `json:"component_id,omitempty"`
	Level                    string                 `json:"level,omitempty"`
	OtherLevel               string                 `json:"other_level,omitempty"`
	DisplayString            string                 `json:"display_string,omitempty"`
	RestrictionsApply        bool                   `json:"restrictions_apply,omitempty"`
	RepositoryProcessingNote string                 `json:"repository_processing_note,omitempty"`
	Parent                   map[string]interface{} `json:"parent,omitempty"`
	Resource                 map[string]interface{} `json:"resource,omitempty"`
	Series                   map[string]interface{} `json:"series,omitempty"`
	Position                 int                    `json:"position,omitempty"`
	Instances                []*Instance            `json:"instances,omitempty"`
	Notes                    []*NoteText            `json:"notes,omitempty"`
	HasUnpublishedAncester   bool                   `json:"has_unpublished_ancestor,omitempty"`
}

// ArchivalRecordChildren JSONModel(:archival_record_children)
type ArchivalRecordChildren struct {
	Children []*ArchivalObject `json:"children,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Field string `json:"field,omitempty"`
	Value bool   `json:"value,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Op string `json:"op,omitempty"` // ENUM as: string AND OR NOT
	//FIXME: this needs to be re-thought, do I use an interface type, a struct?
	Subqueries map[string]interface{} `json:"subqueries,omitempty"` // One of 	JSONModel(:boolean_query) object,JSONModel(:field_query) object,JSONModel(:boolean_field_query) object,JSONModel(:date_field_query) object

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI           string                 `json:"uri,omitempty"`
	Identifier    string                 `json:"identifier,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Publish       bool                   `json:"publish,omitempty"` //NOTE: default should true
	PathFromRoot  map[string]interface{} `json:"path_from_root,omitempty"`
	LinkedRecords map[string]interface{} `json:"linked_records,omitempty"`
	Creator       map[string]interface{} `json:"creator,omitempty"`

	LockVersion    json.Number `json:"lock_version,Number"`
	JSONModelType  string      `json:"jsonmodel_type,omitempty"`
	CreatedBy      string      `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string      `json:"last_modified_by,omitempty"`
	UserMTime      string      `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string      `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string      `json:"create_time,omitempty,omitempty"`
}

// ClassificationTerm JSONModel(:classification_term)
type ClassificationTerm struct {
	URI           string                 `json:"uri,omitempty"`
	Identifier    string                 `json:"identifier,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Publish       bool                   `json:"publish,omitempty"` //NOTE: default should true
	PathFromRoot  map[string]interface{} `json:"path_from_root,omitempty"`
	LinkedRecords map[string]interface{} `json:"linked_records,omitempty"`
	Creator       map[string]interface{} `json:"creator,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Position       int                    `json:"position,omitempty"`
	Parent         map[string]interface{} `json:"parent,omitempty"`
	Classification map[string]interface{} `json:"classification,omitempty"`
}

// RecordTree JSONModel(:record_tree)
type RecordTree struct {
	URI         string `json:"uri,omitempty"`
	ID          int    `json:"id,omitempty"`
	RecordURI   string `json:"record_uri,omitempty"`
	Title       string `json:"title,omitempty"`
	Suppressed  bool   `json:"suppressed,omitempty"`
	Publish     bool   `json:"publish,omitempty"`
	HasChildren bool   `json:"has_children,omitempty"`
	NodeType    string `json:"node_type,omitempty"`

	LockVersion    json.Number `json:"lock_version,Number"`
	JSONModelType  string      `json:"jsonmodel_type,omitempty"`
	CreatedBy      string      `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string      `json:"last_modified_by,omitempty"`
	UserMTime      string      `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string      `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string      `json:"create_time,omitempty,omitempty"`
}

// ClassificationTree JSONModel(:classification_tree)
type ClassificationTree struct {
	URI         string `json:"uri,omitempty"`
	ID          int    `json:"id,omitempty"`
	RecordURI   string `json:"record_uri,omitempty"`
	Title       string `json:"title,omitempty"`
	Suppressed  bool   `json:"suppressed,omitempty"`
	Publish     bool   `json:"publish,omitempty"`
	HasChildren bool   `json:"has_children,omitempty"`
	NodeType    string `json:"node_type,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Identifier string                `json:"identifier,omitempty"`
	Children   []*ClassificationTree `json:"children,omitempty"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ContainerProfileKey string               `json:"container_profile_key,omitempty"`
	Type1               string               `json:"type_1,omitempty"`
	Indicator1          string               `json:"indicator_1,omitempty"`
	Barcode1            string               `json:"Barcode_1,omitempty"`
	Type2               string               `json:"type_2,omitempty"`
	Indicator2          string               `json:"indicator_2,omitempty"`
	Type3               string               `json:"type_3,omitempty"`
	Indicator3          string               `json:"indicator_3"`
	ContainerExtent     string               `json:"container_extent,omitempty"`
	ContainerExtentType string               `json:"container_extent_type,omitempty"`
	ContainerLocations  []*ContainerLocation `json:"container_locations,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Status    string                 `json:"status,omitempty"`
	StartDate *Date                  `json:"start_date,omitempty"`
	EndDate   *Date                  `json:"end_date,omitempty"`
	Note      string                 `json:"note,omitempty"`
	Ref       string                 `json:"location,omitempty"`
	Resolved  map[string]interface{} `json:"_resolved,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI             string `json:"uri,omitempty"`
	Name            string `json:"name,omitempty"`
	URL             string `json:"url,omitempty"`
	DimensionUnits  string `json:"dimension_units,omitempty"`
	ExtentDimension string `json:"extent_dimension,omitempty" ` //ENUM as: height width depth
	Height          string `json:"height,omitempty"`
	Width           string `json:"width,omitempty"`
	Depth           string `json:"width,omitempty"`
	DisplayString   string `json:"display_string,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	DateType   string `json:"date_type,omitempty"`
	Label      string `json:"label,omitempty"`
	Certainty  string `json:"certainty,omitempty"`
	Expression string `json:"expression,omitempty"`
	Begin      string `json:"begin,omitempty"`
	End        string `json:"end,omitempty"`
	Era        string `json:"era,omitempty"`
	Calendar   string `json:"calendar,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Comparator string `json:"comparator,omitempty"` // ENUM as: greater_than lesser_than equal
	Field      string `json:"field,omitempty"`
	Value      *Date  `json:"value,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Scope        string    `json:"scope,omitempty"`
	Description  string    `json:"description,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	Disposition  string    `json:"disposition,omitempty"`
	Notification bool      `json:"notification,omitempty"`
	Date         *Date     `json:"date,omitempty"`
	Extents      []*Extent `json:"extents,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI        string                 `json:"uri,omitempty"`
	RecordType string                 `json:"record_type,omitempty"` //ENUM of : archival_object digital_object_component resource accession subject digital_object agent_person agent_family agent_software agent_corporate_entity event location classification classification_term
	Defaults   map[string]interface{} `json:"defaults,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ShowSuppressed             bool   `json:"show_suppressed,omitempty"`
	Publish                    bool   `json:"publish,omitempty"`
	AccessionBrowseColumn1     string `json:"accession_browse_column_1,omitempty"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn2     string `json:"accession_browse_column_2,omitempty"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn3     string `json:"accession_browse_column_3,omitempty"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn4     string `json:"accession_browse_column_4,omitempty"`      // enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	AccessionBrowseColumn5     string `json:"accession_browse_column_5,omitempty"`      //  enum string identifier accession_date acquisition_type resource_type restrictions_apply access_restrictions use_restrictions publish no_value
	ResourceBrowseColumn1      string `json:"resource_browse_column_1,omitempty"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn2      string `json:"resource_browse_column_2,omitempty"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn3      string `json:"resource_browse_column_3,omitempty"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn4      string `json:"resource_browse_column_4,omitempty"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	ResourceBrowseColumn5      string `json:"resource_browse_column_5,omitempty"`       // enum string identifier resource_type level language restrictions ead_id finding_aid_status publish no_value
	DigitalObjectBrowseColumn1 string `json:"digital_object_browse_column_1,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn2 string `json:"digital_object_browse_column_2,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn3 string `json:"digital_object_browse_column_3,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn4 string `json:"digital_object_browse_column_4,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DigitalObjectBrowseColumn5 string `json:"digital_object_browse_column_5,omitempty"` // enum string digital_object_id digital_object_type level restrictions publish no_value
	DefaultValues              bool   `json:"default_values,omitempty"`
	NoteOrder                  string `json:"note_order,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ID                int                      `json:"id,omitempty"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	FileVersions         []*FileVersion           `json:"file_versions,omitempty"`
	Restrictions         bool                     `json:"restrictions,omitmepty"`
	Tree                 map[string]interface{}   `json:"tree,omitmepty"`
	Notes                []map[string]interface{} `json:"notes,omitmepty"`
	CollectionManagement *CollectionManagement    `json:"collection_management,omitempty"`
	UserDefined          []map[string]interface{} `json:"user_defined,omitmepty"`
	LinkedInstances      []map[string]interface{} `json:"linked_instances,omitempty"`
}

// DigitalObjectComponent JSONModel(:digital_object_component)
type DigitalObjectComponent struct {
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitempty"`
	Title             string                   `json:"title,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Publish           bool                     `json:"publish"`
	Subjects          []map[string]interface{} `json:"subjects"`
	LinkedEvents      []map[string]interface{} `json:"linked_events,omitempty"`
	Extents           []*Extent                `json:"extents,omitempty"`
	Dates             []*Date                  `json:"dates,omitempty"`
	ExternalDocuments []map[string]interface{} `json:"external_documents,omitempty"`
	RightsStatements  []*RightsStatement       `json:"rights_statements,omitempty"`
	LinkedAgents      []*Agent                 `json:"linked_agents,omitempty"`
	Suppressed        bool                     `json:"suppressed,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	ComponentID            string                 `json:"component_id,omitempty"`
	Label                  string                 `json:"label,omitempty"`
	DisplayString          string                 `json:"display_string,omitempty"`
	FileVersions           []*FileVersion         `json:"file_versions,omitempty"`
	Parent                 map[string]interface{} `json:"parent,omitempty"`
	DigitalObject          *DigitalObject         `json:"digital_object,omitempty"`
	Position               int                    `json:"position,omitempty"`
	Notes                  []*NoteText            `json:"notes,omitempty"`
	HasUnpublishedAncestor bool                   `json:"has_unpublished_ancestor,omitempty"`
}

// DigitalObjectTree JSONModel(:digital_object_tree)
type DigitalObjectTree struct {
	URI         string `json:"uri,omitempty"`
	ID          int    `json:"id,omitempty"`
	RecordURI   string `json:"record_uri,omitempty"`
	Title       string `json:"title,omitempty"`
	Suppressed  bool   `json:"suppressed,omitempty"`
	Publish     bool   `json:"publish"`
	HasChildren bool   `json:"has_children,omitempty"`
	NodeType    string `json:"node_type,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Level             string               `json:"level,omitempty"`
	DigitalObjectType string               `json:"digital_object_type,omitempty"`
	FileVersions      []*FileVersion       `json:"file_versions,omitempty"`
	Children          []*DigitalObjectTree `json:"children,omitempty"`
}

// DigitalRecordChildren JSONModel(:digital_record_children)
type DigitalRecordChildren struct {
	Children []*DigitalObjectComponent `json:"children,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI               string              `json:"uri,omitempty"`
	Name              string              `json:"name,omitempty"`
	DefaultValue      string              `json:"default_value,omitempty"`
	Editable          bool                `json:"editable,omitempty"`
	Relationships     []string            `json:"relationships,omitempty"`
	EnumerationValues []*EnumerationValue `json:"enumeration_values,omitempty"`
	Values            []string            `json:"values,omitempty"`
	ReadonlyValues    []string            `json:"readonly_values,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI     string       `json:"uri,omitempty"`
	EnumURI *Enumeration `json:"enum_uri,omitempty"`
	From    string       `json:"from,omitempty"`
	To      string       `json:"to,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI        string `json:"uri,omitempty"`
	Value      string `json:"value,omitempty"`
	Position   int    `json:"position,omitempty"`
	Suppressed bool   `json:"suppressed,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitempty"`
	ExternalDocuments []map[string]interface{} `json:"external_documents,omitempty"`
	EventType         string                   `json:"event_type,omitempty"`
	Date              *Date                    `json:"date,omitempty"`
	Timestamp         string                   `json:"timestamp,omitempty"`
	Outcome           string                   `json:"outcome,omitempty"`
	OutcomeNote       string                   `json:"outcome_note,omitempty"`
	Suppressed        bool                     `json:"suppressed,omitempty"`
	LinkedAgents      []*Agent                 `json:"linked_agents,omitempty"`
	LinkedRecords     map[string]interface{}   `json:"linked_records,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Portion          string `json:"portion"`
	Number           string `json:"number"`
	ExtentType       string `json:"extent_type"`
	ContainerSummary string `json:"container_summary"`
	PhysicalDetails  string `json:"physical_details"`
	Dimensions       string `json:"dimensions"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Negated bool   `json:"negated,omitempty"`
	Field   string `json:"field,omitempty"`
	Value   string `json:"value,omitempty"`
	Literal bool   `json:"literal,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Identifier            string `json:"identifier,omitempty"`
	FileURI               string `json:"file_uri,omitempty"`
	Publish               bool   `json:"publish"`
	UseStatement          string `json:"use_statement,omitempty"`
	XLinkActuateAttribute string `json:"xlink_actuate_attribute,omitempty"`
	XLinkShowAttribute    string `json:"xlink_show_attribute,omitempty"`
	FileFormatName        string `json:"file_format_name,omitempty"`
	FileFormatVersion     string `json:"file_format_version,omitempty"`
	FileSizeBytes         int    `json:"file_size_bytes,omitempty"`
	Checksum              string `json:"checksum,omitempty"`
	ChecksumMethod        string `json:"checksum_method,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Find          string `json:"find,omitempty"`
	Replace       string `json:"replace,omitempty"`
	RecordType    string `json:"record_type,omitempty"`
	Property      string `json:"property,omitempty"`
	BaseRecordURI string `json:"base_record_uri,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI               string   `json:"uri,omitempty"`
	GroupCode         string   `json:"group_code,omitempty"`
	Description       string   `json:"description,omitempty"`
	MemberUsernames   []string `json:"member_usernames,omitempty"`
	GrantsPermissions []string `json:"grants_permissions,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Filenames  []string `json:"filenames,omitempty"`
	ImportType string   `json:"import_type,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	InstanceType  string                 `json:"instance_type,omitempty"`
	Container     *Container             `json:"container,omitempty"`
	SubContainer  *SubContainer          `json:"sub_container,omitempty"`
	DigitalObject map[string]interface{} `json:"digital_object,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI           string                 `json:"uri,omitempty"`
	JobType       string                 `json:"job_type,omitempty"`
	Job           map[string]interface{} `json:"job,omitempty"`
	JobParams     string                 `json:"job_params,omitempty"`
	TimeSubmitted string                 `json:"time_submitted,omitempty"`
	TimeStarted   string                 `json:"time_started,omitempty"`
	TimeFinished  string                 `json:"time_finished,omitempty"`
	Owner         string                 `json:"owner"`
	Status        string                 `json:"status"` // enum string running completed canceled queued failed default queued
	QueuePosition int                    `json:"queue_position,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI                  string        `json:"uri,omitempty"`
	Title                string        `json:"title,omitempty"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitempty"`
	Building             string        `json:"building,omitempty"`
	Floor                string        `json:"Floor,omitempty"`
	Room                 string        `json:"Room,omitempty"`
	Area                 string        `json:"area,omitempty"`
	Barcode              string        `json:"barcode,omitempty"`
	Classification       `json:"string,omitempty"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitempty"`
	Coordinate2Label     string `json:"coordinate_2_label,omitempty"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitempty"`
	Coordinate3Label     string `json:"coordinate_3_label,omitempty"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitempty"`
	Temporary            string `json:"temporary,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI                  string        `json:"uri,omitempty"`
	Title                string        `json:"title,omitempty"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitempty"`
	Building             string        `json:"building,omitempty"`
	Floor                string        `json:"Floor,omitempty"`
	Room                 string        `json:"Room,omitempty"`
	Area                 string        `json:"area,omitempty"`
	Barcode              string        `json:"barcode,omitempty"`
	Classification       `json:"string,omitempty"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitempty"`
	Coordinate2Label     string `json:"coordinate_2_label,omitempty"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitempty"`
	Coordinate3Label     string `json:"coordinate_3_label,omitempty"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitempty"`
	Temporary            string `json:"temporary,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Locations        []*Location            `json:"locations,omitempty"`
	Coordinate1Range map[string]interface{} `json:"coordinate_1_range,omitempty"`
	Coordinate2Range map[string]interface{} `json:"coordinate_2_range,omitempty"`
	Coordinate3Range map[string]interface{} `json:"coordinate_3_range,omitempty"`
}

// LocationBatchUpdate JSONModel(:location_batch_update)
type LocationBatchUpdate struct {
	URI                  string        `json:"uri,omitempty"`
	Title                string        `json:"title,omitempty"`
	ExternalIDs          []*ExternalID `json:"external_ids,omitempty"`
	Building             string        `json:"building,omitempty"`
	Floor                string        `json:"Floor,omitempty"`
	Room                 string        `json:"Room,omitempty"`
	Area                 string        `json:"area,omitempty"`
	Barcode              string        `json:"barcode,omitempty"`
	Classification       `json:"string,omitempty"`
	Coordinate1Label     string `json:"coordinatel_1_label"`
	Coordinate1Indicator string `json:"coordinate_1_indicator,omitempty"`
	Coordinate2Label     string `json:"coordinate_2_label,omitempty"`
	Coordinate2Indicator string `json:"coordinate_2_indicator,omitempty"`
	Coordinate3Label     string `json:"coordinate_3_label,omitempty"`
	Coordinate3Indicator string `json:"coordinate_3_indicator,omitempty"`
	Temporary            string `json:"temporary,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	RecordURIs []*Location `json:"record_uris,omitempty"`
}

// MergeRequest JSONModel(:merge_request)
type MergeRequest struct {
	URI     string                 `json:"uri,omitempty"`
	Target  map[string]interface{} `json:"target,omitempty"`
	Victims map[string]interface{} `json:"victims,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	AuthorityID          string  `json:"authority_id,omitempty"`
	Dates                string  `json:"dates,omitempty"`
	UseDates             []*Date `json:"use_dates,omitempty"`
	Qualifier            string  `json:"qualifier,omitempty"`
	Source               string  `json:"source,omitempty"`
	Rules                string  `json:"rules,omitempty"`
	Authorized           bool    `json:"authorized,omitempty"`
	IsDisplayName        bool    `json:"is_display_name,omitempty"`
	SortName             string  `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	PrimaryName      string `json:"primary_name,omitempty"`
	SubordinateName1 string `json:"subordinate_name_1,omitempty"`
	SubordinateName2 string `json:"subordinate_name_2,omitempty"`
	Number           string `json:"number,omitempty"`
}

// NameFamily JSONModel(:name_family)
type NameFamily struct {
	AuthorityID          string  `json:"authority_id,omitempty"`
	Dates                string  `json:"dates,omitempty"`
	UseDates             []*Date `json:"use_dates,omitempty"`
	Qualifier            string  `json:"qualifier,omitempty"`
	Source               string  `json:"source,omitempty"`
	Rules                string  `json:"rules,omitempty"`
	Authorized           bool    `json:"authorized,omitempty"`
	IsDisplayName        bool    `json:"is_display_name,omitempty"`
	SortName             string  `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	FamilyName string `json:"family_name,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
}

// NameForm JSONModel(:name_form)
type NameForm struct {
	URI      string `json:"uri,omitempty"`
	Kind     string `json:"kind,omitempty"`
	SortName string `json:"sort_name,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	AuthorityID          string  `json:"authority_id,omitempty"`
	Dates                string  `json:"dates,omitempty"`
	UseDates             []*Date `json:"use_dates,omitempty"`
	Qualifier            string  `json:"qualifier,omitempty"`
	Source               string  `json:"source,omitempty"`
	Rules                string  `json:"rules,omitempty"`
	Authorized           bool    `json:"authorized,omitempty"`
	IsDisplayName        bool    `json:"is_display_name,omitempty"`
	SortName             string  `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitempty"` //NOTE: default should be true

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	PrimaryName string `json:"primary_name,omitempty"`
	Title       string `json:"title,omitempty"`
	NameOrder   string `json:"name_order,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	RestOfName  string `json:"rest_of_name,omitempty"`
	Suffix      string `json:"suffix,omitempty"`
	FullerForm  string `json:"fuller_form,omitempty"`
	Number      string `json:"number,omitempty"`
}

// NameSoftware JSONModel(:name_software)
type NameSoftware struct {
	AuthorityID          string  `json:"authority_id,omitempty"`
	Dates                string  `json:"dates,omitempty"`
	UseDates             []*Date `json:"use_dates,omitempty"`
	Qualifier            string  `json:"qualifier,omitempty"`
	Source               string  `json:"source,omitempty"`
	Rules                string  `json:"rules,omitempty"`
	Authorized           bool    `json:"authorized,omitempty"`
	IsDisplayName        bool    `json:"is_display_name,omitempty"`
	SortName             string  `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool    `json:"sort_name_auto_generate,omitempty"` //NOTE: default should be true

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	SoftwareName string `json:"software_name,omitempty"`
	Version      string `json:"version,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
}

// NoteAbstract JSONModel(:note_abstract)
type NoteAbstract struct {
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitempty"`
}

// NoteBibliography JSONModel(:note_bibliography)
type NoteBibliography struct {
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitempty"`
	Type    string   `json:"type,omitempty"`
	Items   []string `json:"items,omitempty"`
}

// NoteBiogHist JSONModel(:note_bioghist)
type NoteBiogHist struct {
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Title   string   `json:"title,omitempty"`
	Publish bool     `json:"publish"`
	Items   []string `json:"items,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string               `json:"content,omitempty"`
	XLink   map[string]interface{} `json:"xlink,omitempty"`
}

// NoteDefinedlist JSONModel(:note_definedlist)
type NoteDefinedlist struct {
	Title   string   `json:"title,omitempty"`
	Publish bool     `json:"publish"`
	Items   []string `json:"items,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitempty"`
	Type    string   `json:"type,omitempty"`
}

// NoteIndex JSONModel(:note_index)
type NoteIndex struct {
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string               `json:"content,omitempty"`
	Type    string                 `json:"type,omitempty"`
	Items   map[string]interface{} `json:"items,omitempty"`
}

// NoteIndexItem JSONModel(:note_index_item)
type NoteIndexItem struct {
	Value         string                 `json:"value,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Reference     string                 `json:"reference,omitempty"`
	ReferenceText string                 `json:"reference_text,omitempty"`
	ReferenceRef  map[string]interface{} `json:"reference_ref,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Type              string             `json:"type,omitempty"`
	RightsRestriction *RightsRestriction `json:"rights_restriction,omitempty"`
	Subnotes          map[string]interface{}
}

// NoteOrderedlist JSONModel(:note_orderedlist)
type NoteOrderedlist struct {
	Title       string   `json:"title,omitempty"`
	Publish     bool     `json:"publish"`
	Enumeration string   `json:"enumeration,omitempty"`
	Items       []string `json:"items,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Levels  []*NoteOutlineLevel `json:"levels,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Items map[string]interface{} `json:"items,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Label         string `json:"label,omitempty"`
	Publish       bool   `json:"publish"`
	PersistentID  string `json:"persistent_id,omitempty"`
	IngestProblem string `json:"ingest_problem,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	Content []string `json:"content,omitempty"`
	Type    string   `json:"type,omitempty"`
}

// NoteText JSONModel(:note_text)
type NoteText struct {
	Content string `json:"content,omitempty"`
	Publish bool   `json:"publish"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	PermissionCode string `json:"permission_code,omitempty"`
	Description    string `json:"description,omitempty"`
	Level          string `json:"level,omitempty"` // enum string repository global

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI      string    `json:"uri,omitempty"`
	UserID   int       `json:"user_id,omitempty"`
	Defaults *Defaults `json:"defaults,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Source string `json:"source,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	URI        string                 `json:"uri,omitempty"`
	Name       string                 `json:"name,omitempty"`
	RecordType string                 `json:"record_type,omitempty"` // enum string archival_object digital_object_component
	Order      []string               `json:"order,omitempty"`
	Visible    []string               `json:"visible,omitempty"`
	Defaults   map[string]interface{} `json:"defaults,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ReportType string `json:"report_type,omitempty"`
	Format     string `json:"format,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ContactPersons        string                 `json:"contact_persons,omitempty"`
	AgentRepresentation   map[string]interface{} `json:"agent_representation,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	AgentRepresentation *AgentCorporateEntity  `json:"agent_representation,omitempty"`

	LockVersion    json.Number `json:"lock_version,Number"`
	JSONModelType  string      `json:"jsonmodel_type,omitempty"`
	CreatedBy      string      `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string      `json:"last_modified_by,omitempty"`
	UserMTime      string      `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string      `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string      `json:"create_time,omitempty,omitempty"`
}

// Resource JSONModel(:resource)
type Resource struct {
	ID                int                      `json:"id,omitempty"`
	XMLName           xml.Name                 `json:"-"`
	URI               string                   `json:"uri,omitempty"`
	ExternalIDs       []*ExternalID            `json:"external_ids,omitempty"`
	Title             string                   `json:"title,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Publish           bool                     `json:"publish,omitempty"`
	Subjects          []map[string]interface{} `json:"subjects,omitempty"`
	LinkedEvents      []map[string]interface{} `json:"linked_events,omitempty"`
	Extents           []*Extent                `json:"extents,omitempty"`
	Dates             []*Date                  `json:"dates,omitempty"`
	ExternalDocuments []map[string]interface{} `json:"external_documents,omitempty"`

	//	RightsStatements  []*RightsStatement       `json:"rights_statement"`
	RightsStatements []interface{} `json:"rights_statements,omitempty"`
	LinkedAgents     []*Agent      `json:"linked_agents,ommitempty"`
	Suppressed       bool          `json:"suppressed,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty"`
	CreateTime     string            `json:"create_time,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`

	ID0          string                 `json:"id_0,omitempty"`
	ID1          string                 `json:"id_1,omitempty"`
	ID2          string                 `json:"id_2,omitempty"`
	ID3          string                 `json:"id_3,omitempty"`
	Level        string                 `json:"level,omitempty"`
	OtherLevel   string                 `json:"other_level,omitempty"`
	ResourceType string                 `json:"resource_type,omitempty"`
	Tree         map[string]interface{} `json:"tree,omitempty"`

	Restrictions                bool                     `json:"restrictioons,omitempty"`
	RepositoryProcessingNote    string                   `json:"repository_processing_note,omitempty"`
	EADID                       string                   `xml:"control>recordid" json:"ead_id,omitempty"`
	EADLocation                 string                   `xml:"control>location" json:"ead_location,omitempty"`
	FindingAidTitle             string                   `xml:"control>filedesc>titlestmt>titleproper" json:"finding_aid_title,omitempty"`
	FindingAidSubtitle          string                   `xml:"control>filedesc>titlestmt>subtitle" json:"finding_aid_subtitle,omitempty"`
	FindingAidFileTitle         string                   `xml:"control>filedesc>titlestmt>filing_title" json:"find_aid_filing_title,omitempty"`
	FindingAidDate              string                   `json:"finding_aid_date,omitempty"`
	FindingAidAuthor            string                   `xml:"control>filedesc>titlestmt>author" json:"finding_aid_author,omitempty"`
	FindingAidDescriptionRultes string                   `json:"finding_aid_decription_rules,omitempty"`
	FindingAidLanguage          string                   `json:"finding_aid_language,omitempty"`
	FindingAidSponsor           string                   `xml:"control>filedesc>titlestmt>sponsor" json:"finding_aid_sponsor,omitempty"`
	FindingAidEditionStatement  string                   `json:"finding_aid_edition_statement,omitempty"`
	FindingAidSeriesStatement   string                   `json:"finding_aid_series_statement,omitempty"`
	FindingAidStatus            string                   `json:"finging_aid_status,omitempty"`
	FindingAidNote              string                   `json:"finding_aid_note,omitempty"`
	RevisionStatements          []*RevisionStatement     `json:"revision_statements,omitempty"`
	Instances                   []*Instance              `json:"instances,omitempty"`
	Deaccessions                []*Deaccession           `json:"deaccession,omitempty"`
	CollectionManagement        *CollectionManagement    `json:"collection_management,omitempty"`
	UserDefined                 *UserDefined             `json:"user_defined,omitempty"`
	ReleatedAccessions          []map[string]interface{} `json:"related_accessions,omitempty"`
	Classifications             []map[string]interface{} `json:"classifications,omitempty"`
	Notes                       []map[string]interface{} `json:"notes,omitempty"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	StatuteCitation        string                   `json:"statute_citation,omitempty"`
	Jurisdiction           string                   `json:"jurisdiction,omitempty"`
	TypeNote               string                   `json:"type_note,omitempty"`
	Permissions            string                   `json:"permissions,omitempty"`
	Restrictions           string                   `json:"restrictions"`
	RestrictionStartDate   *Date                    `json:"restrictions_start_date,omitempty"`
	RestrictionEndDate     *Date                    `json:"restriction_end_date,omitempty"`
	GrantedNote            string                   `json:"granted_note,omitempty"`
	ExternalDocuments      []map[string]interface{} `json:"external_documents"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	Terms                     []map[string]interface{} `json:"terms,omitempty"` // uri_or_object
	Vocabulary                string                   `json:"vocabulary,omitempty"`
	AuthorityID               string                   `json:"authority_id,omitempty"`
	ExternalDocuments         []map[string]interface{} `json:"external_documents"`

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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

	LockVersion    json.Number       `json:"lock_version,Number"`
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
	ID    int                      `json:"id,omitempty"`
	URI   string                   `json:"uri,omitempty"`
	RefID string                   `json:"ref_id,omitempty"`
	Name  string                   `json:"name,omitempty"`
	Terms []map[string]interface{} `json:"terms,omitempty"`

	LockVersion    json.Number       `json:"lock_version,Number"`
	JSONModelType  string            `json:"jsonmodel_type,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty,omitempty"`
	CreateTime     string            `json:"create_time,omitempty,omitempty"`
	Repository     map[string]string `json:"repository,omitempty"`
}

//
// String functions for cait public structures
//
func stringify(o interface{}) string {
	src, _ := json.Marshal(o)
	return string(src)
}

// String convert NoteText struct as a JSON formatted string
func (cait *NoteText) String() string {
	return stringify(cait)
}

// String convert an ArchicesSpaceAPI struct as a JSON formatted string
func (cait *ArchivesSpaceAPI) String() string {
	return stringify(cait)
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
