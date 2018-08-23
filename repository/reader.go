package repository

import (
	"github.com/knakk/sparql"
)

// Reader reads from the data store
type Reader interface {
	QueryEverything() (*sparql.Results, error)
	QueryByID(id string) (*sparql.Results, error)
	QueryByIDAndPredicate(id string, predcate string) (*sparql.Results, error)
	QueryThroughNode(id string, localPredicate string, localType string, remotePredicate string) (*sparql.Results, error)
}
