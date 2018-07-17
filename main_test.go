package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/actions"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
)

func TestTouchAction(t *testing.T) {
	msg := &message.Message{
		Action: "touch",
	}
	action := actionForMessage(msg, nil)

	assert.IsType(t, &actions.TouchAction{}, action)
}

func TestRebuildAction(t *testing.T) {
	msg := &message.Message{
		Action: "rebuild",
	}
	action := actionForMessage(msg, nil)

	assert.IsType(t, &actions.RebuildAction{}, action)
}
