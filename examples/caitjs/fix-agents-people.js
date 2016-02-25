//
// fix-agent-peoples.js - this is a clean up script post Java MigrationTool import.
//
// In some of our agents/people records the name did not migrate
// properly with the Java based Migration tool.  This script looks at each agent
// record, if the primary_name is populated with a comma and rest_of_name is empty
// this program will split the name with everything after the comma put in rest_of_ name field.
//
// Set the name_order to inverted for display_name and names.
//
// All our agents/people should be "published"
//

// splitName takes a nameObject split's the primary_name at the comma into primary_name, rest_of_name properties
function splitName(nameObject) {
    "use strict";
    var parts = [];

    if (nameObject.primary_name === undefined || nameObject.rest_of_name === undefined) {
        console.log("error in nameObject", JSON.stringify(nameObject));
        os.exit(1);
    }

    // Don't make any changes as the rest_of_name field appears populated
    if (nameObject.rest_of_name.trim().length !== 0){
        return nameObject;
    }
    parts = nameObject.primary_name.split(",", 2);
    nameObject.primary_name = parts[0].trim();
    nameObject.rest_of_name = parts[1].trim();
    return nameObject;
}

// fixAgentPerson gets the agent record from id and fixes name, publish and name_order
function fixAgentPerson(id) {
    "use strict";
    var agent = {};

    agent = api.getAgent("people", id);
    // console.log("DEBUG before agent\n", JSON.stringify(agent, null, "\t"));
    agent.publish = true;// Set publish flag to true
    agent.name_order = "inverted";// Set name order to inverted (primary_name, rest_of_name)
    if (agent.display_name.primary_name.indexOf(",") > -1) {
        agent.display_name = splitName(agent.display_name);
    }
    for (var i = 0; i < agent.names.length; i++) {
        if (agent.names[i].primary_name.indexOf(",") > -1) {
            agent.names[i] = splitName(agent.names[i]);
        }
        agent.names[i].name_order = "inverted";// Force name order to direct
    }
    // console.log("DEBUG agent after\n", JSON.stringify(agent, null, "\t"));
    // console.log(JSON.stringify(agent));
    return api.updateAgent(agent);
}

//
// Main process to fix the agent record
//
(function () {
    "use strict";
    var args = [],
        arg = "",
        agent_id = 0,
        resp = {};

    args = os.args();
    if (args.length === 0) {
        console.log("USAGE: caitjs fix-agent-peoples.js AGENT_ID");
        os.exit(1);
    }
    api.login();

    args.forEach(function (arg) {
        var agent_id = 0;
        agent_id = parseInt(arg);
        if (! agent_id) {
            console.log("Error converting args into agent_id", typeof agent_id, "for", args.join(", "));
            os.exit(1);
        }
        resp = fixAgentPerson(agent_id);
        console.log(JSON.stringify(resp));
    });
}());
