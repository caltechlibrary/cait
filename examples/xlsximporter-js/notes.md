This is an example JavaScript file for importing Digital Objects from a
Excel spreadsheet with the following columns.

+ Digital Object ID
+ Title (map to title)
+ Series
+ Keywords
+ Name \_and Subject
+ url_online oral history URL
+ Oral History Text by Item ID::Text Description
+ Oral History Text by Item ID::Search Text

The Mapping follows these rules:

+ Digital Object ID maps to URI with for /repositories/2/digital_obejcts/ID
+ Title maps to title
+ "Series" maps to subject of type function
+ "Keywords" maps to subject of type topical
+ "Name \_and Subject" map to creator or subject based on the content of "Series"
    + Oral History -> Creator (interviewer/interviewee), Subject (interviewee)
    + Film & Video -> Subject
    + Institute Publications -> Creator
    + Manuscript Collection -> Creator
    + Watson Lecture -> Creator
    + Alumni Day -> Creator
    + Commencement -> Creator
    + Institute Publications -> Creator
    + "" -> Creator
+ "url_online oral history URL" maps to unique identifier for Digital Object, File Versions-> File URI (publish should be checked), Notes -> General Notes:Persistent ID
+ "Oral History Text by Item ID::Search Text" maps to Notes -> General Notes:Content
+ "Oral History Text by Item ID::Text Description" maps to Notes -> General Note:Label

All Publish fields should be true
