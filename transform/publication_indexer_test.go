package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func createPublicationResource() models.Resource {
	data := make(map[string]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	book, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Book")

	title, _ := rdf.NewLiteral("Hello world!")
	abstract1, _ := rdf.NewLiteral("Abstract 1")
	cites1, _ := rdf.NewIRI("http://example.com/paper1")
	link1, _ := rdf.NewIRI("http://example.com/link1")
	id1, _ := rdf.NewIRI("http://example.com/ident1")
	doi1, _ := rdf.NewIRI("https://doi.org/10.1109/5.771073")
	id, _ := rdf.NewIRI("http://example.com/record1")
	data["id"] = id
	data["type"] = document
	data["subtype"] = book
	data["title"] = title
	data["abstract"] = abstract1
	data["cites"] = cites1
	data["uri"] = link1
	data["identifier"] = id1
	data["doi"] = doi1

	return models.NewResource(data)
}

func TestPublicationResourceToDoc(t *testing.T) {
	indexer := &PublicationIndexer{}

	resource := createPublicationResource()
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource.(*models.Publication), in)

	assert.Equal(t, "Hello world!", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
	assert.Equal(t, "https://doi.org/10.1109/5.771073", doc.Get("doi_ssim"))

}
