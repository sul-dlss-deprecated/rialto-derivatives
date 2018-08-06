package transform

import (
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// CompositeIndexer delegates to subindexers to transform resources to solr Documents
type CompositeIndexer struct {
	publicationIndexer Indexer
	personIndexer      Indexer
	defaultIndexer     Indexer
}

// NewCompositeIndexer creates a new CompositeIndexer instance
func NewCompositeIndexer() *CompositeIndexer {
	return &CompositeIndexer{
		publicationIndexer: &PublicationIndexer{},
		personIndexer:      &PersonIndexer{},
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
	types := resource.ValueOf("type")
	if types != nil {
		doc.Set("type_ssi", types[0].String())
	} else {
		log.Printf("No resource types exist for %s", resource)
	}
	var indexer Indexer
	if resource.IsPublication() {
		indexer = m.publicationIndexer
	} else if resource.IsPerson() {
		indexer = m.personIndexer
	} else {
		indexer = m.defaultIndexer
	}
	return indexer.Index(resource, doc)
}
