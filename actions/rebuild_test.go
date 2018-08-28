package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// MockedWriter is a mocked object that implements the Writer interface
type MockedWriter struct {
	mock.Mock
}

func (f *MockedWriter) Add(docs []models.Resource) error {
	return nil
}

func (f *MockedWriter) RemoveAll() error {
	return nil
}

// MockedRepository is a mocked object that implements the Repository interface
type MockedRepository struct {
	mock.Mock
}

func (f *MockedRepository) AllResources(fun func([]models.Resource) error) error {
	return nil
}

func (f *MockedRepository) SubjectToResource(id string) (models.Resource, error) {
	args := f.Called(id)
	return args.Get(0).(models.Resource), nil
}

func (f *MockedRepository) QueryForInstitution(id string) (*string, error) {
	return nil, nil
}

func TestRebuildRepository(t *testing.T) {
	msg := &message.Message{}
	fakeWriter := new(MockedWriter)
	fakeRepo := new(MockedRepository)
	reg := &runtime.Registry{
		Writer:    fakeWriter,
		Canonical: fakeRepo,
	}
	action := NewRebuildAction(reg)
	err := action.Run(msg)

	assert.Nil(t, err)
}
