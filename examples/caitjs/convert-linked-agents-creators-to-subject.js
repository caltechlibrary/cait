//
// convert-linked-agents-creators-to-subject.js
//
// convert creator linked agents to subject linked agents
//
function checkAccessionAndUpdate(id, cnt) {
    var a = api.getAccession(2, id), i = 0, changed = false;
    if (a.linked_agents !== undefined && a.linked_agents.length > 0) {
        changed = false;
        for (var i = 0; i < a.linked_agents.length; i++) {
            if (a.linked_agents[i].role === "creator") {
                a.linked_agents[i].role = "subject";
                changed = true;
            }
        }
        if (changed === true) {
            api.updateAccession(a);
        }
    }
    if (cnt % 100 === 0) {
        console.log("processed", cnt, Date());
    }
}

console.log("Logging in");
api.login();
console.log("Get a list of accession ids");
ids = api.listAccessions(2);
console.log("Sort accession ids");
ids.sort();
console.log("For each accession id, check and update");
ids.forEach(checkAccessionAndUpdate);
