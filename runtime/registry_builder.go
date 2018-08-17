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
	repo := buildRepository()
	solr := buildSolrClient(repo)
	db := buildDatabase(repo)
	writer := derivative.NewCompositeWriter(db, solr)
	return NewRegistry(repo, writer)
}

func buildSolrClient(repo repository.Repository) *derivative.SolrClient {
	indexer := transform.NewCompositeIndexer(repo)

	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	return derivative.NewSolrClient(host, collection, indexer)
}

func buildRepository() repository.Repository {
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := repository.NewSparqlReader(endpoint)
	return repository.NewService(sparqlReader)
}

func buildDatabase(repo repository.Repository) *derivative.PostgresClient {
	conf := derivative.NewPostgresConfig().
		WithUser(os.Getenv("RDS_USERNAME")).
		WithPassword(os.Getenv("RDS_PASSWORD")).
		WithDbname(os.Getenv("RDS_DB_NAME")).
		WithHost(os.Getenv("RDS_HOSTNAME")).
		WithPort(os.Getenv("RDS_PORT"))

	return derivative.NewPostgresClient(conf, repo)
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
