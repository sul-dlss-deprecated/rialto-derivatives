package message

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Message a message (from SQS) that the application understands
type Message struct {
	Action   string
	Entities []string
	Body     interface{}
}

// SQSBody is the body of an SQS message
type SQSBody struct {
	Message string
}

// ParseSQS transforms an SQSMessage into a message
func ParseSQS(record events.SQSMessage) (*Message, error) {
	body := &SQSBody{}
	err := json.Unmarshal([]byte(record.Body), body)
	if err != nil {
		return nil, err
	}
	msg := &Message{}
	err = json.Unmarshal([]byte(body.Message), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// ParseSNS transforms a SNSEventRecord into a message
func ParseSNS(record events.SNSEventRecord) (*Message, error) {
	data := record.SNS.Message
	msg := &Message{}
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
