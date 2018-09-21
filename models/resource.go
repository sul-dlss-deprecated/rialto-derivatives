package models

import (
	"fmt"

	"github.com/knakk/rdf"
)

// Resource is the interface type for resources
type Resource interface {
	Subject() string
}

// NewResource creates a new instance of the resource
func NewResource(data map[string]rdf.Term) Resource {
	if data["type"] == nil {
		return nil
	}
	switch t := data["type"].String(); t {
	case "http://xmlns.com/foaf/0.1/Organization":
		return NewOrganization(data)
	case "http://xmlns.com/foaf/0.1/Person":
		return NewPerson(data)
	case "http://vivoweb.org/ontology/core#Grant":
		return NewGrant(data)
	case "http://www.w3.org/2004/02/skos/core#Concept":
		return NewConcept(data)
	case "http://purl.org/ontology/bibo/Document":
		return NewPublication(data)
	case "http://xmlns.com/foaf/0.1/Project":
		return NewProject(data)
	default:
		panic(fmt.Errorf("No type for %v", data))
	}
}
