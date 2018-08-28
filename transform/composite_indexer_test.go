package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestResourceToDoc(t *testing.T) {
	service := repository.NewService(new(MockedReader))
	indexer := NewCompositeIndexer(service)
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	book, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Book")

	title, _ := rdf.NewLiteral("Hello world!")
	data["id"] = id
	data["type"] = document
	data["title"] = title
	data["subtype"] = book

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
