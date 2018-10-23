package repository

import (
	"github.com/knakk/sparql"
)

// Reader reads from the data store
type Reader interface {
	QueryEverything(f func(*sparql.Results) error) error
	QueryByIDs(ids []string) ([]*sparql.Results, error)
	GetPositionOrganizationInfo(personID string) (*sparql.Results, error)
	GetCountriesInfo(personID string) (*sparql.Results, error)
	GetAuthorInfo(publicationID string) (*sparql.Results, error)
	GetConceptInfo(subject string) (*sparql.Results, error)
}
