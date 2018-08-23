package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
)

func TestTouchAction(t *testing.T) {
	msg := &message.Message{
		Action: "touch",
	}
	action := DispatchMessage(msg, nil)

	assert.IsType(t, &TouchAction{}, action)
}

func TestRebuildAction(t *testing.T) {
	msg := &message.Message{
		Action: "rebuild",
	}
	action := DispatchMessage(msg, nil)

	assert.IsType(t, &RebuildAction{}, action)
}
