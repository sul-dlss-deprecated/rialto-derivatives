package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// DefaultIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	if resource.Titles() != nil {
		doc.Set("title_ssi", resource.Titles()[0].String())
	}
	return doc
}
