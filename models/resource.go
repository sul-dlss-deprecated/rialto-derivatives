package models

import "github.com/knakk/rdf"

// Resource represents the data we get from Neptune
type Resource struct {
	Subject string
	data    map[string][]rdf.Term
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

// Abstract returns the bibo:abstract assertions
func (r *Resource) Abstract() []rdf.Term {
	return r.data[AbstractPredicate]
}

// DOI returns the bibo:doi assertions
func (r *Resource) DOI() []rdf.Term {
	return r.data[DoiPredicate]
}

// Cites returns the bibo:cites assertions
func (r *Resource) Cites() []rdf.Term {
	return r.data[CitesPredicate]
}

// Identifier returns the bibo:identifier assertions
func (r *Resource) Identifier() []rdf.Term {
	return r.data[IdentifierPredicate]
}

// Link returns the bibo:uri assertions
func (r *Resource) Link() []rdf.Term {
	return r.data[LinkPredicate]
}

// TODO: author, Profiles, editor -- these go through a vivo:Authorship node.
