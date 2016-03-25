//
// accessions-agent-roles.xlsx.js will crawls the accessions data
// and build an Excel xlsx workbook of accession title, accession ids, agent name(s), agent role(s)
//
(function () {
    "use strict";

    var workbookName = "accessions-agents-roles.xlsx",
        rowCount = 0,
        totalAccessionIDs = 0,
        res = {},
        workbook = {},
        accessionIDs = [];

    // Setup empty workbook
    workbook = {
        "Accession Agent Roles": [
            [
                "Accession URI",
                "Accession Title",
                "Creator Agents URI",
                "Creator Agents Names",
                "Subject Agents URI",
                "Subject Agents Names",
                "Completed"
            ]
        ]
    };


    function appendRow(wk, uri, title, ca_uri, ca_names, sa_uri, sa_names) {
            var table = wk["Accession Agent Roles"];
            table.push([uri, title, ca_uri, ca_names, sa_uri, sa_names, ""]);
    }

    function show(obj) {
        console.log(JSON.stringify(obj, null, "  "));
    }

    function addAccession(id, i) {
        var accession = {},
            uri = "",
            title = "",
            ca_uri = [],
            ca_names = [],
            sa_uri = [],
            sa_names = [];

        accession = api.getAccession(2, id);
        uri = ["", "repostiories", 2, "accessions", id].join("/");
        title = accession.title;
        //show(accession);// DEBUG
        if (accession.linked_agents !== undefined) {
            // linkedAgents = accession.linked_agents;
            accession.linked_agents.forEach(function (link) {
                if (link.role === "creator") {
                    ca_uri.push(link.ref);
                } else {
                    sa_uri.push(link.ref);
                }
            })

        }
        ca_uri.forEach(function (uri) {
            var parts = uri.split("/"),
                id = 0;
            id = parts[(parts.length - 1)];
            agent = api.getAgent("people", id);
            if (agent.display_name !== undefined && agent.display_name.sort_name !== undefined) {
                ca_names.push(agent.display_name.sort_name);
            }
        })
        sa_uri.forEach(function (uri) {
            var parts = uri.split("/"),
                id = 0;
            id = parts[(parts.length - 1)];
            agent = api.getAgent("people", id);
            if (agent.display_name !== undefined && agent.display_name.sort_name !== undefined) {
                sa_names.push(agent.display_name.sort_name);
            }
        })
        if (ca_uri.length > 0 || sa_uri.length > 0) {
            rowCount++;
            appendRow(workbook, uri, title, ca_uri.join("; "), ca_names.join("; "), sa_uri.join("; "), sa_names.join("; "));
        }
        if ((i % 100) === 1) {
            console.log("Processed", i, "/", totalAccessionIDs, "accessions,", rowCount, "roles found");
            xlsx.write(workbookName, workbook);
        }
    }

    console.log("authenitcating", os.getEnv("CAIT_API_URL"));
    res = api.login();
    if (res.isAuth !== true) {
        console.log("Login failed");
        os.exit(1);
    }


    console.log("getting accession ids");
    accessionIDs = api.listAccessions(2);
    totalAccessionIDs = accessionIDs.length;
    console.log("Processed", 0, "/", totalAccessionIDs, "accessions,", rowCount, "roles found");
    accessionIDs.forEach(addAccession);
    console.log("Processed", accessionIDs.length, "/", totalAccessionIDs, "accessions,", rowCount, "roles found");
    console.log("Writing out accessions-agent-roles.xlsx")
    xlsx.write(workbookName, workbook);
}());
