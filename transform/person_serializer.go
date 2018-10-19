package transform

import (
	"encoding/json"
	"fmt"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// PersonSerializer transforms person resource types into JSON Documents
type PersonSerializer struct {
	repo repository.Repository
}

type person struct {
	Name         string    `json:"name"`
	Departments  *[]string `json:"departments"`
	Institutions *[]string `json:"institutionalAffiliations"`
	Country      *string   `json:"country"`
}

// NewPersonSerializer makes a new instance of the PersonSerializer
func NewPersonSerializer(repo repository.Repository) *PersonSerializer {
	return &PersonSerializer{repo: repo}
}

// Serialize returns the Person resource as a JSON string.
// Must include the following properties:
//
//   name (string)
//   department ([URI])
//   institution ([URI])
func (m *PersonSerializer) Serialize(resource *models.Person) string {
	p := &person{
		Name:         m.retrieveAssociatedName(resource),
		Departments:  m.retrievePositionOrganizationURIs(resource.DepartmentOrgs),
		Institutions: m.retrievePositionOrganizationURIs(resource.InstitutionOrgs),
	}

	if resource.Country != "" {
		p.Country = &resource.Country
	}

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

func (m *PersonSerializer) retrievePositionOrganizationURIs(resources []*models.PositionOrganization) *[]string {
	orgs := make([]string, len(resources))
	for n, resource := range resources {
		orgs[n] = resource.URI
	}
	return &orgs
}
