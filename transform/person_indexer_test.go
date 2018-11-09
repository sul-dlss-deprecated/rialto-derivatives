package transform

import (
	"testing"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

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

func (f *MockedReader) QueryByIDs(ids []string) ([]*sparql.Results, error) {
	args := f.Called(ids)
	return args.Get(0).([]*sparql.Results), args.Error(1)
}

func (f *MockedReader) GetOrganizationInfo(id *string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetPositionOrganizationInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetAuthorInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetConceptInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetGrantInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetCountriesInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetPersonSubtypesInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func (f *MockedReader) GetIdentifierInfo(id string) (*sparql.Results, error) {
	return nil, nil
}

func TestPersonResourceWithName(t *testing.T) {
	fakeSparql := new(MockedReader)

	// TODO Fix
	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := &models.Person{Firstname: "Harry", Lastname: "Potter", URI: "http://example.com/record1"}

	in := make(solr.Document)
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Harry Potter", doc.Get("name_tsim"))
}

func TestPersonWithoutDepartment(t *testing.T) {
	fakeSparql := new(MockedReader)

	indexer := &PersonIndexer{
		Canonical: repository.NewService(fakeSparql),
	}

	resource := &models.Person{Firstname: "Hermione", Lastname: "Granger", URI: "http://example.com/record1"}
	in := make(solr.Document)
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Hermione Granger", doc.Get("name_tsim"))
}
