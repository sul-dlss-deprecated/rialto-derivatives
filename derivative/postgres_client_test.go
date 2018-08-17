package derivative

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestPostgresAdd(t *testing.T) {
	conf := NewPostgresConfig().WithDbname("rialto_test").WithSSL(false)
	client := NewPostgresClient(conf)
	client.RemoveAll()

	data := make(map[string][]rdf.Term)
	document, _ := rdf.NewIRI("http://purl.org/ontology/bibo/Document")
	title, _ := rdf.NewLiteral("Hello world!")

	data[models.Predicates["rdf"]["type"]] = []rdf.Term{document}
	data[models.Predicates["dct"]["title"]] = []rdf.Term{title}

	resource := models.NewResource("http://example.com/record1", data)

	err := client.addPerson(resource)

	assert.Nil(t, err)

}
