package models

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/stretchr/testify/assert"
)

func TestPublicationResource(t *testing.T) {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := NewResource("http://example.com/record1", data)

	assert.Equal(t, "Hello world!", resource.ValueOf("title")[0].String())
	assert.Equal(t, "http://purl.org/ontology/bibo/Document", resource.ValueOf("type")[0].String())
	assert.True(t, resource.IsPublication())

}

func TestPersonResource(t *testing.T) {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")
	name, _ := rdf.NewIRI("http://example.com/name1")

	data[Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[Predicates["vcard"]["hasName"]] = []rdf.Term{name}

	resource := NewResource("http://example.com/record1", data)

	assert.Equal(t, "http://example.com/name1", resource.ValueOf("name")[0].String())
	assert.Equal(t, "http://xmlns.com/foaf/0.1/Person", resource.ValueOf("type")[0].String())
	assert.True(t, resource.IsPerson())

}

func TestOrganizationResource(t *testing.T) {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")
	name, _ := rdf.NewLiteral("Cornell")

	data[Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[Predicates["skos"]["prefLabel"]] = []rdf.Term{name}
	resource := NewResource("http://example.com/record1", data)

	assert.Equal(t, "Cornell", resource.ValueOf("orgName")[0].String())
	assert.Equal(t, "http://xmlns.com/foaf/0.1/Organization", resource.ValueOf("type")[0].String())
	assert.True(t, resource.IsOrganization())
}

func TestGrantResource(t *testing.T) {
	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Grant")
	title, _ := rdf.NewLiteral("Hydra in a Box")

	data[Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[Predicates["skos"]["prefLabel"]] = []rdf.Term{title}
	resource := NewResource("http://example.com/record1", data)

	assert.Equal(t, "Hydra in a Box", resource.ValueOf("grantName")[0].String())
	assert.Equal(t, "http://vivoweb.org/ontology/core#Grant", resource.ValueOf("type")[0].String())
	assert.True(t, resource.IsGrant())
}
