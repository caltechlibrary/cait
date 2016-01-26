/**
 * This is an example JavaScript file for importing Digital Objects from a
 * Excel spreadsheet with the following columns.
 * + Digital Object ID
 * + Title
 * + Series
 * + Keywords
 * + Name _and Subject
 * + url_online oral history URL
 * + Oral History Text by Item ID::Text Description
 * + Oral History Text by Item ID::Search Text
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
 *    datesarray JSONModel(:date) object
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
    apiURI = "";

// Log into the ArchivesSpace API and save the token for re-use.
function login() {
    var username = Getenv("ASPACE_USERNAME"),
        password = Getenv("ASPACE_PASSWORD"),
        data = "",
        uri = "";
    apiURI = Getenv("ASPACE_API_URL"),
    src = HttpPost(apiURI + '/users/' + username + '/login', 'multipart/form-data', encodeURI('password='+password));
    data = JSON.parse(src);
    apiToken = data.session;
}

// Login to the API.
login();
//console.log('\texport ASPACE_API_TOKEN: ' + apiToken);
//console.log('\t curl -H "X-ArchivesSpace-Session: ' + apiToken + '" ' + apiURI);

function callback(row) {
    //FIXME: look up accession that is related
    //FIXME: Need to figure out exactly how the row object maps to ArchivesSpace's
    // concept of a Digital Object.

    //console.log("DEBUG row: " + JSON.stringify(row));
    var obj = {
        uri: "/repositories/2/digital_object",
        title: row["Title"],
        publish: true,
        subject: [],
        extents: [],
        dates: [],
        external_documents: [],
        rights_statements: [],
        linked_agents: [],
        suppressed: false,
        restrictions: false
    };
    //console.log("DEBUG obj: " + JSON.stringify(obj));
    return {path: "", source: obj, error: ""};
}
