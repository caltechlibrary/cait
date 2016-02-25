//
// mk-agent-name-direct.js - if the name is stored indirect, correct the ordering and switch to direct
//
var args = os.args();

if (args.length === 0) {
    console.log("USAGE: caitjs mk-agent-name-direct.js");
    os.exit(1);
}
api.login();

args.forEach(function (arg, i) {
    var lock_version = 0,
        primary_name = "",
        rest_of_name = "",
        sort_name = "",
        sort = [],
        parts = [];

    id = parseInt(arg);
    if (! id) {
        console.log("Error converting arg", i, "of", args.join(", "));
        os.exit(1);
    }
    agent = api.getAgent("people", id);
    if (agent.display_name.primary_name.indexOf(",") > -1) {
        parts = agent.display_name.primary_name.split(",");
        primary_name = parts[0];
        rest_of_name = agent.display_name.rest_of_name.trim();
        if (rest_of_name.length === 0){
            rest_of_name = parts[1].trim();
        } else {
            rest_of_name = [rest_of_name, parts[1]].join(", ");
        }
        /*
        if (agent.display_name.qualifier.length > 0) {
            sort_name = [primary_name, ", ", rest_of_name, " (", agent.display_name.qualifier, ")"].join("");
        } else {
            sort_name = [primary_name, ", ", rest_of_name];
        }
        */
        //agent.display_name.sort_name = sort_name;
        agent.display_name.primary_name = primary_name;
        agent.display_name.rest_of_name = rest_of_name;
        agent.display_name.name_order = "direct";
    }
    for (var i = 0; i < agent.names.length; i++) {
        if (agent.names[i].primary_name.indexOf(",") > -1) {
            parts = agent.names[i].primary_name.split(",", 2);
            primary_name = parts[0];
            rest_of_name = agent.names[i].rest_of_name.trim();
            if (rest_of_name.length === 0){
                rest_of_name = parts[1].trim();
            } else {
                rest_of_name = [rest_of_name, parts[1]].join(", ");
            }
            /*
            if (agent.names[i].qualifier.length > 0) {
                sort_name = [primary_name, ", ", rest_of_name, " (", agent.names[i].qualifier, ")"].join("");
            } else {
                sort_name = [primary_name, ", ", rest_of_name];
            }
            */
            //agent.names[i].sort_name = sort_name;
            agent.names[i].primary_name = primary_name;
            agent.names[i].rest_of_name = rest_of_name;
            agent.names[i].name_order = "direct";
        }
    }
    /*
    if (agent.title != agent.display_name.sort_name) {
        agent.title = agent.display_name.sort_name
    }*/
    console.log(JSON.stringify(agent.display_name));
    resp = api.updateAgent(agent);
    console.log(JSON.stringify(resp));
});
