package message

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Message a message (from SNS) that the application understands
type Message struct {
	Action   string
	Entities []string
	Body     interface{}
}

// Parse transforms a SNSEventRecord into a message
func Parse(record events.SNSEventRecord) (*Message, error) {
	data := record.SNS.Message
	msg := &Message{}
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
