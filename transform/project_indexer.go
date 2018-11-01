package transform

import (
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// ProjectIndexer transforms organization resource types into solr Documents
type ProjectIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *ProjectIndexer) Index(resource *models.Project, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Project")
	doc.Set("title_tesi", resource.Title)
	doc.Set("alternative_title_tesim", resource.AlternativeTitle)
	doc.Set("start_date_ssi", resource.StartDate)
	doc.Set("end_date_ssi", resource.EndDate)

	return doc
}
