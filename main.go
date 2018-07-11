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
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Running Lambda function\n")
	registry := buildServiceRegistry()

	for _, record := range snsEvent.Records {
		msg, _ := message.Parse(record)
		err := actionForMessage(msg, registry).Run(msg)
		if err != nil {
			panic(err)
		}
	}
}

func buildServiceRegistry() *runtime.Registry {
	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	client := derivative.NewSolrClient(host, collection)
	indexer := &transform.Indexer{}
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)

	return runtime.NewRegistry(client, indexer, sparqlReader)
}

func actionForMessage(msg *message.Message, registry *runtime.Registry) actions.Action {
	switch msg.Action {
	case "touch":
		return actions.NewTouchAction(registry)
	case "delete":
		return actions.NewRebuildAction(registry)
	}
	log.Panicf("Unknown action '%s'", msg.Action)
	return nil
}

func main() {
	lambda.Start(Handler)
}
