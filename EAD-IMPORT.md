
# cait

## ead import

The prior system used Excel spreadsheets in conjunction with SCread to create EADs.  _cait_ supports
a similar workflow leveraging the _xlsximporter_ tool.  

## Excel Spreadsheets

The Excel spreadsheets need to be in the ".xlsx" XML based format (aka. Workbook).  This is file
is expected to be a series of sheets with the following header row --

+ Box 
+ Folder
+ Arrangement 
+ Title (& see also)|Nominations and Recommentdations
+ Dates
+ See/Header
+ Series
+ Subseries
+ Lists
+ PhysDesc
+ Notes
+ Digital Archival Object role
+ DAO Href
+ DAO title
+ DAO Description

Each row is a record with empty fields holding the values defined by previous rows. Example would be if
row five indicated the "arrangement" was "box 34" then row 6 would be assumed to refer to "box 34" until
such time as a new value in the arrangement column was encountered.

### Sheetnames

There are a few reserved sheetnames.

+ Config - this sheet is used to link to archivesspace data, it includes things like relevant accession ids
+ Tools - this is used by the scread process
+ SCREAD-DACS X-WALK - this is used by the scread process
+ Comments - this is used by the scread process

With the exception of *Config* the other sheets are just ignored.


# cait

## Ead importer

### Legacy Workflow

The prior system used Excel spreadsheets in conjunction with SCread to create EADs.  _cait_ may support
a similar workflow leveraging the _xlsximporter_ tool.

### Legacy EADs

Caltech EADs hosted at California Online Archives can be downloaded from http://voro.cdlib.org/oac-ead/prime2002/caltech/
These EADs are in older formats with some not being passing the import process in ArchivesSpace 1.4.2 and 1.5.0.

## Ead notes structures found by viewing a sample at OAC

### EADHeader

#### FileDesc

+ FileDesc
    + Titlestmt
        + []Titleproper (Example: primary title and filing type titles)
        + []author (Example: a single comma delimited list of who processed the EAD)
    + Pulbicationstmt
        + Publisher (Example: Caltech Archives)
        + Address
            + []Addressline (multiple lines of official address, Phone, Fax, Email, and URL to website)
        + Date 
        + P (Example: California Institute of Technology. All rights reserved.)

#### ProfileDesc

+ ProfileDesc
    + Creation
    + LangUsage

#### RevisionDesc

+ RevisionDesc
    + []Change
        + Date
        + Item

### FrontMatter

### ArchDesc

## CA Excel Workbook Organization

The Excel spreadsheets need to be in the ".xlsx" XML based format (aka. Workbook).  This is file
is expected to be a series of sheets with the following header row (winds up be DID in ArchDesc) --

+ Box  (int)
+ Folder (int)
+ Arrangement 
    + Series: concat(Arrangement, Title)
        + Sub-sereies: concat(Arrangement, Title)
            + Box: concat(Arrangement, Title
            + Header (optional): SeeHeader()
            + Note (optional): Title, Notes
            + Subscript (optional): Title
            + see also: Title
            + (other): Title
        + Box: concat(Arrangement, Title
        + Header (optional): SeeHeader()
        + Note (optional): Title, Notes
        + Subscript (optional): Title
        + see also: Title
        + (other): Title
+ Title (& see also)|Nominations and Recommentdations
+ Dates
+ See/Header (attach to the folder being described in the row, continue until new header is defined)
+ Series
+ Subseries (if zero, no subseries is defined)
+ Lists
+ PhysDesc
+ Notes
+ Digital Archival Object role
+ DAO Href
+ DAO title
+ DAO Description

Each row is a record with empty fields holding the values defined by previous rows. Example would be if
row five indicated the "arrangement" was "box 34" then row 6 would be assumed to refer to "box 34" until
such time as a new value in the arrangement column was encountered.

### Sheetnames

There are a few reserved sheetnames.

+ Config - this sheet is used to link to archivesspace data, it includes things like relevant accession ids
+ Tools - this is used by the scread process
+ SCREAD-DACS X-WALK - this is used by the scread process
+ Comments - this is used by the scread process

With the exception of *Config* the other sheets are just ignored.

