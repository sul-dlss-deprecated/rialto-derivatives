package models

import (
	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
)

// Grant is an award for some project(s) or work(s), usually attached to one or
// more lead agents (PIs) whether people or departments, and awarded or funded
// by an organization or agency.
type Grant struct {
	URI         string
	Name        string
	Assigned    *Labeled
	PI          *Labeled
	Identifiers []string
	Start       string
	End         string
}

// NewGrant instantiates a grant from sparql results
func NewGrant(data map[string]rdf.Term) *Grant {
	obj := &Grant{
		URI:         data["id"].String(),
		Name:        data["name"].String(),
		Identifiers: []string{},
	}

	obj.PI = &Labeled{data["pi"].String(), data["pi_label"].String()}

	if data["assigned"] != nil {
		obj.Assigned = &Labeled{data["assigned"].String(), data["assigned_label"].String()}
	}

	if start := data["start"]; start != nil {
		obj.Start = start.String()
	}

	if end := data["end"]; end != nil {
		obj.End = end.String()
	}

	return obj
}

// Subject returns the resources Subject URI
func (c *Grant) Subject() string {
	return c.URI
}

// SetIdentifierInfo allow identifiers to be passed in
func (c *Grant) SetIdentifierInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		identifier := solution["id"].String()
		c.Identifiers = append(c.Identifiers, identifier)
	}
}
