package models

import (
	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
)

// Person is a human actor involved in creating works
type Person struct {
	URI              string
	Subtype          string
	Firstname        string
	Lastname         string
	Organization     *string // URI
	DepartmentLabel  *string
	DepartmentURI    *string // URI
	SchoolLabel      *string
	SchoolURI        *string // URI
	InstitutionLabel *string
	InstitutionURI   *string // URI
}

// NewPerson instantiates a person from sparql results
func NewPerson(data map[string]rdf.Term) *Person {
	obj := &Person{
		URI:     data["id"].String(),
		Subtype: data["subtype"].String(),
	}
	if firstname := data["firstname"]; firstname != nil {
		obj.Firstname = firstname.String()
	}

	if lastname := data["lastname"]; lastname != nil {
		obj.Lastname = lastname.String()
	}

	if organization := data["org"]; organization != nil {
		str := organization.String()
		obj.Organization = &str
	}

	return obj
}

// Subject returns the resources Subject URI
func (c Person) Subject() string {
	return c.URI
}

// SetOrganizationInfo allow organization relationships to be passed in
func (c *Person) SetOrganizationInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		name := solution["name"].String()
		uri := solution["org"].String()
		switch solution["type"].String() {
		case "http://vivoweb.org/ontology/core#Department":
			c.DepartmentLabel = &name
			c.DepartmentURI = &uri
		case "http://vivoweb.org/ontology/core#School":
			c.SchoolLabel = &name
			c.SchoolURI = &uri
		case "http://vivoweb.org/ontology/core#University":
			c.InstitutionLabel = &name
			c.InstitutionURI = &uri
		}
	}
}
