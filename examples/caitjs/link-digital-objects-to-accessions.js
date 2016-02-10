//
// Link digital object (by title) to their accession (which has matching title)
// This script demonstrates.
//
// + api.login()
// + api.listAccessions(repo_id)
// + api.listDigitalObjects(repo_id)
// + api.updateAccession(accession)
//
// The script assumes you're working with repository id no. 2
//

//
// dastardly global variables because this is a quick and dirty solution
//
var repo_id = 2,
    accession_ids = [],
    accession_titles = {},
    digital_object_ids = [],
    res,
    i = 0;

//
// helpful functions
//
function mapAccessionTitle(accession_id) {
    var accession;
    accession = api.getAccession(repo_id, accession_id);
    if (accession_titles[accession.title] === undefined) {
        accession_titles[accession.title] = [];
    }
    if (accession_titles[accession.title] !== undefined){
        accession_titles[accession.title].push(accession_id);
    }
}

function makeDigitalObjectInstance(digital_object) {
    /*example instances array
  "instances": [
    {
      "create_time": "2016-02-10T01:04:54Z",
      "created_by": "admin",
      "digital_object": {
        "ref": "/repositories/2/digital_objects/278"
      },
      "instance_type": "digital_object",
      "jsonmodel_type": "instance",
      "last_modified_by": "admin",
      "lock_version": 0,
      "system_mtime": "2016-02-10T01:04:54Z",
      "user_mtime": "2016-02-10T01:04:54Z"
    }
  ],
    */

    // Minimal fields to add to our accession instances
    return {
        instance_type: "digital_object",
        digital_object: {
            ref: digital_object.uri
        },
        jsonmodel_type: "instance"
    }
}
function linkAccession(digital_object, accession_ids) {
    console.log(digital_object.title, " --> ", JSON.stringify(accession_ids));
    accession_ids.forEach(function (id) {
        accession = api.getAccession(repo_id, id);
        if (accession.instances === undefined) {
            accession.instances = [];
        }
        accession.instances.push(makeDigitalObjectInstance(digital_object));
        console.log(JSON.stringify(api.updateAccession(accession)));
    });
}

//
// Main processing starts here...
//
console.log("Authenticating");
res = api.login();
if (res.error !== undefined) {
    console.log("Login error", res.error);
    os.exit();
}

console.log("Getting accession id list");
accession_ids = api.listAccessions(2);
console.log("Scanning", accession_ids.length, "accessions for titles");
for (i = 0; i < accession_ids.length; i++) {
    id = accession_ids[i];
    mapAccessionTitle(id);
    if ((i % 100 === 0) && (i > 0)) {
        console.log("Scanned", i, "accessions for titles");
    }
}
console.log("Scanned", i, "accessions for titles");

console.log("Getting a list of digital object ids");
digital_object_ids = api.listDigitalObjects(2);
console.log("Processing", digital_object_ids.length, "digital objects");
for (i = 0; i < digital_object_ids.length; i++) {
    digital_object_id = digital_object_ids[i];
    digital_object = api.getDigitalObject(repo_id, digital_object_id)
    // NOTE: look up title to find accession_id(s)
    if (digital_object.title !== undefined && accession_titles[digital_object.title] !== undefined) {
        linkAccession(digital_object, accession_titles[digital_object.title]);
    }

    if ((i % 10 === 0) && (i > 0)) {
        console.log("Processed", i, "digital objects");
    }
}
