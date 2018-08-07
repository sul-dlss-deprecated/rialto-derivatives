package models

import "github.com/knakk/rdf"

// Resource is the interface type for resources
type Resource interface {
	ValueOf(name string) []rdf.Term
	IsConcept() bool
	IsGrant() bool
	IsOrganization() bool
	IsPerson() bool
	IsProject() bool
	IsPublication() bool
	Subject() string
}

// RdfBackedResource represents the data we get from Neptune
type RdfBackedResource struct {
	subject string
	data    map[string][]rdf.Term
}

var property = map[string]string{
	// All resource types:
	"type": Predicates["rdf"]["type"],

	// Publication resource types
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

	// Grant resources
	"grantName": Predicates["skos"]["prefLabel"],

	// For person resources
	"name": Predicates["vcard"]["hasName"],

	// For name resources
	"given-name":  Predicates["vcard"]["given-name"],
	"family-name": Predicates["vcard"]["family-name"],

	// Organization resources
	"orgName": Predicates["skos"]["prefLabel"],

	// Project resources
	"hasStartDate": Predicates["frapo"]["hasStartDate"],
	"hasEndDate":   Predicates["frapo"]["hasEndDate"],

	// Concept resources
	"label": Predicates["skos"]["prefLabel"],
}

// NewResource creates a new instance of the resource
func NewResource(subject string, data map[string][]rdf.Term) Resource {
	return &RdfBackedResource{subject: subject, data: data}
}

// Subject returns the identifier of the resource
func (r *RdfBackedResource) Subject() string {
	return r.subject
}

// IsPublication returns true if the type is a publiction
func (r *RdfBackedResource) IsPublication() bool {
	return r.isTypeIn(publicationTypes)
}

// IsPerson returns true if the type is a person
func (r *RdfBackedResource) IsPerson() bool {
	return r.isTypeIn(personTypes)
}

// IsOrganization returns true if the type is a organization
func (r *RdfBackedResource) IsOrganization() bool {
	return r.isTypeIn(organizationTypes)
}

// IsGrant returns true if the type is a grant
func (r *RdfBackedResource) IsGrant() bool {
	return r.isTypeIn(grantTypes)
}

// IsProject returns true if the type is a project
func (r *RdfBackedResource) IsProject() bool {
	return r.isTypeIn(projectTypes)
}

// IsConcept returns true if the type is a concept
func (r *RdfBackedResource) IsConcept() bool {
	return r.isTypeIn(conceptTypes)
}

// isTypeIn returns true if the resource type is in the provided list
func (r *RdfBackedResource) isTypeIn(desiredTypes []string) bool {
	resourceTypes := r.ValueOf("type")
	for _, desired := range desiredTypes {
		for _, resType := range resourceTypes {
			if desired == resType.String() {
				return true
			}
		}
	}
	return false
}

// ValueOf returns the rdf assertions for the predicate registerd under the supplied property name
func (r RdfBackedResource) ValueOf(name string) []rdf.Term {
	return r.data[property[name]]
}

// TODO: -- these go through a vivo:Authorship node.
// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.
