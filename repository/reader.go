package repository

import "github.com/knakk/sparql"

// Reader reads from the data store
type Reader interface {
	QueryEverything() (*sparql.Results, error)
	QueryByID(id string) (*sparql.Results, error)
}
