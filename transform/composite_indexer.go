package transform

import (
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// CompositeIndexer delegates to subindexers to transform resources to solr Documents
type CompositeIndexer struct {
	publicationIndexer Indexer
	defaultIndexer     Indexer
}

// NewCompositeIndexer creates a new CompositeIndexer instance
func NewCompositeIndexer() *CompositeIndexer {
	return &CompositeIndexer{
		publicationIndexer: &PublicationIndexer{},
		defaultIndexer:     &DefaultIndexer{},
	}
}

// Map transforms a collection of resources into a collection of Solr Documents
func (m *CompositeIndexer) Map(resources []models.Resource) []solr.Document {
	docs := make([]solr.Document, len(resources))
	for i, v := range resources {
		docs[i] = m.mapOne(v)
	}
	return docs
}

// mapOne sets the id and type and then delegates to the type specific indexer
func (m *CompositeIndexer) mapOne(resource models.Resource) solr.Document {
	doc := make(solr.Document)
	doc.Set("id", resource.Subject)
	if resource.ResourceTypes() != nil {
		doc.Set("type_ssi", resource.ResourceTypes()[0].String())
	} else {
		log.Printf("No resource types exist for %s", resource)
	}

	if resource.IsPublication() {
		return m.publicationIndexer.Index(resource, doc)
	}
	return m.defaultIndexer.Index(resource, doc)
}
