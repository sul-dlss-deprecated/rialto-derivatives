package models

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
)

func TestPublicationResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	subtype, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Book")
	title, _ := rdf.NewLiteral("Hello world!")
	created, _ := rdf.NewLiteral("2004-06-11?")
	identifier, _ := rdf.NewLiteral("1234567")

	data["id"] = id
	data["type"] = document
	data["subtype"] = subtype
	data["title"] = title
	data["created"] = created
	data["identifier"] = identifier

	resource := NewResource(data)
	assert.IsType(t, &Publication{}, resource)
}

func TestPersonResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")
	fname, _ := rdf.NewLiteral("Justin")
	lname, _ := rdf.NewLiteral("Coyne")
	student, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Student")

	data["id"] = id
	data["type"] = document
	data["subtype"] = student
	data["firstname"] = fname
	data["lastname"] = lname

	resource := NewResource(data)
	assert.IsType(t, &Person{}, resource)
}

func TestOrganizationResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record2")
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	subtype, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Department")
	abbreviation, _ := rdf.NewLiteral("FOO")
	parent, _ := rdf.NewIRI("http://example.com/record1")

	name, _ := rdf.NewLiteral("Cornell")

	data["id"] = id
	data["type"] = document
	data["subtype"] = subtype
	data["name"] = name
	data["abbreviation"] = abbreviation
	data["parent"] = parent
	resource := NewResource(data)
	assert.IsType(t, &Organization{}, resource)
}

func TestGrantResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Grant")
	name, _ := rdf.NewLiteral("Hydra in a Box")

	data["id"] = id
	data["type"] = document
	data["name"] = name
	resource := NewResource(data)
	assert.IsType(t, &Grant{}, resource)
}

func TestProjectResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Project")
	title, _ := rdf.NewLiteral("Hydra in a Box")
	alt, _ := rdf.NewLiteral("Hybox")
	startdate, _ := rdf.NewLiteral("Oct 1st")
	enddate, _ := rdf.NewLiteral("23rd century")

	data["id"] = id
	data["type"] = document
	data["title"] = title
	data["alternativetitle"] = alt
	data["startdate"] = startdate
	data["enddate"] = enddate
	resource := NewResource(data)
	assert.IsType(t, &Project{}, resource)
}

func TestConceptResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://www.w3.org/2004/02/skos/core#Concept")
	label, _ := rdf.NewLiteral("Hydra in a Box")

	data["id"] = id
	data["type"] = document
	data["label"] = label
	resource := NewResource(data)
	assert.IsType(t, &Concept{}, resource)
}
