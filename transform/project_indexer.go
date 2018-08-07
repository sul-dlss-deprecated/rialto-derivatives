package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// ProjectIndexer transforms organization resource types into solr Documents
type ProjectIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *ProjectIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"title":          "title_tesi",
		"alternateTitle": "alternative_title_tesim",
		"hasStartDate":   "start_date_ssi",
		"hasEndDate":     "end_date_ssi",
	}

	doc = indexMapping(resource, doc, mapping)

	return doc
}
