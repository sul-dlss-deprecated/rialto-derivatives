package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss/rialto-derivatives/actions"
	"github.com/sul-dlss/rialto-derivatives/derivative"
	"github.com/sul-dlss/rialto-derivatives/message"
	"github.com/sul-dlss/rialto-derivatives/repository"
	"github.com/sul-dlss/rialto-derivatives/runtime"
	"github.com/sul-dlss/rialto-derivatives/transform"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	repo := repository.BuildRepository()
	registry := runtime.NewRegistry(repo, buildSolrClient(repo))
	for _, record := range sqsEvent.Records {
		msg, err := message.Parse(record)
		if err != nil {
			panic(err)
		}

		if err = actions.DispatchMessage(msg, registry).Run(msg); err != nil {
			panic(err)
		}
	}
}

func buildSolrClient(repo repository.Repository) *derivative.SolrClient {
	indexer := transform.NewCompositeIndexer(repo)

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	return derivative.NewSolrClient(host, collection, indexer)
}

func main() {
	lambda.Start(Handler)
}
