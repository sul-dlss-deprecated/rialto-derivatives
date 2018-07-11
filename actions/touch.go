package actions

import (
	"github.com/sul-dlss-labs/rialto-lambda/derivative"
	"github.com/sul-dlss-labs/rialto-lambda/message"
	"github.com/sul-dlss-labs/rialto-lambda/models"
	"github.com/sul-dlss-labs/rialto-lambda/transform"
)

type TouchAction struct {
	solrClient *derivative.SolrClient
	indexer    *transform.Indexer
}

// Action is a request type that this service can handle. Currently "touch" and "delete"
type Action interface {
	Run(*message.Message) error
}

// NewTouchAction creates a Touch (update) action
func NewTouchAction(solrClient *derivative.SolrClient, indexer *transform.Indexer) Action {
	return &TouchAction{solrClient: solrClient, indexer: indexer}
}

func (r *TouchAction) Run(message *message.Message) error {
	resourceList := []models.Resource{}
	resourceList = append(resourceList, r.recordToResource(message))

	err := r.solrClient.Add(r.indexer.Map(resourceList))
	return err
}

// This will take an SNS message and retrieve a resource from Neptune
func (r *TouchAction) recordToResource(msg *message.Message) models.Resource {
	// TODO: data should come from neptune
	data := map[string]interface{}{"rdf:type": "bar", "dc:title": "whatever"}
	return models.NewResource(data)
}
