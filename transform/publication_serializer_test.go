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

func TestToSQLPublicationResource(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		URI:   "http://example.com/publication1",
		Title: "New developments in the management of narcolepsy",
	}

	sql, values := indexer.SQLForInsert(resource)

	assert.Equal(t, `INSERT INTO "publications" ("uri", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uri) DO UPDATE SET metadata=$2, updated_at=$4 WHERE publications.uri=$1`, sql)
	assert.Equal(t, `{"title":"New developments in the management of narcolepsy"}`, values[1])
}
