package transform

import (
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PublicationIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource *models.Publication, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Publication")
	doc.Set("title_tesi", resource.Title)
	doc.Set("created_ssim", resource.Created)

	if resource.DOI != nil {
		doc.Set("doi_ssim", *resource.DOI)
	}

	if resource.Abstract != nil {
		doc.Set("abstract_tesim", *resource.Abstract)
	}
	m.indexAuthors(resource, &doc)
	m.indexConcepts(resource, &doc)
	m.indexGrants(resource, &doc)
	doc.Set("identifiers_ssim", resource.Identifiers)

	if resource.Description != nil {
		doc.Set("description_tesim", *resource.Description)
	}

	if resource.Publisher != nil {
		doc.Set("publisher_ssim", *resource.Publisher)
	}

	if resource.CreatedYear != 0 {
		doc.Set("created_year_isim", resource.CreatedYear)
	}

	if resource.Identifiers != nil {
		doc.Set("identifiers_ssim", resource.Identifiers)
	}

	return doc
}

func (m *PublicationIndexer) indexAuthors(resource *models.Publication, doc *solr.Document) {
	var authors = []string{}
	var authorLabels = []string{}
	for _, author := range resource.Authors {
		authors = append(authors, author.URI)
		authorLabels = append(authorLabels, author.Label)
	}
	doc.Set("authors_ssim", authors)
	doc.Set("author_labels_tsim", authorLabels)
}

func (m *PublicationIndexer) indexConcepts(resource *models.Publication, doc *solr.Document) {
	var concepts = []string{}
	var conceptLabels = []string{}
	for _, concept := range resource.Concepts {
		concepts = append(concepts, concept.URI)
		conceptLabels = append(conceptLabels, concept.Label)
	}

	doc.Set("concepts_ssim", concepts)
	doc.Set("concept_labels_ssim", conceptLabels)
}

func (m *PublicationIndexer) indexGrants(resource *models.Publication, doc *solr.Document) {
	var grants = []string{}
	var grantLabels = []string{}
	for _, grant := range resource.Grants {
		grants = append(grants, grant.URI)
		grantLabels = append(grantLabels, grant.Label)
	}

	doc.Set("grants_ssim", grants)
	doc.Set("grant_labels_ssim", grantLabels)
}
