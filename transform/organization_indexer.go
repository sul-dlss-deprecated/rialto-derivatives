package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// OrganizationIndexer transforms organization resource types into solr Documents
type OrganizationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *OrganizationIndexer) Index(resource *models.Organization, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Organization")
	doc.Set("organization_type_ssi", resource.Subtype)
	doc.Set("title_tesi", resource.Name)
	doc.Set("abbreviation_ssim", resource.Abbreviation)
	doc.Set("parent_ssim", resource.Parent)
	return doc
}
