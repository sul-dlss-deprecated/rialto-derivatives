package repository

import (
	"os"
)

// BuildRepository instantiates the repository from the environment variables
func BuildRepository() Repository {
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := NewSparqlReader(endpoint)
	return NewService(sparqlReader)
}
