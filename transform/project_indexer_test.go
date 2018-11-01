package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestProjectIndexerToDoc(t *testing.T) {
	indexer := &ProjectIndexer{}
	data := make(map[string]rdf.Term)
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Project")
	title, _ := rdf.NewLiteral("Rialto")
	id, _ := rdf.NewIRI("http://example.com/record1")
	data["id"] = id
	data["type"] = document
	data["title"] = title

	resource := models.NewResource(data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource.(*models.Project), in)

	assert.Equal(t, "Rialto", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
