package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// OrganizationIndexer transforms organization resource types into solr Documents
type OrganizationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *OrganizationIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"orgName": "title_tesi",
	}

	doc = indexMapping(resource, doc, mapping)

	return doc
}
