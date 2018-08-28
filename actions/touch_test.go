package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

func TestRecordToResourceList(t *testing.T) {
	repo := new(MockedRepository)
	repo.On("SubjectToResource", "http://example.com/record2").
		Return(&models.Person{URI: "http://example.com/record2"})
	reg := &runtime.Registry{
		Canonical: repo,
	}
	msg := &message.Message{Entities: []string{"http://example.com/record2"}}
	action := NewTouchAction(reg)
	list, _ := action.(*TouchAction).recordToResourceList(msg)

	assert.Equal(t, "http://example.com/record2", list[0].Subject())
}
