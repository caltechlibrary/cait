//
// This is a prototye of the xlsx2ead process in JavaScript.
//
var response = {},
    accession = {},
    accessionInfo = {};

// Login
response = api.login();
console.log(JSON.stringify(response));
// Get the accession info from Workbook
accessionInfo = Workbook.accessionInfo();
// Get the accession from API
accession = api.getAccession(accessionInfo.repo_id, accessionInfo.accession_id);
//
// 1. Update the Accession spreadsheet with the Accession information
// 2. Save the spreadsheet with the updated information
// 3. Create an empty EAD3 structure
// 5. Update EAD3 info from Accession Sheet
// 6. Loop through remaining sheets populating the containers
// 7. Write the EAD to disc as XML
//



