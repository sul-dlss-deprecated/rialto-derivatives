package actions

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

type RebuildAction struct {
	solrClient *derivative.SolrClient
	indexer    *transform.Indexer
}

// NewRebuildAction creates a Rebuild (delete) action
func NewRebuildAction(solrClient *derivative.SolrClient, indexer *transform.Indexer) Action {
	return &RebuildAction{solrClient: solrClient, indexer: indexer}
}

func (r *RebuildAction) Run(message *message.Message) error {
	err := r.solrClient.RemoveAll()
	if err != nil {
		return err
	}
	resourceList, err := r.queryResources()
	if err != nil {
		return err
	}

	err = r.solrClient.Add(r.indexer.Map(resourceList))
	return err
}

func (r *RebuildAction) queryResources() ([]models.Resource, error) {
	list := []models.Resource{}
	// TODO: Populate this list from Neptune
	return list, nil
}
