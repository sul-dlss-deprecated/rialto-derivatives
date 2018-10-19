package transform

import (
	"encoding/json"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// PublicationSerializer transforms publication resource types into JSON Documents
type PublicationSerializer struct {
}

type publication struct {
	Title string `json:"title"`
}

// NewPublicationSerializer makes a new instance of the PublicationSerializer
func NewPublicationSerializer() *PublicationSerializer {
	return &PublicationSerializer{}
}

// Serialize returns the Publication resource as a JSON string.
func (m *PublicationSerializer) Serialize(resource *models.Publication) string {
	p := &publication{
		Title: resource.Title,
	}

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}
