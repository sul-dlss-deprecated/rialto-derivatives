package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/knakk/sparql"
)

const tripleLimit = 100000
const idVariable = "?id"

// SparqlReader represents the functions we do on the triplestore
type SparqlReader struct {
	repo *sparql.Repo
}

// NewSparqlReader creates a new instance of the sparqlReader for the provided endpoint
func NewSparqlReader(url string) *SparqlReader {
	repo, err := sparql.NewRepo(url,
		sparql.Timeout(time.Second*60),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &SparqlReader{repo: repo}
}

// QueryEverything returns all triples in the datastore
// TODO: this could overflow our memory if we get too many records in the store.
//       Should we set a lower limit and paginate the result set?
func (r *SparqlReader) QueryEverything(f func(*sparql.Results) error) error {
	if err := r.queryPeople(f); err != nil {
		return err
	}
	if err := r.queryOrganizations(f); err != nil {
		return err
	}
	if err := r.queryGrants(f); err != nil {
		return err
	}
	if err := r.queryConcepts(f); err != nil {
		return err
	}

	return r.queryPublications(f)
}

// Calls sparqlForOffset once for each page to generate the query
// Calls f() on each page of results
func (r *SparqlReader) queryPage(sparqlForOffset func(offset int) string, f func(*sparql.Results) error) error {
	page := 0
	for {
		offset := page * tripleLimit
		query := sparqlForOffset(offset)
		log.Printf("[SPARQL] %s", query)
		results, err := r.repo.Query(query)
		if err != nil {
			return err
		}
		if resultCount := len(results.Solutions()); resultCount == 0 {
			break
		}
		if err = f(results); err != nil {
			return err
		}
		page++
	}
	return nil
}

func (r *SparqlReader) queryPeople(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?subtype ?firstname ?lastname ?org
			WHERE {
				?id a ?type .
				?id a <http://xmlns.com/foaf/0.1/Person> .
				?id a ?subtype .
				%s
				FILTER ( ?subtype NOT IN (<http://xmlns.com/foaf/0.1/Person>, <http://xmlns.com/foaf/0.1/Agent>))
				FILTER ( ?type = <http://xmlns.com/foaf/0.1/Person>)
				OPTIONAL {
					?id <http://www.w3.org/2006/vcard/ns#hasName> ?n .
					?n <http://www.w3.org/2006/vcard/ns#given-name> ?firstname .
					?n <http://www.w3.org/2006/vcard/ns#family-name> ?lastname .
				}
				OPTIONAL {
					?id <http://vivoweb.org/ontology/core#relatedBy> ?pos .
					?pos a <http://vivoweb.org/ontology/core#Position> .
					?pos <http://vivoweb.org/ontology/core#relates> ?org .
					?org a <http://xmlns.com/foaf/0.1/Organization> .
				}
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// GetOrganizationInfo retrieves a hierarchical list of organizations the give organization subject is part of
func (r *SparqlReader) GetOrganizationInfo(org *string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?org ?type ?name
  						 WHERE {
				FILTER ( ?id = <%s> )
			 	FILTER ( ?type NOT IN (<http://xmlns.com/foaf/0.1/Organization>, <http://xmlns.com/foaf/0.1/Agent>))

        ?id <http://purl.obolibrary.org/obo/BFO_0000050>* ?org .
        ?org <http://www.w3.org/2004/02/skos/core#prefLabel> ?name .
        ?org a ?type .
			}
			ORDER BY ?org OFFSET 0 LIMIT 100`, *org)
	return r.repo.Query(query)
}

func (r *SparqlReader) queryOrganizations(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?subtype ?name ?altLabel ?abbreviation ?parent
			WHERE {
			  ?id a <http://xmlns.com/foaf/0.1/Organization> .
				?id a ?type .
				?id <http://www.w3.org/2004/02/skos/core#prefLabel>|<http://www.w3.org/2000/01/rdf-schema#label> ?name .
			  ?id a ?subtype .
				%s
			  FILTER ( ?subtype NOT IN (<http://xmlns.com/foaf/0.1/Organization>, <http://xmlns.com/foaf/0.1/Agent>))
				FILTER ( ?type = <http://xmlns.com/foaf/0.1/Organization>)
			  OPTIONAL {
			    ?id <http://www.w3.org/2004/02/skos/core#altLabel> ?altLabel .
					?id <http://vivoweb.org/ontology/core#abbreviation> ?abbreviation .
					?id <http://purl.obolibrary.org/obo/BFO_0000050> ?parent .
			  }
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

func (r *SparqlReader) queryGrants(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?name
			WHERE {
			  ?id a <http://vivoweb.org/ontology/core#Grant> .
				?id a ?type .
			  ?id <http://www.w3.org/2004/02/skos/core#prefLabel>|<http://www.w3.org/2000/01/rdf-schema#label> ?name .
				%s
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// filter creates a SPARQL filter fragement to match only on the passed id.
func (r *SparqlReader) filter(ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return fmt.Sprintf("FILTER ( ?id = <%s> )", ids[0])
}

func (r *SparqlReader) queryConcepts(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?label
			WHERE {
			  ?id a <http://www.w3.org/2004/02/skos/core#Concept> .
				?id a ?type .
			  ?id <http://www.w3.org/2004/02/skos/core#prefLabel>|<http://www.w3.org/2000/01/rdf-schema#label> ?label .
				%s
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// 	Publication resource types
// 	"abstract":       Predicates["bibo"]["abstract"],
// 	"doi":            Predicates["bibo"]["doi"],
// 	"cites":          Predicates["bibo"]["cites"],
// 	"identifier":     Predicates["bibo"]["identifier"],
// 	"link":           Predicates["bibo"]["uri"],
// 	"created":        Predicates["dct"]["created"],
// 	"journalIssue":   Predicates["dct"]["hasPart"],
// 	"subject":        Predicates["dct"]["subject"],
// 	"title":          Predicates["dct"]["title"],
// 	"alternateTitle": Predicates["dct"]["alternate"],
// 	"description":    Predicates["vivo"]["description"],
// 	"fundedBy":       Predicates["vivo"]["hasFundingVehicle"],
// 	"publisher":      Predicates["vivo"]["publisher"],
// 	"sponsor":        Predicates["vivo"]["informationResourceSupportedBy"],
// 	"hasInstrument":  Predicates["gcis"]["hasInstrument"],
// 	"sameAs":         Predicates["owl"]["sameAs"],
// these go through a vivo:Authorship node.
// author 	vivo:relatedBy vivo:Authorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Author of the publication.
// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.
func (r *SparqlReader) queryPublications(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?subtype ?title
			WHERE {
			  ?id a <http://purl.org/ontology/bibo/Document> .
				?id a ?type .
			  ?id <http://purl.org/dc/terms/title> ?title .
				%s
				OPTIONAL {
					?id a ?subtype .
					FILTER ( ?subtype != <http://purl.org/ontology/bibo/Document> )
				}
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// Project resources
func (r *SparqlReader) queryProjects(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?title ?startdate ?enddate
				WHERE {
				  ?id a <http://xmlns.com/foaf/0.1/Project> .
					?id a ?type .
				  ?id <http://purl.org/dc/terms/title> ?title .
					%s
					OPTIONAL {
						?id <http://purl.org/cerif/frapo/hasStartDate> ?startdate .
						?id <http://purl.org/cerif/frapo/hasEndDate> ?enddate .
					}
				}
				ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// QueryByID returns all triples that match the subject the datastore
func (r *SparqlReader) QueryByID(id string) (*sparql.Results, error) {
	// First step is to find out what type this object has:
	doctype, err := r.queryTypeForID(id)
	if err != nil {
		return nil, err
	}
	return r.QueryByTypeAndID(doctype, id)
}

func (r *SparqlReader) queryTypeForID(id string) (string, error) {
	query := fmt.Sprintf(`SELECT ?type WHERE {
		<%s> a ?type .
		FILTER ( ?type IN (<http://xmlns.com/foaf/0.1/Organization>,<http://xmlns.com/foaf/0.1/Person>,<http://vivoweb.org/ontology/core#Grant>,<http://www.w3.org/2004/02/skos/core#Concept>,<http://purl.org/ontology/bibo/Document>,<http://xmlns.com/foaf/0.1/Project>))
	}
	LIMIT 1`, id)
	res, err := r.repo.Query(query)
	if err != nil {
		return "", err
	}
	return res.Solutions()[0]["type"].String(), nil
}

// QueryByTypeAndID issues the query for the provided type with the given ID
func (r *SparqlReader) QueryByTypeAndID(doctype string, id string) (*sparql.Results, error) {
	var retval *sparql.Results
	copyToLocal := func(arg1 *sparql.Results) error {
		retval = arg1
		return nil
	}
	switch t := doctype; t {
	case "http://xmlns.com/foaf/0.1/Organization":
		r.queryOrganizations(copyToLocal, id)
	case "http://xmlns.com/foaf/0.1/Person":
		r.queryPeople(copyToLocal, id)
	case "http://vivoweb.org/ontology/core#Grant":
		r.queryGrants(copyToLocal, id)
	case "http://www.w3.org/2004/02/skos/core#Concept":
		r.queryConcepts(copyToLocal, id)
	case "http://purl.org/ontology/bibo/Document":
		r.queryPublications(copyToLocal, id)
	case "http://xmlns.com/foaf/0.1/Project":
		r.queryProjects(copyToLocal, id)
	default:
		return nil, fmt.Errorf("No type for %v", id)
	}
	return retval, nil
}
