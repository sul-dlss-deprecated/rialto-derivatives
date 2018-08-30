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

func createPersonResource() *models.Person {
	data := make(map[string]rdf.Term)
	fname, _ := rdf.NewLiteral("Christina")
	lname, _ := rdf.NewLiteral("Harlow")
	id, _ := rdf.NewIRI("http://example.com/record1")
	person, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")

	data["id"] = id
	data["type"] = person
	data["firstname"] = fname
	data["lastname"] = lname

	return models.NewResource(data).(*models.Person)
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

func (f *MockedReader) QueryEverything(fun func(*sparql.Results) error) error {
	return nil
}

func (f *MockedReader) QueryByID(id string) (*sparql.Results, error) {
	args := f.Called(id)
	return args.Get(0).(*sparql.Results), args.Error(1)
}

func (f *MockedReader) QueryByIDAndPredicate(id string, predicate string) (*sparql.Results, error) {
	args := f.Called(id)
	return args.Get(0).(*sparql.Results), args.Error(1)
}

func TestPersonResourceWithName(t *testing.T) {
	fakeSparql := new(MockedReader)

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

	resource := &models.Person{Firstname: "Harry", Lastname: "Potter", URI: "http://example.com/record1"}

	in := make(solr.Document)
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Harry Potter", doc.Get("name_ssim"))
}

func TestPersonWithoutDepartment(t *testing.T) {
	fakeSparql := new(MockedReader)

	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := &models.Person{Firstname: "Hermione", Lastname: "Granger", URI: "http://example.com/record1"}
	in := make(solr.Document)
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Hermione Granger", doc.Get("name_ssim"))
}
