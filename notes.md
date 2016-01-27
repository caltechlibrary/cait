Importing Digital Objects from a Excel spreadsheet with the following columns.

+ "Digital Object ID" Column A
+ "Title (map to title)" Column B
+ "Series" Column C
+ "Keywords" Column D
+ "Name \_and Subject" Column E
+ "url_online oral history URL" Column F
+ "Oral History Text by Item ID::Text Description"
+ "Oral History Text by Item ID::Search Text"

The Mapping follows these rules:

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
+ "Oral History Text by Item ID::Text Description" maps to Notes -> General Notes:Content
+ "Oral History Text by Item ID::Text Description" maps to Notes -> General Note:Label

All Publish fields should be true

Note: Need to add meta data to add "cards" results on Google Search.
