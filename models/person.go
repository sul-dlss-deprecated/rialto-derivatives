package models

import (
	"log"

	"github.com/knakk/rdf"
)

// Person is a human actor involved in creating works
type Person struct {
	URI        string
	Subtype    string
	Firstname  string
	Lastname   string
	Department *string // URI
}

// NewPerson instantiates a person from sparql results
func NewPerson(data map[string]rdf.Term) *Person {
	obj := &Person{
		URI:     data["id"].String(),
		Subtype: data["subtype"].String(),
	}
	log.Printf("Creating person with %v", data)

	if firstname := data["firstname"]; firstname != nil {
		obj.Firstname = firstname.String()
	}

	if lastname := data["lastname"]; lastname != nil {
		obj.Lastname = lastname.String()
	}

	if department := data["department"]; department != nil {
		str := department.String()
		obj.Department = &str
	}

	return obj
}

// Subject returns the resources Subject URI
func (c *Person) Subject() string {
	return c.URI
}
