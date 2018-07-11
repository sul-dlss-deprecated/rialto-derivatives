package actions

import (
	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/runtime"
)

// RebuildAction drops the repository and then rebuilds it
type RebuildAction struct {
	registry *runtime.Registry
}

// NewRebuildAction creates a Rebuild (delete) action
func NewRebuildAction(registry *runtime.Registry) Action {
	return &RebuildAction{registry: registry}
}

// Run does the rebuilding
func (r *RebuildAction) Run(message *message.Message) error {
	err := r.registry.Derivatives.RemoveAll()
	if err != nil {
		return err
	}
	resourceList, err := r.queryResources()
	if err != nil {
		return err
	}

	err = r.registry.Derivatives.Add(r.registry.Indexer.Map(resourceList))
	return err
}

// Return a list of resources populated by querying for everything in the triplestore
func (r *RebuildAction) queryResources() ([]models.Resource, error) {
	response, err := r.registry.Canonical.QueryEverything()
	if err != nil {
		return nil, err
	}
	// Solutions look like this:
	// [map[o:AA00 s:http://rialto.stanford.edu/stanford p:http://rialto.stanford.edu/vocab/organizationCodes]]
	// log.Printf("Solutions %s", response.Solutions())
	list := r.toResourceList(r.groupBySubject(response.Solutions()))
	return list, nil
}

func (r *RebuildAction) toResourceList(grouped map[string]map[string][]rdf.Term) []models.Resource {
	list := []models.Resource{}
	for subject, predicates := range grouped {
		list = append(list, models.NewResource(subject, predicates))
	}
	return list
}

// Returns a map with subjects as keys and a map with predicates as values
func (r *RebuildAction) groupBySubject(raw []map[string]rdf.Term) map[string]map[string][]rdf.Term {
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
