/**
 * This is an example JavaScript file for importing Digital Objects from a
 * Excel spreadsheet with the following columns.
 * + Digital Object ID
 * + Title (map to title)
 * + Series
 * + Keywords
 * + Name _and Subject
 * + url_online oral history URL
 * + Oral History Text by Item ID::Text Description
 * + Oral History Text by Item ID::Search Text
 *
 * The Mapping follows these rules:
 * + "Series" maps to subject of type function
 * + "Keywords" maps to subject of type topical
 * + "Name \_and Subject" map to creator or subject based on the content of "Series"
 *      + Oral History -> Creator (interviewer/interviewee), Subject (interviewee)
 *      + Film & Video -> Subject
 *      + Institute Publications -> Creator
 *      + Manuscript Collection -> Creator
 *      + Watson Lecture -> Creator
 *      + Alumni Day -> Creator
 *      + Commencement -> Creator
 *      + Institute Publications -> Creator
 *      + "" -> Creator
 * + "url_online oral history URL" maps to unique identifier for Digital Object, File Versions-> File URI (publish should be checked), Notes -> General Notes:Persistent ID
 * + "Oral History Text by Item ID::Text Description" maps to Notes -> General Notes:Content
 * + "Oral History Text by Item ID::Text Description" maps to Notes -> General Note:Label
 * All Publish fields should be true
 *
 * These need to map into a JSONModel(:digital_object) in ArchivesSpace.
 *
 * ```
 *    uri string
 *    external_ids  array JSONModel(:external_id) object
 *    title string [1..16384]
 *    language string
 *    publish boolean
 *    subjects array object
 *    linked_events array object
 *    extents array JSONModel(:extent) object
 *    dates array JSONModel(:date) object
 *    external_documents array JSONModel(:external_document) object
 *    rights_statements array JSONModel(:rights_statement) object
 *    linked_agents array object
 *    suppressed boolean
 *    lock_version integer,string
 *    jsonmodel_type string
 *    created_by string
 *    last_modified_by string
 *    user_mtime date-time
 *    system_mtime date-time
 *    create_time date-time
 *    repository object
 *    digital_object_id string [..255]
 *    level string
 *    digital_object_type string
 *    file_versions array JSONModel(:file_version) object
 *    restrictions boolean default false
 *    tree object
 *    notesarray [object Object],[object Object]
 *    collection_management JSONModel(:collection_management) object
 *    user_defined JSONModel(:user_defined) object
 *    linked_instances array
 * ```
 */
var apiToken = "",
    apiURI = "",
    apiUsername = "admin";

// Log into the ArchivesSpace API and save the token for re-use.
function login() {
    var password = Getenv("ASPACE_PASSWORD"),
        data = {},
        src = "";
    apiUsername = Getenv("ASPACE_USERNAME"),
    apiURI = Getenv("ASPACE_API_URL"),
    src = HttpPost(apiURI + '/users/' + apiUsername + '/login', 'multipart/form-data', encodeURI('password='+password));
    data = JSON.parse(src);
    apiToken = data.session;
}

// Login to the API.
login();
//console.log('\texport ASPACE_API_TOKEN: ' + apiToken);
//console.log('\t curl -H "X-ArchivesSpace-Session: ' + apiToken + '" ' + apiURI);

function callback(row) {
    //FIXME: need the current date/time in various formats...
    //FIXME: look up accession that is related so I can populate linked_instances
    //FIXME: Need to figure out exactly how the row object maps to ArchivesSpace's
    // concept of a Digital Object.

    //console.log("DEBUG row: " + JSON.stringify(row));
    var obj = {
        uri: "/repositories/2/digital_object",
        title: row["Title"],
        publish: true,
        subjects: [],
        extents: [],
        dates: [
            {
                date_type: "single",
                label: "migration",
                certainty: "",
                expression: "2016 January 26",
                begin: "2016-01-26",
                era: "",
                lock_version: 0,
                jsonmodel_type: "date",
                created_by: apiUsername,
                last_modified_by: apiUsername,
                user_mtime: "2016-01-26T00:00:00Z",
                system_mtime: "2016-01-26T00:00:00Z",
                create_time: "2016-01-26T00:00:00Z"
            }
        ],
        created_by: apiUsername,
        last_modified_by: apiUsername,
        user_mtime: "2016-01-26T00:00:00Z",
        system_mtime: "2016-01-26T00:00:00Z",
        create_time: "2016-01-26T00:00:00Z",
        repository: {
            ref: '/repositories/2'
        },
        external_documents: [],
        rights_statements: [],
        linked_agents: [],
        suppressed: false,
        restrictions: false,
        jsonmodel_type: "digital_object"
    };
    //console.log("DEBUG obj: " + JSON.stringify(obj));
    return {path: "", source: obj, error: ""};
}
