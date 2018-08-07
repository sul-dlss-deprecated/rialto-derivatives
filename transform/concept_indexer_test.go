package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestConceptIndexerToDoc(t *testing.T) {
	indexer := &ConceptIndexer{}
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://www.w3.org/2004/02/skos/core#Concept")
	title, _ := rdf.NewLiteral("animals")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{title}

	resource := models.NewResource("http://www.example.com/animals", data)
	in := make(solr.Document)
	in.Set("id", "http://www.example.com/animals")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"animals"}, doc.Get("title_tesi"))
	assert.Equal(t, "http://www.example.com/animals", doc.Get("id"))
}
