package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss-labs/rialto-derivatives/actions"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	registry := runtime.BuildServiceRegistry()

	for _, record := range snsEvent.Records {
		msg, err := message.Parse(record)
		if err != nil {
			panic(err)
		}

		if err = actionForMessage(msg, registry).Run(msg); err != nil {
			panic(err)
		}
	}
}

func actionForMessage(msg *message.Message, registry *runtime.Registry) actions.Action {
	switch msg.Action {
	case "touch":
		log.Printf("Running Touch action\n")

		return actions.NewTouchAction(registry)
	case "rebuild":
		log.Printf("Running Rebuild action\n")
		return actions.NewRebuildAction(registry)
	}
	log.Panicf("Unknown action '%s'. Allowed actions: touch, rebuild.", msg.Action)
	return nil
}

func main() {
	lambda.Start(Handler)
}
