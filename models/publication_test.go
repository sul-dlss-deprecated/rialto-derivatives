package models

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
)

func TestNewPublicationMinimalFields(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	title, _ := rdf.NewLiteral("Hello world!")
	created, _ := rdf.NewLiteral("2004-06-11?")
	identifier, _ := rdf.NewLiteral("1234567")

	data["id"] = id
	data["title"] = title
	data["created"] = created
	data["identifier"] = identifier

	resource := NewPublication(data)
	assert.IsType(t, &Publication{}, resource)
	assert.Equal(t, title.String(), resource.Title)
	assert.Equal(t, id.String(), resource.Subject())
}

func TestNewPublicationAllFields(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	book, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Book")
	title, _ := rdf.NewLiteral("Predicting Cancer Drug Response")
	created, _ := rdf.NewLiteral("2004-06-11?")
	identifier, _ := rdf.NewLiteral("1234567")
	doi, _ := rdf.NewLiteral("10.1073/pnas.1103105108")
	abstract, _ := rdf.NewLiteral("Regulation of gene expression at the transcriptional level is achieved by complex interactions")
	publisher, _ := rdf.NewIRI("http://example.com/publishing_house")
	description, _ := rdf.NewLiteral("A scholarly paper")

	data["id"] = id
	data["subtype"] = book
	data["title"] = title
	data["created"] = created
	data["identifier"] = identifier
	data["doi"] = doi
	data["abstract"] = abstract
	data["publisher"] = publisher
	data["description"] = description
	data["created"] = created

	resource := NewPublication(data)
	assert.IsType(t, &Publication{}, resource)
	assert.Equal(t, *resource.Subtype, book.String())
	assert.Equal(t, resource.Created, created.String())
	assert.Equal(t, resource.Identifier, identifier.String())
	assert.Equal(t, *resource.DOI, doi.String())
	assert.Equal(t, *resource.Abstract, abstract.String())
	assert.Equal(t, *resource.Publisher, publisher.String())
	assert.Equal(t, *resource.Description, description.String())
	assert.Equal(t, resource.CreatedYear, 2004)

}
