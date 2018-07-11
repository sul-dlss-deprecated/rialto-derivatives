package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss-labs/rialto-derivatives/actions"
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Running Lambda function\n")

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	client := derivative.NewSolrClient(host, collection)
	indexer := &transform.Indexer{}

	for _, record := range snsEvent.Records {
		msg, _ := message.Parse(record)
		err := actionForMessage(msg, client, indexer).Run(msg)
		if err != nil {
			panic(err)
		}
	}
}

func actionForMessage(msg *message.Message, client *derivative.SolrClient, indexer *transform.Indexer) actions.Action {
	switch msg.Action {
	case "touch":
		return actions.NewTouchAction(client, indexer)
	case "delete":
		return actions.NewRebuildAction(client, indexer)
	}
	log.Panicf("Unknown action '%s'", msg.Action)
	return nil
}

func main() {
	lambda.Start(Handler)
}
