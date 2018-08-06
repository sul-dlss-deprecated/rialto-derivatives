package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PersonIndexer transforms person resource types into solr Documents
type PersonIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PersonIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"title": "title_tesi",
	}

	doc = indexMapping(resource, doc, mapping)

	return doc
}
