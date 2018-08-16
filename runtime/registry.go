package runtime

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// Registry is the object that holds all the external services
type Registry struct {
	DbWriter      models.ModelWriter
	DbTransformer *transform.DbTransformer
	IndexWriter   derivative.IndexWriter
	Indexer       *transform.CompositeIndexer
	Canonical     *repository.Service
}

// NewRegistry creates a new instance of the service registry
func NewRegistry(dbWriter models.ModelWriter, dbTransformer *transform.DbTransformer, indexClient derivative.IndexWriter, indexer *transform.CompositeIndexer, sparql repository.Reader) *Registry {
	return &Registry{
		DbWriter:      dbWriter,
		DbTransformer: dbTransformer,
		IndexWriter:   indexClient,
		Indexer:       indexer,
		Canonical:     repository.NewService(sparql),
	}
}
