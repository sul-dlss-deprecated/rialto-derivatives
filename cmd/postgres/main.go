package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sul-dlss/rialto-derivatives/actions"
	"github.com/sul-dlss/rialto-derivatives/derivative"
	"github.com/sul-dlss/rialto-derivatives/message"
	"github.com/sul-dlss/rialto-derivatives/repository"
	"github.com/sul-dlss/rialto-derivatives/runtime"

	// Added for the postgres driver
	_ "github.com/lib/pq"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context, sqsEvent events.SQSEvent) {
	repo := repository.BuildRepository()
	registry := runtime.NewRegistry(repo, buildDatabase(repo))
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

func buildDatabase(repo repository.Repository) *derivative.PostgresClient {
	conf := derivative.NewPostgresConfig().
		WithUser(os.Getenv("RDS_USERNAME")).
		WithPassword(os.Getenv("RDS_PASSWORD")).
		WithDbname(os.Getenv("RDS_DB_NAME")).
		WithHost(os.Getenv("RDS_HOSTNAME")).
		WithPort(os.Getenv("RDS_PORT")).
		WithSSL(os.Getenv("RDS_SSL") == "" || strings.ToLower(os.Getenv("RDS_SSL")) == "true")

	return derivative.NewPostgresClient(conf, repo)
}

func main() {
	lambda.Start(Handler)
}
