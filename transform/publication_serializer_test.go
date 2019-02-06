package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss/rialto-derivatives/models"
)

func TestSerializePublicationResource(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		URI:               "http://example.com/publication1",
		Title:             "New developments in the management of narcolepsy",
		CreatedYear:       2018,
		HasStanfordAuthor: true,
		Concepts: []*models.Concept{&models.Concept{
			URI:   "http://example.com/concept1",
			Label: "Research Area 1"}},
	}

	doc := indexer.Serialize(resource)

	assert.Equal(t, `{"title":"New developments in the management of narcolepsy","created_year":2018,"concepts":["http://example.com/concept1"]}`, doc)
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
	assert.Equal(t, `{"title":"New developments in the management of narcolepsy","created_year":null,"concepts":[]}`, values[1])
}

func TestShouldAddPublicationResource(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		CreatedYear:       2018,
		HasStanfordAuthor: true,
	}

	assert.True(t, indexer.ShouldAdd(resource))
}

func TestShouldNotAddPublicationResourceCreated(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		CreatedYear:       1776,
		HasStanfordAuthor: true,
	}

	assert.False(t, indexer.ShouldAdd(resource))
}

func TestShouldNotAddPublicationResourceStanfordAuthor(t *testing.T) {
	indexer := NewPublicationSerializer()

	resource := &models.Publication{
		CreatedYear:       2020,
		HasStanfordAuthor: false,
	}

	assert.False(t, indexer.ShouldAdd(resource))
}
