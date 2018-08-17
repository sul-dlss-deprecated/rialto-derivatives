package runtime

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// Registry is the object that holds all the external services
type Registry struct {
	Canonical   *repository.Service
	Derivatives derivative.Writer
	Database    derivative.Writer
}

// NewRegistry creates a new instance of the service registry
func NewRegistry(service *repository.Service, solr derivative.Writer, db derivative.Writer) *Registry {
	return &Registry{
		Canonical:   service,
		Derivatives: solr,
		Database:    db,
	}
}
