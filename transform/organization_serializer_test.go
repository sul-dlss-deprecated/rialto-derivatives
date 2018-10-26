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
	assert.Equal(t, `{"type": "http://vivoweb.org/ontology/core#School"}`, org)
}

func TestToSQLOrganizationResource(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	name, _ := rdf.NewLiteral("School of Engineering")
	organization, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")

	data["id"] = id
	data["name"] = name
	data["type"] = organization
	data["subtype"] = school

	resource := models.NewResource(data).(*models.Organization)

	sql, values := (&OrganizationSerializer{}).SQLForInsert(resource)

	assert.Equal(t, `INSERT INTO "organizations" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE organizations.uri=$1`, sql)
	assert.Equal(t, "School of Engineering", values[1])
	assert.Equal(t, `{"type": "http://vivoweb.org/ontology/core#School"}`, values[2])
}
