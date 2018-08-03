package models

// RdfTypePredicate is the predicate for rdf type
const RdfTypePredicate = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"

// AbstractPredicate is the predicate for abstract
const AbstractPredicate = "http://purl.org/ontology/bibo/abstract"

// CitesPredicate is the predicate for cites
const CitesPredicate = "http://purl.org/ontology/bibo/cites"

// DoiPredicate is the predicate for the DOI
const DoiPredicate = "http://purl.org/ontology/bibo/doi"

// IdentifierPredicate is the predicate for the identifier
const IdentifierPredicate = "http://purl.org/ontology/bibo/identifier"

// LinkPredicate is the predicate for the link
const LinkPredicate = "http://purl.org/ontology/bibo/uri"

// TitlePredicate is the rdf predicate we are using for title
const TitlePredicate = "http://purl.org/dc/terms/title"

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
