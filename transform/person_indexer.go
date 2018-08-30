package transform

import (
	"fmt"

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
	doc.Set("name_ssim", m.retrieveAssociatedName(resource))

	// 2. department
	doc.Set("department_ssim", resource.Department)
	// TODO 3. institution

	return doc
}

func (m *PersonIndexer) retrieveAssociatedName(resource *models.Person) string {
	givenName := resource.Firstname
	familyName := resource.Lastname
	return fmt.Sprintf("%v %v", givenName, familyName)
}
