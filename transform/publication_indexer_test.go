package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func createPublicationResource() models.Resource {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")
	abstract1, _ := rdf.NewLiteral("Abstract 1")
	abstract2, _ := rdf.NewLiteral("Abstract 2")
	cites1, _ := rdf.NewIRI("http://example.com/paper1")
	cites2, _ := rdf.NewIRI("http://example.com/paper2")
	link1, _ := rdf.NewIRI("http://example.com/link1")
	link2, _ := rdf.NewIRI("http://example.com/link2")
	id1, _ := rdf.NewIRI("http://example.com/ident1")
	id2, _ := rdf.NewIRI("http://example.com/ident2")
	doi1, _ := rdf.NewIRI("https://doi.org/10.1109/5.771073")
	doi2, _ := rdf.NewIRI("https://doi.org/10.1000/182")

	data[models.RdfTypePredicate] = []rdf.Term{document}
	data[models.TitlePredicate] = []rdf.Term{title}
	data[models.AbstractPredicate] = []rdf.Term{abstract1, abstract2}
	data[models.CitesPredicate] = []rdf.Term{cites1, cites2}
	data[models.LinkPredicate] = []rdf.Term{link1, link2}
	data[models.IdentifierPredicate] = []rdf.Term{id1, id2}
	data[models.DoiPredicate] = []rdf.Term{doi1, doi2}

	return models.NewResource("http://example.com/record1", data)
}

func TestPublicationResourceToDoc(t *testing.T) {
	indexer := &PublicationIndexer{}

	resource := createPublicationResource()
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Hello world!"}, doc.Get("title_ssi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
	assert.Equal(t, []string{"https://doi.org/10.1109/5.771073", "https://doi.org/10.1000/182"}, doc.Get("doi_ssim"))

}
