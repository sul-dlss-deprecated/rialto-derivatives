package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestOrganizationIndexerToDoc(t *testing.T) {
	indexer := &OrganizationIndexer{}
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	organization, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Division")
	title, _ := rdf.NewLiteral("Women's Fencing Program")
	parent, _ := rdf.NewIRI("http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs")

	data["id"] = id
	data["type"] = organization
	data["subtype"] = document
	data["name"] = title
	data["parent"] = parent

	resource := models.NewResource(data)
	in := make(solr.Document)
	in.Set("id", "http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs/womens-fencing-program")
	doc := indexer.Index(resource.(*models.Organization), in)

	assert.Equal(t, "Women's Fencing Program", doc.Get("title_tesi"))
	assert.Equal(t, "http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs", doc.Get("parent_ssim"))
	assert.Equal(t, "http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs/womens-fencing-program", doc.Get("id"))
}
