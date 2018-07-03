package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vanng822/go-solr/solr"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Running Lambda function\n")

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	si, _ := solr.NewSolrInterface(host, collection)

	for _, record := range snsEvent.Records {
		doc := recordToDoc(record.SNS.Message)
		docList := []solr.Document{doc}
		params := url.Values{}
		si.Add(docList, 1000, &params)
		si.Commit()
	}
}

func recordToDoc(message string) solr.Document {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		log.Fatal(err)
	}

	return solr.Document(data)
}

func main() {
	lambda.Start(Handler)
}
