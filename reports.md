
# Building reports out of dataset exported from ArchivesSpace

## required

+ Bash, grep, sed
+ [cait](https://caltechlibrary.github.io/cait/)
+ [datatools](https://caltechlibrary.github.io/datatools/)


## Accessions Report

+ URL to edit accession
+ id
+ id_0, id_1
+ .display_string
+ .title
+ .extents (an array of extent objects)
    + extents[i].extent_type
    + extents[i].physical_details
+ .provenance
+ .subjects (an array pointing at another datset)
    + .subjects[i].ref is path to subjects object
+ .dates (an array of date objects)
    + .dates[i].data_type
    + .dates[i].label
    + .dates[i].expression
+ user_defined??
+ access_restrictions??
+ use_restrictions??
+ general_note??
