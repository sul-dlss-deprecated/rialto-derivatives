package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// GrantIndexer transforms grant resource types into solr Documents
type GrantIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *GrantIndexer) Index(resource *models.Grant, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Grant")
	doc.Set("title_tesi", resource.Name)
	return doc
}
