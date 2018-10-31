package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestSerializeConceptResource(t *testing.T) {
	indexer := NewConceptSerializer()

	resource := &models.Concept{
		Label: "Philosophy",
		URI:   "http://example.com/concept1",
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{}`, doc)
}

func TestToSQLConceptResource(t *testing.T) {
	indexer := NewConceptSerializer()

	resource := &models.Concept{
		Label: "Philosophy",
		URI:   "http://example.com/concept1",
	}

	sql, values := indexer.SQLForInsert(resource)

	assert.Equal(t, `INSERT INTO "concepts" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE concepts.uri=$1`, sql)
	assert.Equal(t, "http://example.com/concept1", values[0])
	assert.Equal(t, "Philosophy", values[1])
	assert.Equal(t, `{}`, values[2])
}
