package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestResourceToDoc(t *testing.T) {
	indexer := &Indexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[models.RdfTypePredicate] = []rdf.Term{document}
	data[models.TitlePredicate] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)

	assert.Equal(t, "Hello world!", docs[0].Get("title_ssi"))
	assert.Equal(t, "http://example.com/record1", docs[0].Get("id"))

}
