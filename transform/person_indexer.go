package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

// PersonIndexer transforms person resource types into solr Documents
type PersonIndexer struct {
	Canonical repository.Repository
}

// NewPersonIndexer creates a new instance of the Person indexer
func NewPersonIndexer(repository repository.Repository) *PersonIndexer {
	return &PersonIndexer{Canonical: repository}
}

// Index adds fields from the resource to the Solr Document
// The only sting literals for person are 'email' and 'about' (It's unclear if we need to index these).
// Everything else is a URI
func (m *PersonIndexer) Index(resource *models.Person, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Person")

	// 1. Get the associated name resource
	doc.Set("name_tsim", resource.Name())

	// 2. department
	doc.Set("department_label_ssim", m.retrieveLabels(resource.DepartmentOrgs))

	// 3. school
	doc.Set("school_label_ssim", m.retrieveLabels(resource.SchoolOrgs))

	// 4. institution
	doc.Set("institution_label_ssim", m.retrieveLabels(resource.InstitutionOrgs))

	// 5. countries
	doc.Set("countries_label_ssim", m.retrieveLabels(resource.Countries))

	// 6. subtype
	doc.Set("person_subtype_ssim", m.retrieveLabels(resource.Subtypes))

	return doc
}

func (m *PersonIndexer) retrieveLabels(resources []*models.Labeled) *[]string {
	labels := make([]string, len(resources))
	for n, resource := range resources {
		labels[n] = resource.Label
	}
	return &labels
}
