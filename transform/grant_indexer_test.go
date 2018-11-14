package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestGrantIndexerToDoc(t *testing.T) {
	indexer := &GrantIndexer{}
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Grant")
	name, _ := rdf.NewLiteral("Hydra in a Box")
	pi, _ := rdf.NewIRI("http://example.com/record2")
	piName, _ := rdf.NewLiteral("Bob")
	assigned, _ := rdf.NewIRI("http://example.com/record3")
	assignedName, _ := rdf.NewLiteral("Chocolate Foundation")
	start, _ := rdf.NewLiteral("2018-05-14")

	data["id"] = id
	data["type"] = document
	data["name"] = name
	data["pi"] = pi
	data["pi_label"] = piName
	data["assigned"] = assigned
	data["assigned_label"] = assignedName
	data["start"] = start

	resource := models.NewGrant(data)
	resource.Identifiers = []string{"123456"}
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Hydra in a Box", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
	assert.Equal(t, "http://example.com/record2", doc.Get("pi_ssim"))
	assert.Equal(t, "Bob", doc.Get("pi_label_tsim"))
	assert.Equal(t, "http://example.com/record3", doc.Get("assigned_ssim"))
	assert.Equal(t, "Chocolate Foundation", doc.Get("assigned_label_tsim"))
	assert.Equal(t, "2018-05-14", doc.Get("start_date_ss"))
	assert.Equal(t, nil, doc.Get("end_date_ss"))
	assert.Equal(t, []string{"123456"}, doc.Get("identifiers_ssim"))
}
