package transform

import (
	"fmt"
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// PersonSerializer transforms person resource types into JSON Documents
type PersonSerializer struct {
	repo repository.Repository
}

// NewPersonSerializer makes a new instance of the PersonSerializer
func NewPersonSerializer(repo repository.Repository) *PersonSerializer {
	return &PersonSerializer{repo: repo}
}

// Serialize returns the Person resource as a JSON string.
// Must include the following properties:
//
//   name (string)
//   department (URI)
//   institutionalAffiliation (URI)
func (m *PersonSerializer) Serialize(resource models.Resource) string {
	return fmt.Sprintf(`{"name": "%s"}`, m.retrieveAssociatedName(resource))
}

// TODO: This method is copied from PersonIndexer.  In order to be more efficient,
// we should lookup names before passing to the postgres/solr writers.
func (m *PersonSerializer) retrieveAssociatedName(resource models.Resource) string {
	nameURI := resource.ValueOf("name")
	if len(nameURI) == 0 {
		log.Printf("No name URI found for %s", resource.Subject())
		return ""
	}

	nameResource, err := m.repo.SubjectToResource(nameURI[0].String())

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
