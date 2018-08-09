package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// DbTransformer delegates to subindexers to transform resources to solr Documents
type DbTransformer struct {
}

// NewDbTransformer creates a new DbTransformer instance
func NewDbTransformer() *DbTransformer {
	return &DbTransformer{}
}

// Map transforms a collection of resources into a collection of database rows
func (m *DbTransformer) Map(resources []models.Resource) []interface{} {
	docs := make([]interface{}, len(resources))
	for i, v := range resources {
		docs[i] = m.mapOne(v)
	}
	return docs
}

// mapOne sets the id and type and then delegates to the type specific indexer
func (m *DbTransformer) mapOne(resource models.Resource) interface{} {
	doc := "whatever"
	// doc.Set("id", resource.Subject)
	// types := resource.ValueOf("type")
	// if types != nil {
	// 	doc.Set("type_ssi", types[0].String())
	// } else {
	// 	log.Printf("No resource types exist for %s", resource)
	// }
	// var indexer Indexer
	// if resource.IsPublication() {
	// 	indexer = m.publicationIndexer
	// } else if resource.IsPerson() {
	// 	indexer = m.personIndexer
	// } else {
	// 	indexer = m.defaultIndexer
	// }
	// return indexer.Index(resource, doc)
	return doc
}
