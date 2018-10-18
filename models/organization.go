package models

import "github.com/knakk/rdf"

// Organization is a non-person agent. It may represent a University, School or Department
type Organization struct {
	URI          string
	Subtype      string
	Name         string
	Abbreviation string
	Parent       string
}

// NewOrganization instantiates an organization from sparql results
func NewOrganization(data map[string]rdf.Term) *Organization {
	org := &Organization{
		URI:  data["id"].String(),
		Name: data["name"].String(),
	}

	if subtype := data["subtype"]; subtype != nil {
		org.Subtype = subtype.String()
	}

	if abbrev := data["abbreviation"]; abbrev != nil {
		org.Abbreviation = abbrev.String()
	}

	if parent := data["parent"]; parent != nil {
		org.Parent = parent.String()
	}
	return org
}

// Subject returns the resources Subject URI
func (c *Organization) Subject() string {
	return c.URI
}
