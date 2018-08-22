package actions

import (
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// RebuildAction drops the repository and then rebuilds it
type RebuildAction struct {
	registry *runtime.Registry
}

// NewRebuildAction creates a Rebuild ("rebuild") action
func NewRebuildAction(registry *runtime.Registry) Action {
	return &RebuildAction{registry: registry}
}

// Run does the rebuilding
func (r *RebuildAction) Run(message *message.Message) error {
	err := r.registry.Writer.RemoveAll()
	if err != nil {
		return err
	}
	err = r.queryResources(func(resourceList []models.Resource) error {
		return r.registry.Writer.Add(resourceList)
	})
	return err
}

// Calls the function with a list of resources (in managable sized chunks) populated
// by querying for everything in the triplestore
func (r *RebuildAction) queryResources(f func([]models.Resource) error) error {
	return r.registry.Canonical.AllResources(f)
}
