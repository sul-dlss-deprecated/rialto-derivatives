package message

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestParseSQS(t *testing.T) {
	evtRecord := events.SQSMessage{
		Body: "{\"Message\": \"{\\\"Action\\\": \\\"touch\\\", \\\"Entities\\\":[\\\"http://example.com/foo1\\\"] }\"}",
	}

	event, _ := ParseSQS(evtRecord)
	assert.Equal(t, "touch", event.Action)
	assert.Equal(t, []string{"http://example.com/foo1"}, event.Entities)

}

func TestParseSNS(t *testing.T) {
	evtRecord := events.SNSEventRecord{
		SNS: events.SNSEntity{
			Message: "{\"Action\": \"touch\", \"Entities\":[\"http://example.com/foo1\"] }",
		},
	}
	event, _ := ParseSNS(evtRecord)
	assert.Equal(t, "touch", event.Action)
	assert.Equal(t, []string{"http://example.com/foo1"}, event.Entities)

}
