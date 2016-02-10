//
// In this example we
// + Look at the command line parameters for the repository ID
// + Login to the API
// + Get the repository information
//
var args = os.args(),
    resp = api.login();

if (resp.isAuth === undefined || resp.isAuth === false) {
    console.log("Could not log in", resp.error);
    os.exit(1);
}

if (args.length != 1) {
    console.log("Need to provide the repostory number in the command line args");
    os.exit(1);
}
var repo = api.getRepository(args[0]);
console.log("Repository Info", args[0], ":", JSON.stringify(repo));

