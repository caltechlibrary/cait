/* Generate a JSON list of three resources */
var resource_list = [];
function buildResourceList(val) {
    resource = api.getResource(2, val);
    if (resource) {
        resource_list.push(resource);
    }
}

/* Main processing */
api.login();
resources = api.listResources(2)
resources.forEach(buildResourceList);
console.log(JSON.stringify(resource_list, null, "  "));
