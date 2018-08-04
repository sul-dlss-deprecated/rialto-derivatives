package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestResourceToDoc(t *testing.T) {
	indexer := NewCompositeIndexer()
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)

	assert.Equal(t, []string{"Hello world!"}, docs[0].Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", docs[0].Get("id"))

}

func TestUntypedResourceToDoc(t *testing.T) {
	indexer := NewCompositeIndexer()
	data := make(map[string][]rdf.Term)
	title, _ := rdf.NewLiteral("Hello world!")
	data[models.Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)
	assert.Equal(t, []string{"Hello world!"}, docs[0].Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", docs[0].Get("id"))

}
