package repository

import (
	"fmt"
	"log"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// Repository is an interface that rialto-derivatives reads from as its source
type Repository interface {
	SubjectToResource(subject string) (models.Resource, error)
	AllResources(func([]models.Resource) error) error
}

// Service is the Neptune implementation of the repository
type Service struct {
	reader Reader
}

// NewService creates a new Service instance
func NewService(reader Reader) Repository {
	return &Service{reader: reader}
}

// SubjectToResource takes a subject string and returns a resource
func (m *Service) SubjectToResource(subject string) (models.Resource, error) {
	response, err := m.reader.QueryByID(subject)

	if err != nil {
		return nil, err
	}

	log.Printf("Solutions: %v", response.Solutions())
	list := m.toResourceList(response.Solutions())
	if len(list) == 0 {
		return nil, fmt.Errorf("Record not found: %s", subject)
	}

	return list[0], nil
}

// AllResources returns a full list of resources
func (m *Service) AllResources(f func([]models.Resource) error) error {
	err := m.reader.QueryEverything(func(response *sparql.Results) error {
		// Solutions look like this:
		// [map[o:AA00 s:http://rialto.stanford.edu/stanford p:http://rialto.stanford.edu/vocab/organizationCodes]]
		// log.Printf("Solutions %s", response.Solutions())
		list := m.toResourceList(response.Solutions())
		return f(list)
	})
	return err
}

func (m *Service) toResourceList(solutions []map[string]rdf.Term) []models.Resource {
	list := []models.Resource{}
	for _, solution := range solutions {
		resource := models.NewResource(solution)
		if v, ok := resource.(models.Person); ok {
			// People also need to be informed of their organization membership.
			results, err := m.reader.GetOrganizationInfo(v.Organization)
			if err != nil {
				panic(err)
			}
			v.SetOrganizationInfo(results)
		}
		list = append(list, resource)
	}
	return list
}
