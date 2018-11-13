package transform

import (
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// GrantIndexer transforms grant resource types into solr Documents
type GrantIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *GrantIndexer) Index(resource *models.Grant, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Grant")
	doc.Set("title_tesi", resource.Name)
	if resource.PI != nil {
		doc.Set("pi_label_tsim", resource.PI.Label)
		doc.Set("pi_ssim", resource.PI.URI)
	}
	if resource.Assigned != nil {
		doc.Set("assigned_label_tsim", resource.Assigned.Label)
		doc.Set("assigned_ssim", resource.Assigned.URI)
	}
	if resource.Start != "" {
		doc.Set("start_date_ss", resource.Start)
	}
	if resource.End != "" {
		doc.Set("end_date_ss", resource.End)
	}
	if resource.Identifiers != nil {
		doc.Set("identifiers_ssim", resource.Identifiers)
	}

	return doc
}
