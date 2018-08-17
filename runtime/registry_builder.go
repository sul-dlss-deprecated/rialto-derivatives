package runtime

import (
	"bytes"
	"os"

	// Added for the postgres driver
	_ "github.com/lib/pq"

	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
)

// BuildServiceRegistry builds a Registry with the services initialized from the environment variables
func BuildServiceRegistry() *Registry {
	service := buildSparqlService()
	solr := buildSolrClient(service)
	db := buildDatabase()
	writer := derivative.NewCompositeWriter(db, solr)
	return NewRegistry(service, writer)
}

func buildSolrClient(service *repository.Service) *derivative.SolrClient {
	indexer := transform.NewCompositeIndexer(service)

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	return derivative.NewSolrClient(host, collection, indexer)
}

func buildSparqlService() *repository.Service {
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)
	return repository.NewService(sparqlReader)
}

func buildDatabase() *derivative.PostgresClient {
	conf := derivative.NewPostgresConfig().
		WithUser(os.Getenv("RDS_USERNAME")).
		WithPassword(os.Getenv("RDS_PASSWORD")).
		WithDbname(os.Getenv("RDS_DB_NAME")).
		WithHost(os.Getenv("RDS_HOSTNAME")).
		WithPort(os.Getenv("RDS_PORT"))

	return derivative.NewPostgresClient(conf)
}

func addField(b bytes.Buffer, field string, variable string) {
	val := os.Getenv(variable)
	if len(val) != 0 {
		b.WriteString(field)
		b.WriteString("=")
		b.WriteString(val)
		b.WriteString(" ")
	}
}
