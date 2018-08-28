package actions

import (
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

// Run fetches a record from Neptune and updates the configured writer
func (r *TouchAction) Run(message *message.Message) error {
	resourceList, err := r.recordToResourceList(message)
	if err != nil {
		return err
	}
	err = r.registry.Writer.Add(resourceList)
	return err
}

// This will take an SNS message and retrieve a resource from Neptune
func (r *TouchAction) recordToResourceList(msg *message.Message) ([]models.Resource, error) {
	subject := msg.Entities[0]
	resource, err := r.registry.Canonical.SubjectToResource(subject)
	if err != nil {
		return nil, err
	}
	resourceList := []models.Resource{}
	resourceList = append(resourceList, resource)
	return resourceList, err
}
