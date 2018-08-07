package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// GrantIndexer transforms grant resource types into solr Documents
type GrantIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *GrantIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"grantName": "title_tesi",
	}

	doc = indexMapping(resource, doc, mapping)

	return doc
}
