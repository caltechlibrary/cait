//
// convert-linked-agents-creators-to-subject.js
//
// convert creator linked agents to subject linked agents
//
api.login();
ids = api.listAccessions(2);

ids.forEach(function (id) {
    var a = api.getAccession(2, id), i = 0, changed = false;
    if (a.linked_agents !== undefined && a.linked_agents.length > 0) {
        for (var i = 0; i < a.linked_agents.length; i++) {
            if (a.linked_agents[i].role === "creator") {
                a.linked_agents[i].role = "subject";
                changed = true;
            }
        }
        if (changed === true) {
            console.log("Updating", a.uri);
            api.updateAccession(a);
        }
    }
});
