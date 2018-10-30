package models

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
)

func TestNewGrantMinimalFields(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")

	name, _ := rdf.NewLiteral("Grant #1")
	pi, _ := rdf.NewIRI("http://example.com/record2")
	piName, _ := rdf.NewLiteral("Bob")
	assigned, _ := rdf.NewIRI("http://example.com/record3")
	assignedName, _ := rdf.NewLiteral("Chocolate Foundation")

	data["id"] = id
	data["name"] = name
	data["pi"] = pi
	data["pi_label"] = piName
	data["assigned"] = assigned
	data["assigned_label"] = assignedName

	resource := NewGrant(data)
	assert.IsType(t, &Grant{}, resource)
	assert.Equal(t, id.String(), resource.Subject())
	assert.Equal(t, name.String(), resource.Name)
	assert.Equal(t, pi.String(), resource.PI.URI)
	assert.Equal(t, piName.String(), resource.PI.Label)
	assert.Equal(t, assigned.String(), resource.Assigned.URI)
	assert.Equal(t, assignedName.String(), resource.Assigned.Label)
}
