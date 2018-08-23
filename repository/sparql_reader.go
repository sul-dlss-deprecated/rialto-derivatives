package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/knakk/sparql"
)

const tripleLimit = 10000000

// SparqlReader represents the functions we do on the triplestore
type SparqlReader struct {
	repo *sparql.Repo
}

// NewSparqlReader creates a new instance of the sparqlReader for the provided endpoint
func NewSparqlReader(url string) *SparqlReader {
	repo, err := sparql.NewRepo(url,
		sparql.Timeout(time.Second*10),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &SparqlReader{repo: repo}
}

// QueryEverything returns all triples in the datastore
// TODO: this could overflow our memory if we get too many records in the store.
//       Should we set a lower limit and paginate the result set?
func (r *SparqlReader) QueryEverything() (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?s ?p ?o
		WHERE { ?s ?p ?o .
			?s a ?type .
			FILTER (?type IN (<http://xmlns.com/foaf/0.1/Agent>,<http://vivoweb.org/ontology/core#Grant>,<http://www.w3.org/2004/02/skos/core#Concept>)} LIMIT %v`, tripleLimit)
	return r.repo.Query(query)
}

// QueryByID returns all triples that match the subject the datastore
func (r *SparqlReader) QueryByID(id string) (*sparql.Results, error) {
	query := fmt.Sprintf("SELECT ?p ?o WHERE { <%s> ?p ?o } LIMIT %v", id, tripleLimit)
	return r.repo.Query(query)
}
