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
	err := r.registry.IndexWriter.RemoveAll()
	if err != nil {
		return err
	}
	resourceList, err := r.queryResources()
	if err != nil {
		return err
	}

	err = r.registry.IndexWriter.Add(r.registry.Indexer.Map(resourceList))
	return err
}

// Return a list of resources populated by querying for everything in the triplestore
func (r *RebuildAction) queryResources() ([]models.Resource, error) {
	list, err := r.registry.Canonical.AllResources()
	if err != nil {
		return nil, err
	}
	return list, nil
}
