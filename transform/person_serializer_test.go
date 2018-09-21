package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestSerializePersonResourceWithName(t *testing.T) {
	fakeSparql := new(MockedReader)

	// TODO Fix
	indexer := NewPersonSerializer(repository.NewService(fakeSparql))

	dept := "http://example.com/department1"
	inst := "http://example.com/institution1"
	resource := &models.Person{
		Firstname:      "Harry",
		Lastname:       "Potter",
		URI:            "http://example.com/record1",
		DepartmentURI:  &dept,
		InstitutionURI: &inst,
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"name":"Harry Potter","department":"http://example.com/department1","institutionalAffiliation":"http://example.com/institution1"}`, doc)
}
