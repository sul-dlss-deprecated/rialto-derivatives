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
		msg, _ := message.Parse(record)
		err := actionForMessage(msg, registry).Run(msg)
		if err != nil {
			panic(err)
		}
	}
}

func buildServiceRegistry() *runtime.Registry {
	dbHost := getenvOrDefault("DB_HOST", "localhost")
	dbName := getenvOrDefault("DB_NAME", "rialto")
	dbPort := getenvOrDefault("DB_PORT", "5432")
	dbUsername := getenvOrDefault("DB_USERNAME", "foo")
	dbPassword := getenvOrDefault("DB_PASSWORD", "bar")
	dbClient := derivative.NewPostgresClient(dbHost, dbName, dbPort, dbUsername, dbPassword)
	dbTransformer := transform.NewDbTransformer()

	solrHost := os.Getenv("SOLR_HOST")
	solrCollection := os.Getenv("SOLR_COLLECTION")
	solrClient := derivative.NewSolrClient(solrHost, solrCollection)
	solrIndexer := transform.NewCompositeIndexer()

	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)

	return runtime.NewRegistry(dbClient, dbTransformer, solrClient, solrIndexer, sparqlReader)
}

func getenvOrDefault(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
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
	log.Panicf("Unknown action '%s'", msg.Action)
	return nil
}

func main() {
	lambda.Start(Handler)
}
