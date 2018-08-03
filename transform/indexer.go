package transform

import (
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
	for property, field := range mapping {
		terms := resource.ValueOf(property)
		if terms != nil {
			doc.Set(field, mapToString(terms))
		}
	}
	return doc
}
