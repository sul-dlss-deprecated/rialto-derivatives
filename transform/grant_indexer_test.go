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
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Grant")
	title, _ := rdf.NewLiteral("Hydra in a Box")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Hydra in a Box"}, doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
