package models

import "github.com/knakk/rdf"

// Resource represents the data we get from Neptune
type Resource struct {
	Subject string
	data    map[string][]rdf.Term
}

// RdfTypePredicate is the predicate for rdf type
const RdfTypePredicate = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"

// TitlePredicate is the rdf predicate we are using for title
const TitlePredicate = "http://xmlns.com/foaf/0.1/name"

var publicationTypes = []string{
	"http://purl.org/ontology/bibo/Document",
}

// NewResource creates a new instance of the resource
func NewResource(subject string, data map[string][]rdf.Term) Resource {
	return Resource{Subject: subject, data: data}
}

// IsPublication returns true if the type is a publiction
func (r *Resource) IsPublication() bool {
	resourceTypes := r.ResourceTypes()
	for _, pubType := range publicationTypes {
		for _, resType := range resourceTypes {
			if pubType == resType.String() {
				return true
			}
		}
	}
	return false
}

// ResourceTypes returns the type assertions
func (r *Resource) ResourceTypes() []rdf.Term {
	return r.data[RdfTypePredicate]
}

// Titles returns the title assertions
func (r *Resource) Titles() []rdf.Term {
	return r.data[TitlePredicate]
}
