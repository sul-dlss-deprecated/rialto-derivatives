package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/repository"
)

func TestResourceToDoc(t *testing.T) {
	service := repository.NewService(new(MockedReader))
	indexer := NewCompositeIndexer(service)
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	book, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Book")
	created, _ := rdf.NewLiteral("2004-06-11?")
	identifier, _ := rdf.NewLiteral("1234567")

	title, _ := rdf.NewLiteral("Hello world!")
	data["id"] = id
	data["type"] = document
	data["title"] = title
	data["subtype"] = book
	data["created"] = created
	data["identifier"] = identifier

	resource := models.NewResource(data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)

	assert.Equal(t, "Hello world!", docs[0].Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", docs[0].Get("id"))
}

func TestUntypedResourceToDoc(t *testing.T) {
	service := repository.NewService(new(MockedReader))
	indexer := NewCompositeIndexer(service)
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	title, _ := rdf.NewLiteral("Hello world!")
	data["id"] = id
	data["title"] = title

	resource := models.NewResource(data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)
	assert.Nil(t, docs[0])
}
