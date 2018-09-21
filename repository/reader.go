package repository

import (
	"github.com/knakk/sparql"
)

// Reader reads from the data store
type Reader interface {
	QueryEverything(f func(*sparql.Results) error) error
	QueryByID(id string) (*sparql.Results, error)
	GetOrganizationInfo(id *string) (*sparql.Results, error)
}
