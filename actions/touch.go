package actions

import (
	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// TouchAction updates a single record
type TouchAction struct {
	registry *runtime.Registry
}

// Action is a request type that this service can handle. Currently "touch" and "delete"
type Action interface {
	Run(*message.Message) error
}

// NewTouchAction creates a Touch (update) action
func NewTouchAction(registry *runtime.Registry) Action {
	return &TouchAction{registry: registry}
}

// Run updates the record
func (r *TouchAction) Run(message *message.Message) error {
	resourceList := []models.Resource{}
	resource, err := r.recordToResource(message)
	if err != nil {
		return err
	}

	resourceList = append(resourceList, *resource)

	err = r.registry.Derivatives.Add(r.registry.Indexer.Map(resourceList))
	return err
}

// This will take an SNS message and retrieve a resource from Neptune
func (r *TouchAction) recordToResource(msg *message.Message) (*models.Resource, error) {
	subject := msg.Entities[0]
	response, err := r.registry.Canonical.QueryByID(subject)
	if err != nil {
		return nil, err
	}
	data := map[string][]rdf.Term{}
	for _, triple := range response.Solutions() {
		predicate := triple["p"].String()
		object := triple["o"]

		if data[predicate] == nil {
			// First time we encounter a predicate
			data[predicate] = []rdf.Term{object}
		} else {
			// subsequent encounters
			data[predicate] = append(data[predicate], object)
		}
	}
	resource := models.NewResource(subject, data)
	return &resource, nil
}
