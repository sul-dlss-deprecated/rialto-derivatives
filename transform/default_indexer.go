package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// DefaultIndexer transforms unknown resource types to solr Documents
type DefaultIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *DefaultIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	if resource.Titles() != nil {
		doc.Set("title_ssi", resource.Titles()[0].String())
	}
	return doc
}
