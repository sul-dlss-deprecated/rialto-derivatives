package transform

import (
	"log"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// TypeIndexer adds the type into solr Documents
type TypeIndexer struct {
}

// SolrTypeField is the field that holds the type assertion in Solr
const SolrTypeField = "type_ssi"

// Index adds the type assertion field from the resource to the Solr Document
// This type assertion is going to drive which view is displayed, so it's important for this
// value to be more specific than "Agent", but not so specific as "Student".
// This typically would return Organization or Person
func (m *TypeIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	types := resource.ValueOf("type")
	if types != nil {
		doc.Set(SolrTypeField, m.bestType(types))
	} else {
		log.Printf("No resource types exist for %s", resource)
	}

	return doc
}

func (m *TypeIndexer) bestType(types []rdf.Term) string {
	try := types[0].String()
	if try == agent && len(types) > 1 {
		try = types[1].String()
	}
	return try
}
