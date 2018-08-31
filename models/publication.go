package models

import "github.com/knakk/rdf"

// Publication is a representation of articles, research outputs, datasets, etc.
// If feasible, there should be a link to manifestations of that Work (i.e. DOI).
type Publication struct {
	URI     string
	Subtype *string
	Title   string
	DOI     *string
}

// NewPublication instantiates a publication from sparql results
func NewPublication(data map[string]rdf.Term) *Publication {
	pub := &Publication{
		URI:   data["id"].String(),
		Title: data["title"].String(),
	}

	if subtype := data["subtype"]; subtype != nil {
		subtypeStr := subtype.String()
		pub.Subtype = &subtypeStr
	}

	if doi := data["doi"]; doi != nil {
		doiStr := doi.String()
		pub.DOI = &doiStr
	}
	return pub
}

// Subject returns the resources Subject URI
func (c *Publication) Subject() string {
	return c.URI
}
