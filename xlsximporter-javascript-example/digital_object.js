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
var // Spreadsheet description of columns c??
    cA = "Digital Object ID",
    cB = "Title",
    cC = "Series",
    cD = "Keywords",
    cE = "Name _and Subject",
    cF = "url_online oral history URL",
    cG = "Oral History Text by Item ID::Text Description",
    cH = "Oral History Text by Item ID::Search Text",
    // Convence array to do normalization with
    columnNames = [ cA, cB, cC, cD, cE, cF, cG, cH ],
    // Auth and API vars
    apiToken = "",
    apiURI = Getenv("ASPACE_API_URL"),
    apiUsername = Getenv("ASPACE_USERNAME"),
    apiPassword = Getenv("ASPACE_PASSWORD")
    // Local data locations
    dataDir = Getenv("ASPACE_DATASETS"),
    Subjects = {},
    Titles = [],
    // You could start with object IDs at 1, but this may need to be changed
    // if you have other Digital Objects already ingested.
    ObjectIDOffset = 1;

//
// Helper functions
//

// Take a Date and return it in iso8601 format per
//  https://www.w3.org/TR/NOTE-datetime, https://en.wikipedia.org/wiki/ISO_8601
function iso8601(d) {
    if (d == undefined) {
        d = new Date();
    }
    // ArchivesSpace seems to interpret iso8601 as
    // YYYY-MM-DDThh:mm:ssZ (in UTC)
    return [d.getUTCFullYear(),
        "-",
        ("0" + (d.getUTCMonth() + 1)).slice(-2),
        "-",
        ("0" + (d.getUTCDate())).slice(-2),
        "T",
        ("0" + (d.getUTCHours())).slice(-2),
        ":",
        ("0" + (d.getUTCMinutes())).slice(-2),
        ":",
        ("0" + (d.getUTCSeconds())).slice(-2),
        "Z"
    ].join("");
}

// Take a Date and return it in YYYY-MM-DD format
function yyyymmdd(d) {
    if (d === undefined) {
        d = new Date();
    }
    return [d.getUTCFullYear(),
        "-",
        ("0" + (d.getUTCMonth() + 1)).slice(-2),
        "-",
        ("0" + (d.getUTCDate())).slice(-2)
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
    return [d.getUTCFullYear(),
        " ",
        months[d.getUTCMonth()],
        " ",
        d.getUTCDate()
    ].join("");
}

//
// ArchivesSpace API methods
//

// Log into the ArchivesSpace API and save the token for re-use.
function login() {
    var data = {},
        src = "";
    src = HttpPost(apiURI + '/users/' + apiUsername + '/login', 'multipart/form-data', encodeURI('password='+apiPassword));
    data = JSON.parse(src);
    apiToken = data.session;
}

function getSubjects() {
    var topical = {},
        functional = {};
    subjectIDs = (JSON.parse(HttpGet(apiURI + "/subjects?all_ids=true", [{"X-ArchivesSpace-Session": apiToken}])));
    subjectIDs.forEach(function(id) {
        subject = JSON.parse(HttpGet(apiURI + "/subjects/" + id, [{"X-ArchivesSpace-Session": apiToken}]));
        if (subject.title !== undefined && subject.uri !== undefined) {
            subject.terms.forEach(function (term) {
                if (term.term_type === "function") {
                    console.log("DEBUG subject/function", subject.title, subject.uri);
                    functional[subject.title] = subject.uri;
                }
                if (term.term_type === "topical") {
                    console.log("DEBUG subject/topical", subject.title, subject.uri);
                    topical[subject.title] = subject.uri;
                }
            });

        }
    });
    return {Topical: topical, Functional: functional};
}


function getAccessionTitles() {
    var titles = [],
        titleIDs = [];
    titleIDs = (JSON.parse(HttpGet(apiURI + "/repositories/2/accessions?all_ids=true", [{"X-ArchivesSpace-Session": apiToken}])));
    titleIDs.forEach(function(id) {
        accession = JSON.parse(HttpGet(apiURI + "/repositories/2/accessions/" + id, [{"X-ArchivesSpace-Session": apiToken}]));
        if (accession.title !== undefined && accession.uri !== undefined) {
            console.log("DEBUG accession", accession.title, accession.uri);
            titles.push({title: accession.title, uri: accession.uri});
        }
    });
    return titles;
}

function subjectToURI(label, subjects) {
    s = label;
    if (subjects[s] !== undefined) {
        return subjects[s];
    }
    return "";
}

//
// Initialization
//

login();
Subjects = getSubjects();
Titles = getAccessionTitles();

//
// Main processing and callback
//

// callback() is the primary mapping function
function callback(row) {
    var timestamp = new Date(),
        keys = Object.keys(row),
        obj = {};

    // If we are missing a value for our digital object id, then we have an error
    if (row[cA] === undefined || row[cB] === "") {
        return {path: "", object: "", error: "Missing " + cA}
    }
    if (row[cB] === undefined || row[cB].trim() === "") {
        return {path: "", object: "", error: "Missing " + cB}
    }

    // Normalize the row fields, trim the strings
    columnNames.forEach(function (ky) {
        if (row[ky] === undefined) {
            row[ky] = "";
        } else if (typeof(row[ky]) === "string") {
            s = row[ky];
            row[ky] = s.trim();
        }
    });

    // Our spreadsheet uses row number as ID, but we have some existing
    // Digital Objects so we're going to offset our new object IDs
    objectID = parseInt(row[cA], 10) + ObjectIDOffset;
    obj = {
        digital_object_id: row[cF],
        uri: "/repositories/2/digital_objects/" + objectID,
        title: row[cB],
        publish: true,
        subjects: [],
        extents: [],
        dates: [
            {
                date_type: "single",
                label: "other",
                certainty: "",
                expression: dateExpression(timestamp),
                begin: yyyymmdd(timestamp),
                era: "",
                /*
                lock_version: 0,
                jsonmodel_type: "date",
                created_by: apiUsername,
                last_modified_by: apiUsername,
                user_mtime: iso8601(timestamp),
                system_mtime: iso8601(timestamp),
                create_time: iso8601(timestamp),
                */
            }
        ],
        notes: [],
        file_versions: [
            {
                "file_uri": row[cF],
                "publish": true,
                "jsonmodel_type": "file_version",
                /*
                "created_by": apiUsername,
                "last_modified_by": apiUsername,
                "user_mtime": iso8601(timestamp),
                "system_mtime": iso8601(timestamp),
                "create_time": iso8601(timestamp)
                */
            }
        ],
        /*
        created_by: apiUsername,
        last_modified_by: apiUsername,
        user_mtime: iso8601(timestamp),
        system_mtime: iso8601(timestamp),
        create_time: iso8601(timestamp),
        */
        repository: {
            ref: '/repositories/2'
        },
        external_documents: [],
        rights_statements: [],
        linked_instances: [],
        linked_agents: [],
        suppressed: false,
        restrictions: false,
        jsonmodel_type: "digital_object"
    };
    // Merge notes
    if (row[cH] !== "") {
        obj.notes.push({
            content: [
                row[cH]
            ],
             jsonmodel_type: "note_digital_object",
             label: row[cG],
             persistent_id: row[cF],
             publish: true,
             "type": "note"
        });
    }

    // Merge in our subjects
    subject = subjectToURI(row[cC], Subjects.Functional)
    if (subject != "") {
        obj.subjects.push({ref: subject});
    }
    subject = subjectToURI(row[cD], Subjects.Topical)
    if (subject != "") {
        obj.subjects.push({ref: subject});
    }
    // NOTE: when these Digital Objects are imported they will need to be linked to their accessions
    Titles.forEach(function (item) {
        if (item.title === obj.title) {
            obj.linked_instances.push({ref: item.uri});
        }
    });
    return {path: [dataDir, obj.uri, '.json'].join(""), object: obj, error: ""};
}
