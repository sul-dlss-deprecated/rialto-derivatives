package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss-labs/rialto-lambda/derivative"
	"github.com/sul-dlss-labs/rialto-lambda/models"
	"github.com/sul-dlss-labs/rialto-lambda/transform"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Running Lambda function\n")

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	client := derivative.NewSolrClient(host, collection)
	indexer := &transform.Indexer{}
	resourceList := []models.Resource{}

	for _, record := range snsEvent.Records {
		doc := recordToResource(record.SNS.Message)
		resourceList = append(resourceList, doc)
	}

	err := client.Add(indexer.Map(resourceList))
	if err != nil {
		panic(err)
	}
}

// This will take an SNS message and retrieve a resource from Neptune
func recordToResource(message string) models.Resource {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		log.Fatal(err)
	}

	return models.NewResource(data)
}

func main() {
	lambda.Start(Handler)
}
