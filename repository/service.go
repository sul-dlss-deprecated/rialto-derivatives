package repository

import (
	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// Repository is an interface that rialto-derivatives reads from as its source
type Repository interface {
	SubjectToResource(subject string) (models.Resource, error)
	AllResources() ([]models.Resource, error)
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

	data := map[string][]rdf.Term{}
	for _, triple := range response.Solutions() {
		predicate := triple["p"].String()
		object := triple["o"]

		if data[predicate] == nil {
			// First time we encounter a predicate
			data[predicate] = []rdf.Term{object}
		} else {
			// subsequent encounters
			data[predicate] = append(data[predicate], object)
		}
	}
	resource := models.NewResource(subject, data)
	return resource, nil
}

// AllResources returns a full list of resources
func (m *Service) AllResources() ([]models.Resource, error) {
	response, err := m.reader.QueryEverything()
	if err != nil {
		return nil, err
	}
	// Solutions look like this:
	// [map[o:AA00 s:http://rialto.stanford.edu/stanford p:http://rialto.stanford.edu/vocab/organizationCodes]]
	// log.Printf("Solutions %s", response.Solutions())
	list := m.toResourceList(m.groupBySubject(response.Solutions()))
	return list, nil
}

func (m *Service) toResourceList(grouped map[string]map[string][]rdf.Term) []models.Resource {
	list := []models.Resource{}
	for subject, predicates := range grouped {
		list = append(list, models.NewResource(subject, predicates))
	}
	return list
}

// Returns a map with subjects as keys and a map with predicates as values
func (m *Service) groupBySubject(raw []map[string]rdf.Term) map[string]map[string][]rdf.Term {
	result := map[string]map[string][]rdf.Term{}
	for _, triple := range raw {
		// Subjects and Predicates are always IRIs so it's safe to cast to String
		subject := triple["s"].String()
		predicate := triple["p"].String()
		object := triple["o"]

		if result[subject] == nil {
			// The first time we encounter a subject
			result[subject] = map[string][]rdf.Term{}
		}
		if result[subject][predicate] == nil {
			// First time we encounter a predicate
			result[subject][predicate] = []rdf.Term{object}
		} else {
			// subsequent encounters
			result[subject][predicate] = append(result[subject][predicate], object)
		}
	}
	return result
}
