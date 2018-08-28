package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// ConceptIndexer transforms concept/topic resource types into solr Documents
type ConceptIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *ConceptIndexer) Index(resource *models.Concept, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Concept")
	doc.Set("title_tesi", resource.Label)
	doc.Set("alternative_title_tesim", resource.AlternativeLabel)
	return doc
}
