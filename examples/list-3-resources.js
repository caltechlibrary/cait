/* Generate a JSON list of three resources */
var max_resources = 3,
    resource_list = [];

function buildResourceList(val) {
    if (resource_list.length < max_resources) {
        resource = api.getResource(2, val);
        if (resource) {
            resource_list.push(resource);
        }
    }
}

/* Main processing */
api.login();
resources = api.listResources(2)
resources.forEach(buildResourceList);
//console.log("DEBUG length resource_list: "+resource_list.length);
console.log(JSON.stringify(resource_list));
