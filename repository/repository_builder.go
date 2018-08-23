package repository

import (
	"os"
)

func BuildRepository() Repository {
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	sparqlReader := NewSparqlReader(endpoint)
	return NewService(sparqlReader)
}
