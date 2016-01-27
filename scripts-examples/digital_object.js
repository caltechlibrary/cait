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
 * + Digital Object ID maps to URI with for /repositories/2/digital_obejcts/ID
 * + Title maps to title
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
 * + "Oral History Text by Item ID::Search Text" maps to Notes -> General Notes:Content
 * + "Oral History Text by Item ID::Text Description" maps to Notes -> General Note:Label
 *
 * All Publish fields should be true
 */
var apiToken = "",
    apiURI = "",
    apiUsername = "admin",
    dataDir = Getenv("ASPACE_DATASETS");

// Take a Date and return it in iso8601 format per
// https://www.w3.org/TR/NOTE-datetime, https://en.wikipedia.org/wiki/ISO_8601
function iso8601(d) {
    if (d == undefined) {
        d = new Date();
    }
    // // YYYY-MM-DDThh:mm:ssZ (in UTC)
    // return [d.getUTCYear(),
    //     "-",
    //     ("0" + (d.getUTCMonth() + 1)).slice(-2),
    //     "-",
    //     ("0" + (d.getUTCDate())).slice(-2),
    //     "T",
    //     ("0" + (d.getUTCHours())).slice(-2),
    //     ":",
    //     ("0" + (d.getUTCMinutes())).slice(-2),
    //     ":",
    //     ("0" + (d.getUTCSeconds())).slice(-2),
    //     "Z"
    // ].join("");
    return d.toJSON();
}

// Take a Date and return it in YYYY-MM-DD format
function yyyymmdd(d) {
    if (d === undefined) {
        d = new Date();
    }
    return [d.getYear(),
        "-",
        ("0" + (d.getMonth() + 1)).slice(-2),
        "-",
        ("0" + (d.getDate())).slice(-2)
    ].join("");
}

// Take a Date and return it in '2016 January 26' style format
function dateExpression(d) {
    var months = [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December"
    ];
    if (d === undefined) {
        d = new Date();
    }
    return [d.getFullYear(),
        " ",
        months[d.getMonth()],
        "",
        d.getDate()
    ].join("");
}

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
 //
 // The Mapping follows these rules:
 // + "Digital Object ID" maps to id value in URI
 // + "Title" maps to title
 // + "Series" maps to subject of type function
 // + "Keywords" maps to subject of type topical
 // + "Name \_and Subject" map to creator or subject based on the content of "Series"
 //      + Oral History -> Creator (interviewer/interviewee), Subject (interviewee)
 //      + Film & Video -> Subject
 //      + Institute Publications -> Creator
 //      + Manuscript Collection -> Creator
 //      + Watson Lecture -> Creator
 //      + Alumni Day -> Creator
 //      + Commencement -> Creator
 //      + Institute Publications -> Creator
 //      + "" -> Creator
 // + "url_online oral history URL" maps to unique identifier for Digital Object, File Versions-> File URI (publish should be checked), Notes -> General Notes:Persistent ID
 // + "Oral History Text by Item ID::Search Text" maps to Notes -> General Notes:Content
 // + "Oral History Text by Item ID::Text Description" maps to Notes -> General Note:Label
 // All Publish fields should be true
 //

    //console.log("DEBUG row: " + JSON.stringify(row));
    var timestamp = new Date(),
        obj = {
        uri: "/repositories/2/digital_object/"+row["Digital Object ID"],
        title: row["Title"],
        publish: true,
        subjects: [],
        extents: [],
        dates: [
            {
                date_type: "single",
                label: "migration",
                certainty: "",
                expression: dateExpression(timestamp),
                begin: yyyymmdd(timestamp),
                era: "",
                lock_version: 0,
                jsonmodel_type: "date",
                created_by: apiUsername,
                last_modified_by: apiUsername,
                user_mtime: iso8601(timestamp),
                system_mtime: iso8601(timestamp),
                create_time: iso8601(timestamp),
            }
        ],
        notes: [
            {
                content: [
                    row["Oral History Text by Item ID::Search Text"]
                ],
                 jsonmodel_type: "note_digital_object",
                 label: row["Oral History Text by Item ID::Text Description"],
                 persistent_id: row["url_online oral history URL"],
                 publish: true,
                 "type": "note"
            }
        ],
        file_versions: [
            {
                "file_uri": " http://resolver.caltech.edu/CaltechOH:OH_Clauser_F",
                "publish": true,
                "jsonmodel_type": "file_version",
                "created_by": apiUsername,
                "last_modified_by": apiUsername,
                "user_mtime": iso8601(timestamp),
                "system_mtime": iso8601(timestamp),
                "create_time": iso8601(timestamp)
            }
        ],
        created_by: apiUsername,
        last_modified_by: apiUsername,
        user_mtime: iso8601(timestamp),
        system_mtime: iso8601(timestamp),
        create_time: iso8601(timestamp),
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
    return {path: [dataDir, obj.uri, '.json'].join(""), source: obj, error: ""};
}
