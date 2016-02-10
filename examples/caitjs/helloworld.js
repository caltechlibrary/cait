//
// This is just a simple Hello World in JavaScript.
// It demonstrates getting the username we'll use to
// access the ArchivesSpace REST API.
//
(function () {
    var username = os.getEnv("CAIT_USERNAME");
    console.log("Hello World, hello " + username);
}());
