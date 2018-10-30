package models

import (
	"fmt"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
)

// Person is a human actor involved in creating works
type Person struct {
	URI             string
	Firstname       string
	Lastname        string
	DepartmentOrgs  []*Labeled
	SchoolOrgs      []*Labeled
	InstituteOrgs   []*Labeled
	InstitutionOrgs []*Labeled
	Countries       []*Labeled
	Subtypes        []*Labeled
}

// Labeled is a resource that has a label
type Labeled struct {
	URI   string
	Label string
}

// NewPerson instantiates a person from sparql results
func NewPerson(data map[string]rdf.Term) *Person {
	obj := &Person{
		URI: data["id"].String(),
	}

	if firstname := data["firstname"]; firstname != nil {
		obj.Firstname = firstname.String()
	}

	if lastname := data["lastname"]; lastname != nil {
		obj.Lastname = lastname.String()
	}

	return obj
}

// Subject returns the resources Subject URI
func (c Person) Subject() string {
	return c.URI
}

// Name returns the resources Name
func (c Person) Name() string {
	givenName := c.Firstname
	familyName := c.Lastname
	return fmt.Sprintf("%v %v", givenName, familyName)
}

// SetPositionOrganizationInfo adds organization relationships to a person from sparql results
func (c *Person) SetPositionOrganizationInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		org := &Labeled{solution["org"].String(), solution["name"].String()}
		if solution["type"] != nil {
			switch solution["type"].String() {
			case "http://vivoweb.org/ontology/core#Department":
				c.DepartmentOrgs = append(c.DepartmentOrgs, org)
			case "http://vivoweb.org/ontology/core#School":
				c.SchoolOrgs = append(c.SchoolOrgs, org)
			case "http://vivoweb.org/ontology/core#University":
				c.InstitutionOrgs = append(c.InstitutionOrgs, org)
			case "http://vivoweb.org/ontology/core#Institute":
				c.InstituteOrgs = append(c.InstituteOrgs, org)
			}
		} else {
			// If no type, then assuming an institution
			c.InstitutionOrgs = append(c.InstitutionOrgs, org)
		}
	}
}

// SetCountriesInfo adds countries to a person from sparql results
func (c *Person) SetCountriesInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		country := &Labeled{solution["country"].String(), solution["label"].String()}
		c.Countries = append(c.Countries, country)
	}
}

// SetPersonSubtypesInfo adds subtypes to a person from sparql results
func (c *Person) SetPersonSubtypesInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		// For label, removing http://sul.stanford.edu/rialto/ontology# from URI
		subtype := &Labeled{solution["subtype"].String(), solution["subtype"].String()[40:]}
		c.Subtypes = append(c.Subtypes, subtype)
	}
}
