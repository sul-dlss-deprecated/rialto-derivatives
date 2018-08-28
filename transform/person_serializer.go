package transform

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// PersonSerializer transforms person resource types into JSON Documents
type PersonSerializer struct {
	repo repository.Repository
}

type person struct {
	Name        string  `json:"name"`
	Department  *string `json:"department"`
	Institution *string `json:"institutionalAffiliation"`
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
func (m *PersonSerializer) Serialize(resource *models.Person) string {
	p := &person{
		Name:        m.retrieveAssociatedName(resource),
		Department:  resource.Department,
		Institution: m.retrieveInstitutionURI(resource.Department)}

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// TODO: This method is copied from PersonIndexer. reduce duplication?
func (m *PersonSerializer) retrieveAssociatedName(resource *models.Person) string {
	givenName := resource.Firstname
	familyName := resource.Lastname
	return fmt.Sprintf("%v %v", givenName, familyName)
}

func (m *PersonSerializer) retrieveInstitutionURI(departmentURI *string) *string {
	if departmentURI == nil {
		return nil
	}
	uri, err := m.repo.QueryForInstitution(*departmentURI)
	if err != nil {
		panic(err)
	}
	if uri == nil {
		log.Printf("No institution URI found for %s", *departmentURI)
	}
	return uri
}
