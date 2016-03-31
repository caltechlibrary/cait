//
// list-linked-agents-roles.js
//

api.login();
var ids = api.listAccessions(2),
    linkedAgentRoles = [];

ids.forEach(function (id) {
    var a = api.getAccession(2, id);

    if (a.linked_agents !== undefined && a.linked_agents.length > 0) {
        for (var i = 0; i < a.linked_agents.length; i++) {
            console.log(a.uri, a.linked_agents[i].ref, a.linked_agents[i].role);
            if (linkedAgentRoles.indexOf(a.linked_agents[i].role) < 0) {
                linkedAgentRoles.push(a.linked_agents[i].role);
            }
        }
    }
});

console.log("Roles found: ", JSON.stringify(linkedAgentRoles));
