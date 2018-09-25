package repository

import (
	"strings"
	"testing"

	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSPARQL struct {
	mock.Mock
}

func (m *MockSPARQL) Query(q string) (*sparql.Results, error) {
	args := m.Called(q)
	return args.Get(0).(*sparql.Results), nil
}

func TestQueryTypeForIDWithNoResults(t *testing.T) {
	mockRepo := new(MockSPARQL)
	mockRepo.On("Query", mock.Anything).
		Return(&sparql.Results{})
	reader := &SparqlReader{repo: mockRepo}
	docType, _ := reader.queryTypeForID("http://example.com/record1")

	assert.Equal(t, "", docType)
}

func TestQueryTypeForIDWithResults(t *testing.T) {
	organizationJSON := strings.NewReader(`{
    "head" : {
  "vars" : [ "type" ]
},
"results" : {
  "bindings" : [ {
    "type" : {
      "type" : "uri",
      "value" : "http://vivoweb.org/ontology/core#Department"
    }
  } ]
}
    }`)
	results, _ := sparql.ParseJSON(organizationJSON)
	mockRepo := new(MockSPARQL)
	mockRepo.On("Query", mock.Anything).
		Return(results)
	reader := &SparqlReader{repo: mockRepo}
	docType, _ := reader.queryTypeForID("http://example.com/record1")

	assert.Equal(t, "http://vivoweb.org/ontology/core#Department", docType)
}
