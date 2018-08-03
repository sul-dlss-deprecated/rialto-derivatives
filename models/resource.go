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
	return r.data[Predicates["rdf"]["type"]]
}

// Titles returns the title assertions
func (r *Resource) Titles() []rdf.Term {
	return r.data[Predicates["dct"]["title"]]
}

// Abstract returns the bibo:abstract assertions
func (r *Resource) Abstract() []rdf.Term {
	return r.data[Predicates["bibo"]["abstract"]]
}

// DOI returns the bibo:doi assertions
func (r *Resource) DOI() []rdf.Term {
	return r.data[Predicates["bibo"]["doi"]]
}

// Cites returns the bibo:cites assertions
func (r *Resource) Cites() []rdf.Term {
	return r.data[Predicates["bibo"]["cites"]]
}

// Identifier returns the bibo:identifier assertions
func (r *Resource) Identifier() []rdf.Term {
	return r.data[Predicates["bibo"]["identifier"]]
}

// Link returns the bibo:uri assertions
func (r *Resource) Link() []rdf.Term {
	return r.data[Predicates["bibo"]["uri"]]
}

// TODO: -- these go through a vivo:Authorship node.
// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.

// date of creation 	dct:created 	DateTime string, EDTF 	[1,1] 	Used to describe the creation date of a resource.
// description 	vivo:description 	string-literal 	[0,n] 	Description of the resource.
// funded by 	vivo:hasFundingVehicle 	Grant URI 	[0,n] 	Grant (or contract) providing funding for the publication.
// has instrument 	gcis:hasInstrument 	gcis:Instrument URI 	[0,n] 	A type of tool or device used for a particular task, especially for scientific work, as presented in the publication (specifically for Datasets).
// journal issue 	dcterms:hasPart 	Document URI (Article) 	[0,n] 	Journal is another entity with issue number, label / title, possibly isPartOf URI for the Journal title overall.
// publisher 	vivo:publisher 	URI for foaf:Organization 	[0,n] 	Publisher of the resource.
// sameAs 	owl:sameAs 	URI 	[0,n] 	Other resources (identified via URIs) that are the same as this resource.
// sponsor 	vivo:informationResourceSupportedBy 	Agent URI 	[0,n] 	Institution supporting the publication.
// subject 	dcterms:subject 	Topic / Concept URI 	[0,n] 	Topic or concept the resource is about.
// title 	dcterms:title 	string-literal 	[1,1] 	Title for the resource.
// alternate title
