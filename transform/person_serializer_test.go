package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestSerializePersonResource(t *testing.T) {
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
		Countries: []string{"http://sws.geonames.org/1814991/"},
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"departments":["http://example.com/department1"],"institutionalAffiliations":["http://example.com/institution1"],"countries":["http://sws.geonames.org/1814991/"]}`, doc)
}

func TestToSQLPersonResource(t *testing.T) {
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
		Countries: []string{"http://sws.geonames.org/1814991/"},
	}

	sql, values := indexer.SQLForInsert(resource)

	assert.Equal(t, `INSERT INTO "people" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE people.uri=$1`, sql)
	assert.Equal(t, "Harry Potter", values[1])
	assert.Equal(t, `{"departments":["http://example.com/department1"],"institutionalAffiliations":["http://example.com/institution1"],"countries":["http://sws.geonames.org/1814991/"]}`, values[2])
}
