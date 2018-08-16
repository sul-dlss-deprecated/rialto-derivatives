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

// Index adds fields from the resource to the Solr Document
// The only sting literals for person are 'email' and 'about' (It's unclear if we need to index these).
// Everything else is a URI
func (m *TypeIndexer) Index(resource models.Resource, doc solr.Document) solr.Document {
	types := resource.ValueOf("type")
	if types != nil {
		doc.Set("type_ssi", m.bestType(types))
	} else {
		log.Printf("No resource types exist for %s", resource)
	}

	return doc
}

func (m *TypeIndexer) bestType(types []rdf.Term) string {
	try := types[0].String()
	if try == "http://xmlns.com/foaf/0.1/Agent" && len(types) > 1 {
		try = types[1].String()
	}
	return try
}
