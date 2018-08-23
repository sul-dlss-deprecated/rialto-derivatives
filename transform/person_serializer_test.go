package transform

import (
	"strings"
	"testing"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestSerializePersonResourceWithName(t *testing.T) {
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

	indexer := NewPersonSerializer(repository.NewService(fakeSparql))

	resource := new(MockResource)
	nameURI, _ := rdf.NewIRI("http://example.com/name1")
	resource.On("ValueOf", "name").
		Return([]rdf.Term{nameURI})

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"name":"Harry Potter","department":"http://example.com/department1","institutionalAffiliation":"http://example.com/institution1"}`, doc)
}
