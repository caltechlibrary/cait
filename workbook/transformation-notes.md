
## Analysis

Basic structure

+ <ead>
    + <eadheader> ... (comes from "front matter")
    + <archdesc>
        + <did> (from "front matter")
        + <accessrestrict> (from "front matter")
        + <unrestricted> (from "front matter")
        + <prefercite> (from "front matter")
        + <dsc> this is where our Excel workbook files come into play.


### Kuppermann

Comparing Kuppermann.xlsx and finding aid at http://www.oac.cdlib.org/findaid/ark:/13030/c85q51hb/?query=Kupperman

Three Steps to EAD document

1. Open the HTML page to see the finding aid (this is built directly from the submitted EAD)
2. View Source on the HTML finding aid, look for the "Request Items" link, the link to the XML document is passed in a parameter of "Value"
3. Open the XML link and save to disc http://voro.cdlib.org/oac-ead/prime2002/caltech/kupperma.xml

Excel workbook is contained inside <dsc> element with a a <head> with "Container List".
Each sheet in the Workbook represents a <c01>.

+ Row three is used to populate the <c01> element
    + Arrangement becomes the <unitid/> in <did>
    + Title/Nominations and Recommendations become <unittitle/> in <did>
+ Subseries becomes <c02>
    + Arrangement becomes the <unitid/>
    + Title/Nomintaitons and Recommentdations become <unittitle/> in <did>
+ Box info becomes <c03>
    + Box number becomes <container type="box" label="Box ">$box_no$</container> in <did>
    + Folder No becomes <container type="folder" label="Folder ">$folder_no$</container> in <did>
    + Title/Nominations and Recommendations become <unittitle/> in <did>
    + Date becomes <unitdate> in <did>
+ Header becomes it's own <c03><did> as <unittitle/>

### Rose Bowl Hoax

Comparing RoseBowlHoax.xlsx and finding aid at http://www.oac.cdlib.org/findaid/ark:/13030/c86113qr/entire_text/?query=Rose%20Bowl%20Hoax

Three Steps to EAD document

1. Open the HTML page to see the finding aid (this is built directly from the submitted EAD)
2. View Source on the HTML finding aid, look for the "Request Items" link, the link to the XML document is passed in a parameter of "Value"
3. Open the XML link and save to disc http://voro.cdlib.org/oac-ead/prime2002/caltech/rosebowl.xml

First sheet "Container-List" creates a series of <c01> elements with a single <did>. The series ID column's value is 0 as
is subseries. It looks like we can use the values in box, folder, series, subseries to drive the type of "c" element the
row's data moves into.

Map for <unittitle>, <unitdate> appear the same as found in Kuppermann in the <did> evan at <c01>. If PhysDesc column is
populated then it becomes <physdesc> element in the <did>. Note column becomes <note> with the "c" level element that is
active (i.e. <c01> for RoseBowlHoax.xlsx)


## Questions

+ In <c02> elements the attributes are id, level, score, C-ORDER. How are these defined
    + id apprears to be "id" concatenated with a auto incrementing value
    + level apprears to be always "file" in a <co3>, may also be "series" at <c01> and "subseries" at <c02>

workbook package is focused on processing the container descriptions that come from the formatted spreadsheets entered by ArchivesSpace. They
are transformed into an EAD then imported into ArchivesSpace.  Front matter is pulled from Accession records covering the materials.

A list of EADs available from OAC for Caltech can be seen at http://www.oac.cdlib.org/institutions/California+Institute+of+Technology

