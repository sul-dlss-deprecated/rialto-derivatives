package runtime

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// Registry is the object that holds all the external services
type Registry struct {
	Derivatives derivative.Writer
	Indexer     *transform.CompositeIndexer
	Canonical   *repository.Service
}

// NewRegistry creates a new instance of the service registry
func NewRegistry(solr derivative.Writer, indexer *transform.CompositeIndexer, service *repository.Service) *Registry {
	return &Registry{
		Derivatives: solr,
		Indexer:     indexer,
		Canonical:   service,
	}
}
