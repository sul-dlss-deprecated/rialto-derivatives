package transform

import (
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// Indexer is the interface for objects that transform resources into Solr Documents
type Indexer interface {
	Index(models.Resource, solr.Document) solr.Document
}
