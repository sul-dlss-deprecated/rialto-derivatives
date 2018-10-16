package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestSerializePersonResourceWithName(t *testing.T) {
	fakeSparql := new(MockedReader)

	indexer := NewPersonSerializer(repository.NewService(fakeSparql))

	resource := &models.Person{
		Firstname: "Harry",
		Lastname:  "Potter",
		URI:       "http://example.com/record1",
		DepartmentOrgs: []*models.PositionOrganization{&models.PositionOrganization{
			URI:   "http://example.com/department1",
			Label: "Department 1"}},
		InstitutionOrgs: []*models.PositionOrganization{&models.PositionOrganization{
			URI:   "http://example.com/institution1",
			Label: "Institution 1"}},
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"name":"Harry Potter","departments":["http://example.com/department1"],"institutionalAffiliations":["http://example.com/institution1"]}`, doc)
}
