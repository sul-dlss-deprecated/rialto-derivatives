package models

import "github.com/knakk/rdf"

// Publication is a representation of articles, research outputs, datasets, etc.
// If feasible, there should be a link to manifestations of that Work (i.e. DOI).
type Publication struct {
	URI     string
	Subtype string
	Title   string
	DOI     string
}

// NewPublication instantiates a publication from sparql results
func NewPublication(data map[string]rdf.Term) *Publication {
	pub := &Publication{
		URI:     data["id"].String(),
		Subtype: data["subtype"].String(),
		Title:   data["title"].String(),
	}

	if doi := data["doi"]; doi != nil {
		pub.DOI = doi.String()
	}
	return pub
}

// Subject returns the resources Subject URI
func (c *Publication) Subject() string {
	return c.URI
}
