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
	registry := buildServiceRegistry()

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

func buildServiceRegistry() *runtime.Registry {
	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	client := derivative.NewSolrClient(host, collection)
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)
	service := repository.NewService(sparqlReader)
	indexer := transform.NewCompositeIndexer(service)

	return runtime.NewRegistry(client, indexer, service)
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
