package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// ConceptIndexer transforms concept/topic resource types into solr Documents
type ConceptIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *ConceptIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"label":            "title_tesi",
		"alternativeTitle": "alternative_title_tesim",
	}

	doc = indexMapping(resource, doc, mapping)

	return doc
}
