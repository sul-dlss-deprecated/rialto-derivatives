package runtime

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// Registry is the object that holds all the external services
type Registry struct {
	Canonical repository.Repository
	Writer    derivative.Writer
}

// NewRegistry creates a new instance of the service registry
func NewRegistry(repo repository.Repository, writer derivative.Writer) *Registry {
	return &Registry{
		Canonical: repo,
		Writer:    writer,
	}
}
