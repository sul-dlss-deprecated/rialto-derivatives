package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PublicationIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	mapping := map[string]string{
		"Titles":     "title_ssi",
		"Abstract":   "abstract_tesim",
		"Cites":      "cites_ssim",
		"DOI":        "doi_ssim",
		"Identifier": "identifier_ssim",
		"Link":       "link_ssim",
	}

	doc = indexMapping(resource, doc, mapping)

	// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
	// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
	// date of creation 	dct:created 	DateTime string, EDTF 	[1,1] 	Used to describe the creation date of a resource.
	// description 	vivo:description 	string-literal 	[0,n] 	Description of the resource.
	// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.
	// funded by 	vivo:hasFundingVehicle 	Grant URI 	[0,n] 	Grant (or contract) providing funding for the publication.
	// has instrument 	gcis:hasInstrument 	gcis:Instrument URI 	[0,n] 	A type of tool or device used for a particular task, especially for scientific work, as presented in the publication (specifically for Datasets).
	// journal issue 	dcterms:hasPart 	Document URI (Article) 	[0,n] 	Journal is another entity with issue number, label / title, possibly isPartOf URI for the Journal title overall.
	// publisher 	vivo:publisher 	URI for foaf:Organization 	[0,n] 	Publisher of the resource.
	// sameAs 	owl:sameAs 	URI 	[0,n] 	Other resources (identified via URIs) that are the same as this resource.
	// sponsor 	vivo:informationResourceSupportedBy 	Agent URI 	[0,n] 	Institution supporting the publication.
	// subject 	dcterms:subject 	Topic / Concept URI 	[0,n] 	Topic or concept the resource is about.
	// title 	dcterms:title 	string-literal 	[1,1] 	Title for the resource.
	// alternate title

	return doc
}
