//
// fix-digital_object_id.js reads the JSON blobs of Digital Objects, reformats the ID and then
// updates the revised record.
//
// + api.login()
// + api.listDigitalObjects(repo_id)
// + api.getDigitalObject(repo_id, object_id)
// + api.updateDigitalObject(digital_object)
//
// The script assumes you're working with repository id no. 2
//

//
// dastardly global variables because this is a quick and dirty solution
//
var now = new Date(),
    yr = now.getUTCFullYear(),
    mn = now.getUTCMonth() + 1,
    dy = now.getUTCDate(),
    repo_id = 2,
    digital_object_ids = [],
    res,
    i = 0;

//
// helper functions
//
function lpad(s, size, chr) {
   if (chr === undefined) {
       chr = " ";
   }
   return [chr.repeat(size), s].join("").substr(-1*size);
}

function makeDigitalObjectID(sequenceNo) {
   return [yr, lpad(mn, 2, "0"), lpad(dy, 2, "0"), lpad(sequenceNo, 6, "0")].join("-");
}

function fixDigitalObjectID(s) {
    parts = s.split("/");
    i = parts[parts.length - 1];
    return makeDigitalObjectID(i);
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
digital_object_ids = api.listDigitalObjects(repo_id);
console.log("Processing", digital_object_ids.length, "digital_object_ids");
var re = new RegExp("[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9]");
for (i = 0; i < digital_object_ids.length; i++) {
    object_id = digital_object_ids[i];
    digital_object = api.getDigitalObject(repo_id, object_id);
    digital_object.digital_object_id = fixDigitalObjectID(digital_object.uri);
    res = api.updateDigitalObject(digital_object);
    if (res.status !== "Updated") {
        console.log("ERROR res", JSON.stringify(res), "retrying by appending utime");
        digital_object.digital_object_id += " " + (new Date()).getTime();
        res = api.updateDigitalObject(digital_object);
    }
    if ((i % 10 === 0) && (i > 0)) {
        console.log("Processed", i, "digital_objects");
    }
}
console.log("Processed", i, "digital_objects");
