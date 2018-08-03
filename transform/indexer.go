package transform

import (
	"reflect"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// Indexer is the interface for objects that transform resources into Solr Documents
type Indexer interface {
	Index(models.Resource, solr.Document) solr.Document
}

func mapToString(terms []rdf.Term) []string {
	out := make([]string, len(terms))
	for i, v := range terms {
		out[i] = v.String()
	}
	return out
}

func indexMapping(resource models.Resource, doc solr.Document, mapping map[string]string) solr.Document {
	for method, field := range mapping {
		values := reflect.ValueOf(&resource).MethodByName(method).Call([]reflect.Value{})
		if values != nil {
			terms := values[0].Interface().([]rdf.Term)
			doc.Set(field, mapToString(terms))
		}
	}
	return doc
}
