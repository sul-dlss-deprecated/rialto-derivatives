package models

import "github.com/knakk/rdf"

// Concept represents a subject area or topic. Works, grants, or departments
// may be associated with a Concept. Agents may have a research area that is a Concept.
type Concept struct {
	URI   string
	Label string
}

// NewConcept instantiates a concept from sparql results
func NewConcept(data map[string]rdf.Term) *Concept {
	obj := &Concept{
		URI:   data["id"].String(),
		Label: data["label"].String(),
	}

	return obj
}

// Subject returns the resources Subject URI
func (c *Concept) Subject() string {
	return c.URI
}

// Name returns the resources Name
func (c *Concept) Name() string {
	return c.Label
}
