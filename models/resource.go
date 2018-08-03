package models

import "github.com/knakk/rdf"

// Resource represents the data we get from Neptune
type Resource struct {
	Subject string
	data    map[string][]rdf.Term
}

var property = map[string]string{
	"type":           Predicates["rdf"]["type"],
	"abstract":       Predicates["bibo"]["abstract"],
	"doi":            Predicates["bibo"]["doi"],
	"cites":          Predicates["bibo"]["cites"],
	"identifier":     Predicates["bibo"]["identifier"],
	"link":           Predicates["bibo"]["uri"],
	"created":        Predicates["dct"]["created"],
	"journalIssue":   Predicates["dct"]["hasPart"],
	"subject":        Predicates["dct"]["subject"],
	"title":          Predicates["dct"]["title"],
	"alternateTitle": Predicates["dct"]["alternate"],
	"description":    Predicates["vivo"]["description"],
	"fundedBy":       Predicates["vivo"]["hasFundingVehicle"],
	"publisher":      Predicates["vivo"]["publisher"],
	"sponsor":        Predicates["vivo"]["informationResourceSupportedBy"],
	"hasInstrument":  Predicates["gcis"]["hasInstrument"],
	"sameAs":         Predicates["owl"]["sameAs"],
}

// NewResource creates a new instance of the resource
func NewResource(subject string, data map[string][]rdf.Term) Resource {
	return Resource{Subject: subject, data: data}
}

// IsPublication returns true if the type is a publiction
func (r *Resource) IsPublication() bool {
	resourceTypes := r.ValueOf("type")
	for _, pubType := range publicationTypes {
		for _, resType := range resourceTypes {
			if pubType == resType.String() {
				return true
			}
		}
	}
	return false
}

// ValueOf returns the rdf assertions for the predicate registerd under the supplied property name
func (r *Resource) ValueOf(name string) []rdf.Term {
	return r.data[property[name]]
}

// TODO: -- these go through a vivo:Authorship node.
// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.
