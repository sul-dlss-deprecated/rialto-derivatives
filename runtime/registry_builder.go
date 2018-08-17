package runtime

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// BuildServiceRegistry builds a Registry with the services initialized from the environment variables
func BuildServiceRegistry() *Registry {
	client := buildSolrClient()
	service := buildSparqlService()
	indexer := transform.NewCompositeIndexer(service)
	db := buildDatabase()
	return NewRegistry(client, indexer, service, db)
}

func buildSolrClient() *derivative.SolrClient {
	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	return derivative.NewSolrClient(host, collection)
}

func buildSparqlService() *repository.Service {
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)
	return repository.NewService(sparqlReader)
}

func buildDatabase() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     os.Getenv("RDS_USERNAME"),
		Password: os.Getenv("RDS_PASSWORD"),
		Database: os.Getenv("RDS_DB_NAME"),
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("RDS_HOSTNAME"), os.Getenv("RDS_PORT")),
	})
}
