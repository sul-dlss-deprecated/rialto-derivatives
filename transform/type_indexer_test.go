package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestSingleTypeToDoc(t *testing.T) {
	indexer := &TypeIndexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}

	resource := models.NewResource("http://example.com/record1", data)
	doc := indexer.Index(resource, make(solr.Document))

	assert.Equal(t, "http://purl.org/ontology/bibo/Document", doc.Get("type_ssi"))
}

func TestMultipleTypeToDoc(t *testing.T) {
	indexer := &TypeIndexer{}

	data := make(map[string][]rdf.Term)
	agent, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Agent")
	person, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{agent, person}

	resource := models.NewResource("http://example.com/record1", data)
	doc := indexer.Index(resource, make(solr.Document))
	assert.Equal(t, "http://xmlns.com/foaf/0.1/Person", doc.Get("type_ssi"))

}
