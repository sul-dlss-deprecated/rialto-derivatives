package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestOrganizationSerializer(t *testing.T) {
	data := make(map[string][]rdf.Term)
	name, _ := rdf.NewLiteral("School of Engineering")
	agent, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Agent")
	organization, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")

	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{name}
	data[models.Predicates["rdf"]["type"]] = []rdf.Term{agent, organization, school}

	resource := models.NewResource("http://example.com/record1", data)

	org := (&OrganizationSerializer{}).Serialize(resource)
	assert.Equal(t, `{"name": "School of Engineering", "type": "http://vivoweb.org/ontology/core#School"}`, org)
}
