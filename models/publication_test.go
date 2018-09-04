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

	data["id"] = id
	data["title"] = title

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
	doi, _ := rdf.NewLiteral("10.1073/pnas.1103105108")
	abstract, _ := rdf.NewLiteral("Regulation of gene expression at the transcriptional level is achieved by complex interactions")

	data["id"] = id
	data["subtype"] = book
	data["title"] = title
	data["doi"] = doi
	data["abstract"] = abstract
	resource := NewPublication(data)
	assert.IsType(t, &Publication{}, resource)
	assert.Equal(t, book.String(), *resource.Subtype)
	assert.Equal(t, *resource.DOI, doi.String())
	assert.Equal(t, *resource.Abstract, abstract.String())
}
