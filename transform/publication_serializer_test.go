package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestSerializePublicationResource(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		URI:   "http://example.com/publication1",
		Title: "New developments in the management of narcolepsy",
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"title":"New developments in the management of narcolepsy"}`, doc)
}
