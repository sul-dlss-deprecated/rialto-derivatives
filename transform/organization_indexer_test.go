package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestOrganizationIndexerToDoc(t *testing.T) {
	indexer := &OrganizationIndexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	title, _ := rdf.NewLiteral("Stanford")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Stanford"}, doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
