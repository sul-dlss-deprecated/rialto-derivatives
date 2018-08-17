package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

func TestBuildRegistry(t *testing.T) {
	registry := BuildServiceRegistry()

	assert.IsType(t, &repository.Service{}, registry.Canonical)
	assert.IsType(t, &transform.CompositeIndexer{}, registry.Indexer)
	assert.IsType(t, &derivative.SolrClient{}, registry.Derivatives)
}
