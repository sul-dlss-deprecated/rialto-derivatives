package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss-labs/rialto-lambda/derivative"
	"github.com/sul-dlss-labs/rialto-lambda/message"
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

	for _, record := range snsEvent.Records {
		msg, _ := message.Parse(record)
		switch msg.Action {
		case "touch":
			updateOne(msg, client, indexer)
		case "delete":
			rebuildIndex(msg, client, indexer)
		}
	}
}

func updateOne(msg *message.Message, client *derivative.SolrClient, indexer *transform.Indexer) {
	resourceList := []models.Resource{}
	resourceList = append(resourceList, recordToResource(msg))

	err := client.Add(indexer.Map(resourceList))
	if err != nil {
		panic(err)
	}
}

func rebuildIndex(msg *message.Message, client *derivative.SolrClient, indexer *transform.Indexer) {

}

// This will take an SNS message and retrieve a resource from Neptune
func recordToResource(msg *message.Message) models.Resource {
	// TODO: data should come from neptune
	data := map[string]interface{}{"rdf:type": "bar", "dc:title": "whatever"}
	return models.NewResource(data)
}

func main() {
	lambda.Start(Handler)
}
