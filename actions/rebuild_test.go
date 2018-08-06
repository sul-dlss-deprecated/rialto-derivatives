package actions

import (
	"testing"

	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
	"github.com/vanng822/go-solr/solr"
)

// MockedWriter is a mocked object that implements the Writer interface
type MockedWriter struct {
	mock.Mock
}

func (f *MockedWriter) Add(docs []solr.Document) error {
	return nil
}

func (f *MockedWriter) RemoveAll() error {
	return nil
}

// MockedReader is a mocked object that implements the Reader interface
type MockedReader struct {
	mock.Mock
}

func (f *MockedReader) QueryEverything() (*sparql.Results, error) {
	return &sparql.Results{}, nil
}

func (f *MockedReader) QueryByID(id string) (*sparql.Results, error) {
	return &sparql.Results{}, nil
}

func TestRebuildRepository(t *testing.T) {
	msg := &message.Message{}
	fakeSolr := new(MockedWriter)
	fakeSparql := new(MockedReader)
	reg := &runtime.Registry{
		Derivatives: fakeSolr,
		Canonical:   repository.NewService(fakeSparql),
	}
	action := NewRebuildAction(reg)
	err := action.Run(msg)

	assert.Nil(t, err)
}
