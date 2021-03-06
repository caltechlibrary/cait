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
var now = new Date(),
    // We're working with the repository ID of 2, uri: /repositories/2
    yr = now.getUTCFullYear(),
    mn = now.getUTCMonth() + 1,
    dy = now.getUTCDate(),
    repoID = 2,
    sequenceNo = 0,
    response = {},
    // Spreadsheet description of columns c??
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
    apiUsername = os.getEnv("CAIT_USERNAME"),
    // Local data locations
    dataDir = os.getEnv("CAIT_DATASETS"),
    Subjects = {},
    // You could start with object IDs at 4, but this may need to be changed
    // if you have other Digital Objects already ingested.
    ObjectIDOffset = 4;

//
// Polyfills
//
if (!String.prototype.repeat) {
  String.prototype.repeat = function(count) {
    'use strict';
    if (this == null) {
      throw new TypeError('can\'t convert ' + this + ' to object');
    }
    var str = '' + this;
    count = +count;
    if (count != count) {
      count = 0;
    }
    if (count < 0) {
      throw new RangeError('repeat count must be non-negative');
    }
    if (count == Infinity) {
      throw new RangeError('repeat count must be less than infinity');
    }
    count = Math.floor(count);
    if (str.length == 0 || count == 0) {
      return '';
    }
    // Ensuring count is a 31-bit integer allows us to heavily optimize the
    // main part. But anyway, most current (August 2014) browsers can't handle
    // strings 1 << 28 chars or longer, so:
    if (str.length * count >= 1 << 28) {
      throw new RangeError('repeat count must not overflow maximum string size');
    }
    var rpt = '';
    for (;;) {
      if ((count & 1) == 1) {
        rpt += str;
      }
      count >>>= 1;
      if (count == 0) {
        break;
      }
      str += str;
    }
    // Could we try:
    // return Array(count + 1).join(this);
    return rpt;
  }
}

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
function getSubjects() {
    var topical = {},
        functional = {};
    subjectIDs = api.listSubjects();
    subjectIDs.forEach(function(id) {
        subject = api.getSubject(id);
        if (subject.title !== undefined && subject.uri !== undefined) {
            subject.terms.forEach(function (term) {
                if (term.term_type === "function") {
                    console.log("Collecting subject/function", subject.title, subject.uri);
                    functional[subject.title] = subject.uri;
                }
                if (term.term_type === "topical") {
                    console.log("Collecting subject/topical", subject.title, subject.uri);
                    topical[subject.title] = subject.uri;
                }
            });

        }
    });
    return {Topical: topical, Functional: functional};
}

function subjectToURI(label, subjects) {
    s = label;
    if (subjects[s] !== undefined) {
        return subjects[s];
    }
    return "";
}

function lpad(s, size, chr) {
    if (chr === undefined) {
        chr = " ";
    }
    return [chr.repeat(size), s].join("").substr(-1*size);
}

function makeDigitalObjectID(onlineURL) {
    sequenceNo++;
    return [yr, lpad(mn, 2, "0"), lpad(dy, 2, "0"), lpad(sequenceNo, 6, "0")].join("-");
}

//
// Initialization
//
sequenceNo = 0;
// Make sure the environment varaibles are all set.
["CAIT_API_URL", "CAIT_USERNAME", "CAIT_PASSWORD", "CAIT_DATASETS"].forEach(function (envvar) {
    var s = os.getEnv(envvar);
    if (s == "") {
        console.log("Missing", envvar);
        os.exit(1);
    }
});
console.log("Environment defined, authenticating");
response = api.login();
if (response.error !== undefined) {
    console.log(response.error);
    os.exit(1);
}
console.log("Authenticated");
console.log("Saving results to", dataDir);

Subjects = getSubjects();

//
// Main callback function for processing row data
//
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
        digital_object_id: makeDigitalObjectID(row[cF]),
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
            }
        ],
        notes: [],
        file_versions: [
            {
                "file_uri": row[cF],
                "publish": true,
                "jsonmodel_type": "file_version",
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
    s = [dataDir, obj.uri, '.json'].join("");
    return {path: s, object: obj, error: ""};
}
