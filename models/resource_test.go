package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	data := make(map[string]interface{})
	data["rdf:type"] = "http://purl.org/ontology/bibo/Document"
	data["dc:title"] = "Hello world!"

	resource := NewResource(data)

	assert.Equal(t, "Hello world!", resource.Title())
	assert.Equal(t, "http://purl.org/ontology/bibo/Document", resource.ResourceType())
	assert.True(t, resource.IsPublication())

}
