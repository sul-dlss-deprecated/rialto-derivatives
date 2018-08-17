package runtime

import (
	"os"

	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// BuildServiceRegistry builds a Registry with the services initialized from the environment variables
func BuildServiceRegistry() *Registry {
	client := buildSolrClient()
	service := buildSparqlService()
	indexer := transform.NewCompositeIndexer(service)

	return NewRegistry(client, indexer, service)
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
