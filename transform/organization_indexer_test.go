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
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Division")
	title, _ := rdf.NewLiteral("Women's Fencing Program")
	parent, _ := rdf.NewIRI("http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs")
	abbrev, _ := rdf.NewLiteral("LIAS")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{title}
	data[models.Predicates["obo"]["BFO_0000050"]] = []rdf.Term{parent}
	data[models.Predicates["vivo"]["abbreviation"]] = []rdf.Term{abbrev}

	resource := models.NewResource("http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs/womens-fencing-program", data)
	in := make(solr.Document)
	in.Set("id", "http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs/womens-fencing-program")
	doc := indexer.Index(resource, in)

	assert.Equal(t, []string{"Women's Fencing Program"}, doc.Get("title_tesi"))
	assert.Equal(t, []string{"http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs"}, doc.Get("parent_ssim"))
	assert.Equal(t, []string{"LIAS"}, doc.Get("abbreviation_ssim"))
	assert.Equal(t, "http://rialto.stanford.edu/organizations/department-of-athletics-physical-education-and-recreation/womens-sport-programs/womens-fencing-program", doc.Get("id"))
}
