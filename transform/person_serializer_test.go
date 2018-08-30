package transform

import (
	"strings"
	"testing"

	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestSerializePersonResourceWithName(t *testing.T) {
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

	indexer := NewPersonSerializer(repository.NewService(fakeSparql))

	dept := "http://example.com/department1"
	resource := &models.Person{
		Firstname:  "Harry",
		Lastname:   "Potter",
		URI:        "http://example.com/record1",
		Department: &dept,
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"name":"Harry Potter","department":"http://example.com/department1","institutionalAffiliation":"http://example.com/institution1"}`, doc)
}
