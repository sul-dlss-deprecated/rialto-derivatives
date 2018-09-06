package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PublicationIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource *models.Publication, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Publication")
	doc.Set("title_tesi", resource.Title)
	if resource.DOI != nil {
		doc.Set("doi_ssim", *resource.DOI)
	}

	if resource.Abstract != nil {
		doc.Set("abstract_tesim", *resource.Abstract)
	}

	// "cites":            "cites_ssim",
	// "identifier":       "identifier_ssim",
	// "link":             "link_ssim",
	// "description":      "description_tesim",
	// "fundedBy":         "funded_by_ssim",
	// "publisher":        "publisher_label_tsim", // TODO: Needs URI lookup
	// "sponsor":          "sponsor_label_tsim",   // TODO: Needs URI lookup
	// "hasInstrument":    "has_instrument_ssim",
	// "sameAs":           "same_as_ssim",
	// "journalIssue":     "journal_issue_ssim",
	// "subject":          "subject_label_ssim", // TODO: Needs URI
	// "alternativeTitle": "alternative_title_tesim",

	// TODO: complex lookups
	// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
	// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
	// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.

	// date := resource.DateOfCreation
	// if date != nil {
	// 	// TODO: This may be a ETDF in the resource, but we need it to be a Solr date
	// 	doc.Set("date_created_dtsi", date)
	// }

	return doc
}
