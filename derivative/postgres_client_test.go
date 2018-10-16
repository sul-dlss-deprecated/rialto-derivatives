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

func (f *MockRepository) SubjectsToResources(ids []string) ([]models.Resource, error) {
	args := f.Called(ids)
	return args.Get(0).([]models.Resource), args.Error(1)
}

func (f *MockRepository) AllResources(fun func([]models.Resource) error) error {
	return nil
}

func TestPostgresAddPerson(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	resourceType, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Person")
	student, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#Student")
	fname, _ := rdf.NewLiteral("Barbara")
	lname, _ := rdf.NewLiteral("Liskov")

	data["id"] = id
	data["type"] = resourceType
	data["subtype"] = student
	data["firstname"] = fname
	data["lastname"] = lname

	resource := models.NewResource(data)

	err := client.addPerson(resource.(*models.Person))
	assert.Nil(t, err)

	person, err := client.retrieveOnePerson("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "Barbara Liskov", "departments": [], "institutionalAffiliations": []}`, person)

}

func TestPostgresAddOrganization(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	repo := new(MockRepository)

	client := NewPostgresClient(conf, repo)
	client.RemoveAll()

	data := make(map[string]rdf.Term)
	name, _ := rdf.NewLiteral("School of Engineering")
	school, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#School")
	id, _ := rdf.NewIRI("http://example.com/record1")
	resourceType, _ := rdf.NewIRI("http://xmlns.com/foaf/0.1/Organization")

	data["id"] = id
	data["type"] = resourceType
	data["name"] = name
	data["subtype"] = school

	resource := models.NewResource(data)

	err := client.addOrganization(resource.(*models.Organization))
	assert.Nil(t, err)

	org, err := client.retrieveOneOrganization("http://example.com/record1")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"name": "School of Engineering", "type": "http://vivoweb.org/ontology/core#School"}`, org)
}
