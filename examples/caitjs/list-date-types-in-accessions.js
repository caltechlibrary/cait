
api.login();
var ids = api.listAccessions(2),
    dateTypeList = [],
    collectDateType = function (id) {
        var accession = api.getAccession(2, id),
            dtype = "",
            i = 0;
        if (accession.dates !== undefined && accession.dates.length > 0) {
            for (i = 0; i < accession.dates.length; i++ ) {
                if (dateTypeList.indexOf(accession.dates[i].date_type) < 0) {
                    console.log(accession.dates[i].date_type);
                    dateTypeList.push(accession.dates[i].date_type);
                }
            }
        }
    };

ids.forEach(collectDateType);
dateTypeList.sort();
console.log("Date types:");
console.log(JSON.stringify(dateTypeList, null, "\t"));
