package models

import "github.com/knakk/rdf"

// Grant is an award for some project(s) or work(s), usually attached to one or
// more lead agents (PIs) whether people or departments, and awarded or funded
// by an organization or agency.
type Grant struct {
	URI  string
	Name string
}

// NewGrant instantiates a grant from sparql results
func NewGrant(data map[string]rdf.Term) *Grant {
	obj := &Grant{
		URI:  data["id"].String(),
		Name: data["name"].String(),
	}

	return obj
}

// Subject returns the resources Subject URI
func (c *Grant) Subject() string {
	return c.URI
}
