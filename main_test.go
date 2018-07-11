package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-lambda/actions"
	"github.com/sul-dlss-labs/rialto-lambda/message"
)

func TestRecordToDoc(t *testing.T) {
	msg := &message.Message{
		Action: "touch",
	}
	action := actionForMessage(msg, nil, nil)

	assert.IsType(t, &actions.TouchAction{}, action)
}
