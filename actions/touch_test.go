package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss/rialto-derivatives/message"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/runtime"
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

func (f *MockedRepository) SubjectsToResources(ids []string) ([]models.Resource, error) {
	args := f.Called(ids)
	return args.Get(0).([]models.Resource), nil
}

func TestRecordToResourceList(t *testing.T) {
	repo := new(MockedRepository)
	repo.On("SubjectsToResources", []string{"http://example.com/record2"}).
		Return([]models.Resource{&models.Person{URI: "http://example.com/record2"}})
	reg := &runtime.Registry{
		Canonical: repo,
	}
	msg := &message.Message{Entities: []string{"http://example.com/record2"}}
	action := NewTouchAction(reg)
	list, _ := action.(*TouchAction).recordToResourceList(msg)

	assert.Equal(t, "http://example.com/record2", list[0].Subject())
}
