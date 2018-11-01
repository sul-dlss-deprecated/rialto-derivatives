package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestConceptIndexerToDoc(t *testing.T) {
	indexer := &ConceptIndexer{}
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://www.w3.org/2004/02/skos/core#Concept")
	label, _ := rdf.NewLiteral("animals")
	data["id"] = id
	data["type"] = document
	data["label"] = label

	resource := models.NewResource(data)
	in := make(solr.Document)
	in.Set("id", "http://www.example.com/animals")
	doc := indexer.Index(resource.(*models.Concept), in)

	assert.Equal(t, "animals", doc.Get("title_tesi"))
	assert.Equal(t, "http://www.example.com/animals", doc.Get("id"))
}
