package actions

import (
	"log"

	"github.com/sul-dlss/rialto-derivatives/message"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/runtime"
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

// Run fetches a record from Neptune and updates the configured writer
func (r *TouchAction) Run(message *message.Message) error {
	log.Printf("Received a message with %v entities", len(message.Entities))
	resourceList, err := r.recordToResourceList(message)
	if err != nil {
		return err
	}
	err = r.registry.Writer.Add(resourceList)
	return err
}

// This will take an SNS message and retrieve a resource from Neptune
func (r *TouchAction) recordToResourceList(msg *message.Message) ([]models.Resource, error) {
	resources, err := r.registry.Canonical.SubjectsToResources(msg.Entities)
	if err != nil {
		return nil, err
	}
	resourceList := []models.Resource{}
	resourceList = append(resourceList, resources...)
	return resourceList, err
}
