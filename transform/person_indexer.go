package transform

import (
	"fmt"
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

// PersonIndexer transforms person resource types into solr Documents
type PersonIndexer struct {
	Canonical repository.Repository
}

// NewPersonIndexer creates a new instance of the Person indexer
func NewPersonIndexer(repository repository.Repository) Indexer {
	return &PersonIndexer{Canonical: repository}
}

// Index adds fields from the resource to the Solr Document
// The only sting literals for person are 'email' and 'about' (It's unclear if we need to index these).
// Everything else is a URI
func (m *PersonIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	// 1. Get the associated name resource
	doc.Set("name_ssim", m.retrieveAssociatedName(resource))

	// 2. department
	doc.Set("department_ssim", m.retrieveDepartmentURI(resource))
	// TODO 3. institution

	return doc
}

func (m *PersonIndexer) retrieveAssociatedName(resource models.Resource) string {
	nameURI := resource.ValueOf("name")
	if len(nameURI) == 0 {
		log.Printf("No name URI found for %s", resource.Subject())
		return ""
	}

	nameResource, err := m.Canonical.SubjectToResource(nameURI[0].String())
	if err != nil {
		panic(err)
	}
	givenName := nameResource.ValueOf("given-name")
	familyName := nameResource.ValueOf("family-name")

	if len(givenName) == 0 || len(familyName) == 0 {
		return ""
	}
	return fmt.Sprintf("%v %v", givenName[0], familyName[0])
}

func (m *PersonIndexer) retrieveDepartmentURI(resource models.Resource) *string {
	uri, err := m.Canonical.QueryForDepartment(resource.Subject())
	if err != nil {
		panic(err)
	}
	if uri == nil {
		log.Printf("No department URI found for %s", resource.Subject())
	}
	return uri
}
