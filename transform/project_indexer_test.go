package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestProjectIndexerToDoc(t *testing.T) {
	indexer := &ProjectIndexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Project")
	title, _ := rdf.NewLiteral("Rialto")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Rialto"}, doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
