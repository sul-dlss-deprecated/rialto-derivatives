package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

func TestBuildRegistry(t *testing.T) {
	registry := BuildServiceRegistry()

	assert.IsType(t, &repository.Service{}, registry.Canonical)
	assert.IsType(t, &derivative.CompositeWriter{}, registry.Writer)
}
