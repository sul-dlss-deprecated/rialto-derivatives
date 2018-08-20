package derivative

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// MockRepository is a mocked object that implements the repository interface
type MockRepository struct {
	mock.Mock
}

func (f *MockRepository) SubjectToResource(id string) (models.Resource, error) {
	args := f.Called(id)
	return args.Get(0).(models.Resource), args.Error(1)
}

func (f *MockRepository) AllResources() ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (f *MockRepository) QueryForDepartment(subject string) (*string, error) {
	return nil, nil
}

func makeName(subject string, given string, family string) models.Resource {
	nameData := make(map[string][]rdf.Term)
	fname, _ := rdf.NewLiteral("Barbara")
	lname, _ := rdf.NewLiteral("Liskov")
	nameData[models.Predicates["vcard"]["given-name"]] = []rdf.Term{fname}
	nameData[models.Predicates["vcard"]["family-name"]] = []rdf.Term{lname}
	return models.NewResource(subject, nameData)
}

func TestPostgresAddPerson(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	nameID := "http://example.com/names/123"
	repo.On("SubjectToResource", nameID).
		Return(makeName(nameID, "Barbara", "Liskov"), nil)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	data := make(map[string][]rdf.Term)
	name, _ := rdf.NewIRI(nameID)

	data[models.Predicates["vcard"]["hasName"]] = []rdf.Term{name}

	resource := models.NewResource("http://example.com/record1", data)

	err := client.addPerson(resource)
	assert.Nil(t, err)

	person, err := client.retrieveOnePerson("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "Barbara Liskov", "department": null, "institutionalAffiliation": null}`, person)

}

func TestPostgresAddOrganization(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	data := make(map[string][]rdf.Term)
	name, _ := rdf.NewLiteral("School of Engineering")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")
	data[models.Predicates["skos"]["prefLabel"]] = []rdf.Term{name}
	data[models.Predicates["rdf"]["type"]] = []rdf.Term{school}

	resource := models.NewResource("http://example.com/record1", data)

	err := client.addOrganization(resource)
	assert.Nil(t, err)

	org, err := client.retrieveOneOrganization("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "School of Engineering", "type": "http://vivoweb.org/ontology/core#School"}`, org)
}
