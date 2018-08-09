package models

import (
	"time"
)

// Publication models data stored in a database
type Publication struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	URI       string
	Metadata  publicationMetadata

	Authors []Person `pg:"many2many:people_publications"`
}

type publicationMetadata struct {
}

// Person models data stored in a database
type Person struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	URI       string
	Metadata  personMetadata

	Publications []Publication `pg:"many2many:people_publications"`
}

type personMetadata struct {
	Department               string
	InstitutionalAffiliation string
	Name                     string
}

// Organization models data stored in a database
type Organization struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	URI       string
	Metadata  organizationMetadata
}

type organizationMetadata struct {
	Department string
	Country    string
	Name       string
	Type       string
}

// ModelFromResourceType looks up resource type in type lists and returns name of model
func ModelFromResourceType(resourceType string) interface{} {
	for _, personType := range personTypes {
		if resourceType == personType {
			return Person{}
		}
	}
	for _, publicationType := range publicationTypes {
		if resourceType == publicationType {
			return Publication{}
		}
	}
	return nil
}

// AllModels returns a list of models
func AllModels() []interface{} {
	return []interface{}{
		Person{},
		Publication{},
		Organization{},
	}
}
