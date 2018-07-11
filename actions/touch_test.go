package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-lambda/message"
)

func TestRecordToDoc(t *testing.T) {
	msg := &message.Message{}
	action := &TouchAction{}
	doc := action.recordToResource(msg)

	assert.Equal(t, "whatever", doc.Title())
}
