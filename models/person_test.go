package models

import (
	"strings"
	"testing"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
)

func TestNewPersonMinimalFields(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")

	data["id"] = id

	resource := NewPerson(data)
	assert.IsType(t, &Person{}, resource)
	assert.Equal(t, id.String(), resource.Subject())
}

func TestNewPersonAllFields(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	faculty, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#FacultyMember")

	fname, _ := rdf.NewLiteral("Justin")
	lname, _ := rdf.NewLiteral("Coyne")
	data["id"] = id
	data["subtype"] = faculty
	data["lastname"] = lname
	data["firstname"] = fname

	resource := NewPerson(data)
	assert.IsType(t, &Person{}, resource)
	assert.Equal(t, faculty.String(), resource.Subtype)
	assert.Equal(t, resource.Firstname, fname.String())
	assert.Equal(t, resource.Lastname, lname.String())

}

func TestSetPositionOrganizationInfo(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	faculty, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#FacultyMember")

	data["id"] = id
	data["subtype"] = faculty

	resource := NewPerson(data)

	organizationJSON := strings.NewReader(`{
    "head" : {
  "vars" : [ "org", "type", "name" ]
},
"results" : {
  "bindings" : [ {
    "org" : {
      "type" : "uri",
      "value" : "http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering/electrical-engineering"
    },
    "name" : {
      "type" : "literal",
      "value" : "Electrical Engineering"
    },
    "type" : {
      "type" : "uri",
      "value" : "http://vivoweb.org/ontology/core#Department"
    }
  }, {
	 	"org" : {
	 		"type" : "uri",
	 		"value" : "http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering/nuclear-engineering"
	 	},
	 	"name" : {
	 		"type" : "literal",
	 		"value" : "Nuclear Engineering"
	 	},
	 	"type" : {
	 		"type" : "uri",
	 		"value" : "http://vivoweb.org/ontology/core#Department"
	 	}
	 }, {
    "org" : {
      "type" : "uri",
      "value" : "http://sul.stanford.edu/rialto/agents/orgs/stanford"
    },
    "name" : {
      "type" : "literal",
      "value" : "Stanford University"
    },
    "type" : {
      "type" : "uri",
      "value" : "http://vivoweb.org/ontology/core#University"
    }
  }, {
    "org" : {
      "type" : "uri",
      "value" : "http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering"
    },
    "name" : {
      "type" : "literal",
      "value" : "School of Engineering"
    },
    "type" : {
      "type" : "uri",
      "value" : "http://vivoweb.org/ontology/core#School"
    }
  }, {
    "org" : {
      "type" : "uri",
      "value" : "http://sul.stanford.edu/rialto/agents/orgs/leland-junior-u"
    },
    "name" : {
      "type" : "literal",
      "value" : "Leland Junior University"
    }
  } ]
}
	  }`)
	results, _ := sparql.ParseJSON(organizationJSON)
	resource.SetPositionOrganizationInfo(results)
	assert.IsType(t, &Person{}, resource)
	assert.Equal(t, []*Labeled{
		&Labeled{"http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering/electrical-engineering", "Electrical Engineering"},
		&Labeled{"http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering/nuclear-engineering", "Nuclear Engineering"}},
		resource.DepartmentOrgs)
	assert.Equal(t, []*Labeled{
		&Labeled{"http://sul.stanford.edu/rialto/agents/orgs/school-of-engineering", "School of Engineering"}},
		resource.SchoolOrgs)
	assert.Equal(t, []*Labeled{
		&Labeled{"http://sul.stanford.edu/rialto/agents/orgs/stanford", "Stanford University"},
		&Labeled{"http://sul.stanford.edu/rialto/agents/orgs/leland-junior-u", "Leland Junior University"}},
		resource.InstitutionOrgs)
}

func TestSetCountriesInfo(t *testing.T) {
	data := make(map[string]rdf.Term)
	id, _ := rdf.NewIRI("http://example.com/record1")
	faculty, _ := rdf.NewIRI("http://vivoweb.org/ontology/core#FacultyMember")

	data["id"] = id
	data["subtype"] = faculty

	resource := NewPerson(data)

	countriesJSON := strings.NewReader(`{
    "head" : {
  		"vars" : [ "country" ]
		},
		"results" : {
  		"bindings" : [ {
    		"country" : {
      		"type" : "uri",
      		"value" : "http://sws.geonames.org/1814991/"
    		},
				"label" : {
      		"type" : "literal",
      		"value" : "United States"
    		}
  		} ]
		}
	}`)
	results, _ := sparql.ParseJSON(countriesJSON)
	resource.SetCountriesInfo(results)
	assert.IsType(t, &Person{}, resource)
	assert.Equal(t, []*Labeled{
		&Labeled{"http://sws.geonames.org/1814991/", "United States"}},
		resource.Countries)
}
