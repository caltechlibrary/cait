//
// In this example we sign on to the ArchivesSpace API and then signoff
//
var results;
results = api.login();
console.log("Login executed " + results);
results = api.logout();
console.log("Logout executed " + results);
