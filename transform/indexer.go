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
