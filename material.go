package main

// This defines the XML structure of an .osm changset file
type Osm struct {
	Changesets []Changeset `xml:"changeset"`
}

type Changeset struct {
	Id   int64  `xml:"id,attr"`
	User string `xml:"user,attr"`
	Tags []Tag  `xml:"tag"`
}

type Tag struct {
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

// This defines the known editors later used for analysis of the changsets
var (
	unknownEditor = "_UNKNOWN"
	noEditor      = "_NO_EDITOR"
	knownEditors  = []string{
		"josm",
		"id",
		"potlatch",
		"maps.me",
		"osmand+",
		"vespucci",
		"streetcomplete",
		"osmtools",
		"merkaartor",
		"osm2go",
		unknownEditor,
		noEditor,
	}
)