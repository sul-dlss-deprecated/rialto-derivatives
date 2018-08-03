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
	"http://vivoweb.org/ontology/core#Abstract",
	"http://purl.org/ontology/bibo/Article",
	"http://purl.org/ontology/bibo/Book",
	"http://vivoweb.org/ontology/core#CaseStudy",
	"http://vivoweb.org/ontology/core#Catalog",
	"http://purl.org/spar/fabio/ClinicalGuideline",
	"http://vivoweb.org/ontology/core#ConferencePoster",
	"http://purl.org/ontology/bibo/Manual",
	"http://purl.org/ontology/bibo/Manuscript",
	"http://purl.org/ontology/bibo/Patent",
	"http://purl.org/ontology/bibo/Report",
	"http://vivoweb.org/ontology/core#ResearchProposal",
	"http://vivoweb.org/ontology/core#Score",
	"http://vivoweb.org/ontology/core#Screenplay",
	"http://purl.org/ontology/bibo/Slideshow",
	"http://vivoweb.org/ontology/core#Speech",
	"http://purl.org/ontology/bibo/Standard",
	"http://purl.org/ontology/bibo/Thesis",
	"http://vivoweb.org/ontology/core#Translation",
	"http://purl.org/ontology/bibo/Webpage",
	"http://vivoweb.org/ontology/core#WorkingPaper",
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
