package actions

import (
	"log"

	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// TouchAction updates a single record
type TouchAction struct {
	registry *runtime.Registry
}

// Action is a request type that this service can handle. Currently "touch" and "rebuild"
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
	log.Printf("Retrieved: %s", resourceList)
	docs := r.registry.Indexer.Map(resourceList)
	log.Printf("Writing: %s", docs)
	err = r.registry.Derivatives.Add(docs)
	return err
}

// This will take an SNS message and retrieve a resource from Neptune
func (r *TouchAction) recordToResource(msg *message.Message) (*models.Resource, error) {
	subject := msg.Entities[0]
	return r.registry.Canonical.SubjectToResource(subject)
}
