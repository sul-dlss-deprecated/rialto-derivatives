package actions

import (
	"github.com/sul-dlss-labs/rialto-lambda/derivative"
	"github.com/sul-dlss-labs/rialto-lambda/message"
	"github.com/sul-dlss-labs/rialto-lambda/models"
	"github.com/sul-dlss-labs/rialto-lambda/transform"
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
