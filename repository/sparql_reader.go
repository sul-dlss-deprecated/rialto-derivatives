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

// QueryThroughNode is similar to a HasAndBelongsTo query in an RDBMS where
// we're querying though some intermittent node for the data we need.
// e.g. Person -> <relatedBy> -> Position -> <relates> -> Department
func (r *SparqlReader) QueryThroughNode(id string, localPredicate string, localType string, remotePredicate string) (*sparql.Results, error) {
	query := fmt.Sprintf(`select ?d where
		{ <%s> <%s> ?o .
			?o a <%s> .
			?o <%s> ?d } limit 10`, id, localPredicate, localType, remotePredicate)
	return r.repo.Query(query)
}

// QueryByIDAndPredicate returns the list of objects that match the given subject, predicate tuple
func (r *SparqlReader) QueryByIDAndPredicate(id string, predicate string) (*sparql.Results, error) {
	query := fmt.Sprintf(`select ?o where
		{ <%s> <%s> ?o . } limit 10`, id, predicate)
	return r.repo.Query(query)
}
