package transform

import (
	"strings"
	"testing"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

func createPersonResource() models.Resource {
	data := make(map[string][]rdf.Term)
	name, _ := rdf.NewIRI("http://example.com/name1")
	data[models.Predicates["vcard"]["hasName"]] = []rdf.Term{name}

	return models.NewResource("http://example.com/record1", data)
}

type MockResource struct {
	mock.Mock
}

func (f *MockResource) IsPerson() bool {
	return false
}

func (f *MockResource) IsPublication() bool {
	return false
}

func (f *MockResource) IsGrant() bool {
	return false
}

func (f *MockResource) IsProject() bool {
	return false
}

func (f *MockResource) IsOrganization() bool {
	return false
}

func (f *MockResource) IsConcept() bool {
	return false
}

func (f *MockResource) Subject() string {
	return "http://example.com/record1"
}

func (f *MockResource) ValueOf(key string) []rdf.Term {
	args := f.Called(key)
	return args.Get(0).([]rdf.Term)
}

// MockedReader is a mocked object that implements the Reader interface
type MockedReader struct {
	mock.Mock
}

func (f *MockedReader) QueryEverything() (*sparql.Results, error) {
	return &sparql.Results{}, nil
}

func (f *MockedReader) QueryByID(id string) (*sparql.Results, error) {
	args := f.Called(id)
	return args.Get(0).(*sparql.Results), args.Error(1)
}

func (f *MockedReader) QueryByIDAndPredicate(id string, predicate string) (*sparql.Results, error) {
	args := f.Called(id)
	return args.Get(0).(*sparql.Results), args.Error(1)
}

func (f *MockedReader) QueryThroughNode(id string, localPredicate string, localType string, remotePredicate string) (*sparql.Results, error) {
	args := f.Called(id)
	return args.Get(0).(*sparql.Results), nil
}

func TestPersonResourceWithName(t *testing.T) {
	fakeSparql := new(MockedReader)
	json := strings.NewReader(`{
    "head": { "vars": [ "s" , "p", "o" ] } ,
    "results": {
      "bindings": [
        {
          "s": { "type": "uri" , "value": "http://example.com/name1" } ,
          "p": { "type": "uri" , "value": "http://www.w3.org/2006/vcard/ns#given-name" } ,
          "o": { "type": "literal" , "value": "Harry" }
        },
        {
          "s": { "type": "uri" , "value": "http://example.com/name1" } ,
          "p": { "type": "uri" , "value": "http://www.w3.org/2006/vcard/ns#family-name" } ,
          "o": { "type": "literal" , "value": "Potter" }
        }
      ]
    }
  }`)
	fakeSparql.On("QueryByID", "http://example.com/name1").
		Return(sparql.ParseJSON(json))

	departmentJSON := strings.NewReader(`{
    "head": { "vars": [ "d" ] } ,
    "results": {
      "bindings": [
        {
          "d": { "type": "uri" , "value": "http://example.com/department1" }
        }
      ]
    }
  }`)
	fakeSparql.On("QueryThroughNode", "http://example.com/record1").
		Return(sparql.ParseJSON(departmentJSON))

	institutionJSON := strings.NewReader(`{
	    "head": { "vars": [ "o" ] } ,
	    "results": {
	      "bindings": [
	        {
	          "o": { "type": "uri" , "value": "http://example.com/institution1" }
	        }
	      ]
	    }
	  }`)
	fakeSparql.On("QueryByIDAndPredicate", "http://example.com/department1").
		Return(sparql.ParseJSON(institutionJSON))

	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := new(MockResource)
	nameURI, _ := rdf.NewIRI("http://example.com/name1")
	resource.On("ValueOf", "name").
		Return([]rdf.Term{nameURI})

	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Harry Potter", doc.Get("name_ssim"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}

func TestPersonWithoutNameUriOrDepartment(t *testing.T) {
	fakeSparql := new(MockedReader)
	json := strings.NewReader(`{}`)
	fakeSparql.On("QueryByID", "http://example.com/name1").
		Return(sparql.ParseJSON(json))

	departmentJSON := strings.NewReader(`{}`)
	fakeSparql.On("QueryThroughNode", "http://example.com/record1").
		Return(sparql.ParseJSON(departmentJSON))

	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := new(MockResource)
	resource.On("ValueOf", "name").
		Return([]rdf.Term{})
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "", doc.Get("name_ssim"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}

func TestPersonWhenNameIsNotFound(t *testing.T) {
	fakeSparql := new(MockedReader)
	json := strings.NewReader(`{}`)
	fakeSparql.On("QueryByID", "http://example.com/name1").
		Return(sparql.ParseJSON(json))

	departmentJSON := strings.NewReader(`{}`)
	fakeSparql.On("QueryThroughNode", "http://example.com/record1").
		Return(sparql.ParseJSON(departmentJSON))

	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := createPersonResource()
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "", doc.Get("name_ssim"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
}
