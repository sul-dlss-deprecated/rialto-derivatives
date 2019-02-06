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
	GetIdentifierInfo(subject string) (*sparql.Results, error)
	GetGrantIdentifierInfo(subject string) (*sparql.Results, error)
	GetGrantInfo(subject string) (*sparql.Results, error)
	GetPersonSubtypesInfo(subject string) (*sparql.Results, error)
	GetStanfordAuthorCount(subject string) (*sparql.Results, error)
}
