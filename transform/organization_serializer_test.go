package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestOrganizationSerializer(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	name, _ := rdf.NewLiteral("School of Engineering")
	organization, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")

	data["id"] = id
	data["name"] = name
	data["type"] = organization
	data["subtype"] = school

	resource := models.NewResource(data)

	org := (&OrganizationSerializer{}).Serialize(resource.(*models.Organization))
	assert.Equal(t, `{"name": "School of Engineering", "type": "http://vivoweb.org/ontology/core#School"}`, org)
}
