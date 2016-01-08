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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// The library version
var Version = "0.0.0"

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

func checkEnv(apiURL, apiToken, username, password string) bool {
	if strings.TrimSpace(apiURL) == "" {
		return false
	}
	if strings.TrimSpace(apiToken) == "" {
		if strings.TrimSpace(username) == "" {
			return false
		}
		if strings.TrimSpace(password) == "" {
			log.Println("WARNING: using an empty string for password")
			//return false
		}
	}
	return true
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(apiURL, apiToken, username, password string) *ArchivesSpaceAPI {
	aspace := new(ArchivesSpaceAPI)
	if checkEnv(apiURL, apiToken, username, password) == false {
		log.Fatal("Cannot create a new ArchivesSpace API connection")
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("ArchivesSpace connection not configured, %s", err)
	}
	aspace.URL = u
	aspace.Username = username
	aspace.Password = password
	aspace.AuthToken = apiToken
	return aspace
}

// IsAuth returns true if the auth token has been set, false otherwise
func (aspace *ArchivesSpaceAPI) IsAuth() bool {
	if aspace.AuthToken == "" {
		return false
	}
	return true
}

// Login authenticates against the ArchivesSpace REST API setting the AuthToken
// value in the ArchivesSpaceAPI struct.
func (aspace *ArchivesSpaceAPI) Login() error {
	// See https://golang.org/pkg/net/url/#pkg-examples for example building a URL from parts.
	// Command line example: curl -F "password=admin" "http://localhost:8089/users/admin/login"
	var data map[string]interface{}

	// If we already have a token set then logout and get a new one
	if aspace.IsAuth() == true {
		aspace.Logout()
	}

	u := aspace.URL
	u.Path = fmt.Sprintf("/users/%s/login", aspace.Username)
	form := url.Values{}
	form.Add("password", aspace.Password)

	res, err := http.PostForm(u.String(), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Status != "200 OK" {
		return fmt.Errorf("ArchivesSpace returned HTTP status %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("ArchivesSpace return unreadable body: %s", err)
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("Can't process JSON response %s\n\t%s", content, err)
	}
	aspace.AuthToken = data["session"].(string)
	return nil
}

// Logout clear the authentication token for the session with the API
func (aspace *ArchivesSpaceAPI) Logout() error {
	// Save the token and invalidate the one in our aspace struct.
	token := aspace.AuthToken
	aspace.AuthToken = ""
	// Using the copied token try to logout from the service.
	u := aspace.URL
	u.Path = `/logout`
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-ArchivesSpace-Session", token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// API the common HTTP request processing for interacting with ArchivesSpaceAPI
func (aspace *ArchivesSpaceAPI) API(method string, url string, data interface{}) ([]byte, error) {
	var (
		payload []byte
		err     error
	)
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("Can't create request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	if method == "POST" {
		res, err := client.Do(req)
		defer res.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("Request error: %s", err)
		}
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Read body error: %s", err)
		}
		return content, nil
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	if res.Status != "200 OK" {
		return nil, fmt.Errorf("ArchiveSpace API error %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	return content, nil
}

// CreateAPI is a generalized call to create an object form an interface.
func (aspace *ArchivesSpaceAPI) CreateAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := aspace.API("POST", url, obj)
	if err != nil {
		return nil, err
	}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAPI is a generalized call to get a specific object from an interface
// obj is modified as a side effect
func (aspace *ArchivesSpaceAPI) GetAPI(url string, obj interface{}) error {
	content, err := aspace.API("GET", url, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, obj)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAPI is a generalized call to update an object from an interface.
func (aspace *ArchivesSpaceAPI) UpdateAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := aspace.API("POST", url, obj)
	if err != nil {
		return nil, err
	}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Could not unpack UpdateAPI() response [%s] %s", content, err)
	}
	return data, nil
}

// DeleteAPI is a generalized call to update an object form an interface
func (aspace *ArchivesSpaceAPI) DeleteAPI(url string, obj interface{}) (*ResponseMsg, error) {
	content, err := aspace.API("DELETE", url, obj)
	if err != nil {
		return nil, err
	}

	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Cannnot decode DeleteAPI() response %s", err)
	}
	return data, nil
}

// ListAPI return a list of IDs from an ArchivesSpace instance for given URL
func (aspace *ArchivesSpaceAPI) ListAPI(url string) ([]int, error) {
	content, err := aspace.API("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// [1,2,3,4]
	var ids []int
	err = json.Unmarshal(content, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// CreateRepository will create a respository via the REST API for the
// ArchivesSpace instance defined in the ArchivesSpaceAPI struct.
// It will return the created record.
func (aspace *ArchivesSpaceAPI) CreateRepository(repo *Repository) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = "/repositories"
	return aspace.CreateAPI(u.String(), repo)
}

// GetRepository returns the repository details based on Id
func (aspace *ArchivesSpaceAPI) GetRepository(id int) (*Repository, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf(`/repositories/%d`, id)
	repo := new(Repository)
	err := aspace.GetAPI(u.String(), repo)
	if err != nil {
		return nil, err
	}
	repo.ID = id
	return repo, nil
}

// UpdateRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) UpdateRepository(repo *Repository) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = repo.URI
	return aspace.UpdateAPI(u.String(), repo)
}

// DeleteRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) DeleteRepository(repo *Repository) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/repositories/%d", repo.ID)
	return aspace.DeleteAPI(u.String(), repo)
}

// ListRepositories returns a list of repositories available via the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) ListRepositories() ([]Repository, error) {
	u := *aspace.URL
	u.Path = `/repositories`

	content, err := aspace.API("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"lock_version":0,"repo_code":"1447893780","name":"This is a test generated from go_test","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-19T00:43:00Z","system_mtime":"2015-11-19T00:43:00Z","user_mtime":"2015-11-19T00:43:00Z","jsonmodel_type":"repository","uri":"/repositories/16","agent_representation":{"ref":"/agents/corporate_entities/15"}}
	var repos []Repository
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, err
	}
	// Now I need to populate the repos[?].ID fields
	for i := range repos {
		if id, err := strconv.Atoi(strings.TrimPrefix(repos[i].URI, "/repositories/")); err == nil {
			repos[i].ID = id
		}
	}
	return repos, nil
}

// CreateAgent creates a Agent recod via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) CreateAgent(aType string, agent *Agent) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/agents/%s", aType)
	return aspace.CreateAPI(u.String(), agent)
}

// GetAgent return an Agent via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) GetAgent(agentType string, agentID int) (*Agent, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf(`/agents/%s/%d`, agentType, agentID)

	agent := new(Agent)
	err := aspace.GetAPI(u.String(), agent)
	if err != nil {
		return nil, err
	}

	// Make sure the ID comes from agent.URI
	p := strings.Split(agent.URI, "/")
	id, err := strconv.Atoi(p[len(p)-1])
	if err != nil {
		return agent, err
	}
	agent.ID = id
	return agent, nil
}

// UpdateAgent creates a Agent recod via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) UpdateAgent(agent *Agent) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = agent.URI
	return aspace.UpdateAPI(u.String(), agent)
}

// DeleteAgent creates a Agent record via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) DeleteAgent(agent *Agent) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = agent.URI
	return aspace.DeleteAPI(u.String(), agent)
}

// ListAgents return an array of Agents via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) ListAgents(agentType string) ([]int, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf(`/agents/%s`, agentType)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

// CreateAccession creates a new Accession record in a Repository
func (aspace *ArchivesSpaceAPI) CreateAccession(repoID int, accession *Accession) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/repositories/%d/accessions", repoID)
	return aspace.CreateAPI(u.String(), accession)
}

// GetAccession retrieves an Accession record from a Repository
func (aspace *ArchivesSpaceAPI) GetAccession(repoID, accessionID int) (*Accession, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/repositories/%d/accessions/%d", repoID, accessionID)

	accession := new(Accession)
	err := aspace.GetAPI(u.String(), accession)
	if err != nil {
		return nil, err
	}
	p := strings.Split(accession.URI, "/")
	accession.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return accession, fmt.Errorf("Accession ID parse error %d %s", accession.ID, err)
	}
	return accession, nil
}

// UpdateAccession updates an existing Accession record in a Repository
func (aspace *ArchivesSpaceAPI) UpdateAccession(accession *Accession) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = accession.URI
	return aspace.UpdateAPI(u.String(), accession)
}

// DeleteAccession deleted an Accession record from a Repository
func (aspace *ArchivesSpaceAPI) DeleteAccession(accession *Accession) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = accession.URI
	return aspace.DeleteAPI(u.String(), accession)
}

// ListAccessions return a list of Accession IDs from a Repository
func (aspace *ArchivesSpaceAPI) ListAccessions(repositoryID int) ([]int, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf(`/repositories/%d/accessions`, repositoryID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

// CreateSubject creates a new Subject in ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) CreateSubject(subject *Subject) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = "/subjects"
	return aspace.CreateAPI(u.String(), subject)
}

// GetSubject retrieves a subject record from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) GetSubject(subjectID int) (*Subject, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/subjects/%d", subjectID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"subject","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","vocabulary":"/vocabularies/1"}],"external_documents":[],"uri":"/subjects/1","is_linked_to_published_record":true,"vocabulary":"/vocabularies/1"}
	subject := new(Subject)
	err := aspace.GetAPI(u.String(), subject)
	p := strings.Split(subject.URI, "/")
	subject.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return subject, fmt.Errorf("Accession ID parse error %d %s", subject.ID, err)
	}
	return subject, nil
}

// UpdateSubject updates an existing subject record in an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) UpdateSubject(subject *Subject) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = subject.URI
	return aspace.UpdateAPI(u.String(), subject)
}

// DeleteSubject deletes a subject from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) DeleteSubject(subject *Subject) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = subject.URI
	return aspace.DeleteAPI(u.String(), subject)
}

// ListSubjects return a list of Subject IDs from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) ListSubjects() ([]int, error) {
	u := *aspace.URL
	u.Path = `/subjects`
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

// CreateVocabulary creates a new Vocabulary in ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) CreateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = "/vocabularies"
	return aspace.CreateAPI(u.String(), vocabulary)
}

// GetVocabulary retrieves a vocabulary record from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) GetVocabulary(vocabularyID int) (*Vocabulary, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/vocabularies/%d", vocabularyID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"vocabulary","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","vocabulary":"/vocabularies/1"}],"external_documents":[],"uri":"/vocabularys/1","is_linked_to_published_record":true,"vocabulary":"/vocabularies/1"}
	vocabulary := new(Vocabulary)
	err := aspace.GetAPI(u.String(), vocabulary)
	p := strings.Split(vocabulary.URI, "/")
	vocabulary.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return vocabulary, fmt.Errorf("Accession ID parse error %d %s", vocabulary.ID, err)
	}
	return vocabulary, nil
}

// UpdateVocabulary updates an existing vocabulary record in an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) UpdateVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = vocabulary.URI
	return aspace.UpdateAPI(u.String(), vocabulary)
}

// DeleteVocabulary deletes a vocabulary from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) DeleteVocabulary(vocabulary *Vocabulary) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = vocabulary.URI
	return aspace.DeleteAPI(u.String(), vocabulary)
}

// ListVocabularies return a list of Vocabulary IDs from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) ListVocabularies() ([]int, error) {
	u := *aspace.URL
	u.Path = `/vocabularies`
	/*
		q := u.Query()
		q.Set("all_ids", "true")
		u.RawQuery = q.Encode()
	*/
	content, err := aspace.API("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	var (
		ids          []int
		vocabularies []Vocabulary
	)
	err = json.Unmarshal([]byte(content), &vocabularies)
	for _, val := range vocabularies {
		p := strings.Split(val.URI, "/")
		id, err := strconv.Atoi(p[len(p)-1])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateTerm creates a new Term in ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) CreateTerm(vocabularyID int, term *Term) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/vocabularies/%d/terms", vocabularyID)
	return aspace.CreateAPI(u.String(), term)
}

// GetTerm retrieves a term record from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) GetTerm(vocabularyID, termID int) (*Term, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/vocabularies/%d/terms/%d", vocabularyID, termID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"term","external_ids":[],"publish":true,"terms":[{"lock_version":0,"term":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","term_type":"function","jsonmodel_type":"term","uri":"/terms/1","term":"/terms/1"}],"external_documents":[],"uri":"/terms/1","is_linked_to_published_record":true,"term":"/terms/1"}
	term := new(Term)
	err := aspace.GetAPI(u.String(), term)
	if err != nil {
		return nil, err
	}
	p := strings.Split(term.URI, "/")
	term.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return term, fmt.Errorf("Accession ID parse error %d %s", term.ID, err)
	}
	return term, nil
}

// UpdateTerm updates an existing term record in an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) UpdateTerm(term *Term) (*ResponseMsg, error) {
	u := aspace.URL
	u.Path = term.URI
	return aspace.UpdateAPI(u.String(), term)
}

// DeleteTerm deletes a term from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) DeleteTerm(term *Term) (*ResponseMsg, error) {
	u := aspace.URL
	u.Path = term.URI
	return aspace.DeleteAPI(u.String(), term)
}

// ListTerms return a list of Term IDs from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) ListTerms(vocabularyID int) ([]int, error) {
	u := aspace.URL
	u.Path = fmt.Sprintf(`/vocabularies/%d/terms`, vocabularyID)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

// CreateLocation creates a new Location in ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) CreateLocation(location *Location) (*ResponseMsg, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/locations")
	return aspace.CreateAPI(u.String(), location)
}

// GetLocation retrieves a location record from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) GetLocation(ID int) (*Location, error) {
	u := *aspace.URL
	u.Path = fmt.Sprintf("/locations/%d", ID)

	// content should look something like
	// {"lock_version":121,"title":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T23:16:19Z","user_mtime":"2015-10-19T22:45:07Z","source":"local","jsonmodel_type":"location","external_ids":[],"publish":true,"locations":[{"lock_version":0,"location":"Commencement","created_by":"admin","last_modified_by":"admin","create_time":"2015-10-19T22:45:07Z","system_mtime":"2015-10-19T22:45:07Z","user_mtime":"2015-10-19T22:45:07Z","location_type":"function","jsonmodel_type":"location","uri":"/locations/1","location":"/locations/1"}],"external_documents":[],"uri":"/locations/1","is_linked_to_published_record":true,"location":"/locations/1"}
	location := new(Location)
	err := aspace.GetAPI(u.String(), location)
	if err != nil {
		return nil, err
	}
	p := strings.Split(location.URI, "/")
	location.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return location, fmt.Errorf("Accession ID parse error %d %s", location.ID, err)
	}
	return location, nil
}

// UpdateLocation updates an existing location record in an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) UpdateLocation(location *Location) (*ResponseMsg, error) {
	u := aspace.URL
	u.Path = location.URI
	return aspace.UpdateAPI(u.String(), location)
}

// DeleteLocation deletes a location from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) DeleteLocation(location *Location) (*ResponseMsg, error) {
	u := aspace.URL
	u.Path = location.URI
	return aspace.DeleteAPI(u.String(), location)
}

// ListLocations return a list of Location IDs from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) ListLocations() ([]int, error) {
	u := aspace.URL
	u.Path = fmt.Sprintf(`/locations`)
	q := u.Query()
	q.Set("all_ids", "true")
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

// Search return a list of search results from an ArchivesSpace instance
func (aspace *ArchivesSpaceAPI) Search(opt *map[string]interface{}) ([]int, error) {
	u := aspace.URL
	u.Path = fmt.Sprintf(`/search`)
	q := u.Query()
	for k, v := range *opt {
		q.Set(k, fmt.Sprintf("%s", v))
	}
	if q.Get("page") == "" {
		q.Set("page", "1")
	}
	u.RawQuery = q.Encode()
	return aspace.ListAPI(u.String())
}

//FIXME: need Create, Get, Update, Delete, List functions for DigitalObject, Instances, Extents, Resource, Group, Users
//FIXME: Need Get/query methods for /terms, /search/*

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
