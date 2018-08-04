package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestDefaultIndexerToDoc(t *testing.T) {
	indexer := &DefaultIndexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Hello world!"}, doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
