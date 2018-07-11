package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

func TestRecordToResource(t *testing.T) {
	fakeSparql := new(MockedReader)

	reg := &runtime.Registry{
		Canonical: fakeSparql,
	}
	msg := &message.Message{Entities: []string{"http://example.com/record2"}}
	action := NewTouchAction(reg)
	doc, _ := action.(*TouchAction).recordToResource(msg)

	assert.Equal(t, "http://example.com/record2", doc.Subject)
}
