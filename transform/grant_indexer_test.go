package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestGrantIndexerToDoc(t *testing.T) {
	indexer := &GrantIndexer{}
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Grant")
	name, _ := rdf.NewLiteral("Hydra in a Box")
	data["id"] = id
	data["type"] = document
	data["name"] = name

	resource := models.NewResource(data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource.(*models.Grant), in)

	assert.Equal(t, "Hydra in a Box", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
