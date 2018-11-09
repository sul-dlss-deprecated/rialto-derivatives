package repository

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/knakk/sparql"
)

// Number of records to pull per requst.  If this is too large, then the SPARQL proxy
// lambda will hit a 6mb limit. See: https://github.com/sul-dlss/sparql-loader/issues/44
// When this was set to anything over 10000, it failed.
const tripleLimit = 9000
const idVariable = "?id"

// Number of times to retry connecting to the triplestore
const retries = 90

// SparqlReader represents the functions we do on the triplestore
type SparqlReader struct {
	repo SparqlRepository
}

// SparqlRepository is a repository that handles sparql queries
type SparqlRepository interface {
	Query(q string) (*sparql.Results, error)
}

// RetryingSparqlRepository is a wrapper for a SparqlRepository that will retry failed connections
type RetryingSparqlRepository struct {
	repo *sparql.Repo
}

// Query wraps querying a SparqlRepository that will retry failed connections
func (r *RetryingSparqlRepository) Query(q string) (*sparql.Results, error) {
	results, err := r.repo.Query(q)
	// Retrying because when running under Localstack, triplestore may not be immediately available.
	for i := 0; i < retries && err != nil && strings.Contains(err.Error(), "no such host"); i++ {
		log.Printf("Retrying repo, time %d", i+1)
		time.Sleep(1 * time.Second)
		results, err = r.repo.Query(q)
	}
	return results, err
}

// NewSparqlReader creates a new instance of the sparqlReader for the provided endpoint
func NewSparqlReader(url string) *SparqlReader {
	repo, err := sparql.NewRepo(url,
		sparql.Timeout(time.Second*60),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &SparqlReader{repo: &RetryingSparqlRepository{repo: repo}}
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
		// log.Printf("[SPARQL] %s", query)
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
			return fmt.Sprintf(`SELECT ?id ?type ?firstname ?lastname
			WHERE {
				?id a ?type .
				%s
				FILTER ( ?type = <http://xmlns.com/foaf/0.1/Person>)
				OPTIONAL {
					?id <http://www.w3.org/2006/vcard/ns#hasName> ?n .
					?n <http://www.w3.org/2006/vcard/ns#given-name> ?firstname .
					?n <http://www.w3.org/2006/vcard/ns#family-name> ?lastname .
				}
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// GetPositionOrganizationInfo retrieves a hierarchical list of organizations the given organization subject is part of
func (r *SparqlReader) GetPositionOrganizationInfo(personID string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT DISTINCT ?org ?type ?name
			WHERE {
				?org rdfs:label ?label .
				?org <http://www.w3.org/2004/02/skos/core#prefLabel> ?name .
    		?pos_org a <http://xmlns.com/foaf/0.1/Organization> .
				?pos_org <http://purl.obolibrary.org/obo/BFO_0000050>* ?org .
    		?pos_org <http://vivoweb.org/ontology/core#relatedBy> ?pos .
    		?pos a <http://vivoweb.org/ontology/core#Position> .
    		?pos <http://vivoweb.org/ontology/core#relates> <%s> .
			     OPTIONAL {
				      ?org a ?type .
			 	      FILTER ( ?type NOT IN (<http://xmlns.com/foaf/0.1/Organization>, <http://xmlns.com/foaf/0.1/Agent>))

           }
			}
			ORDER BY ?org OFFSET 0 LIMIT 100`, personID)
	return r.repo.Query(query)
}

// GetCountriesInfo retrieves a list of countries for a person
func (r *SparqlReader) GetCountriesInfo(personID string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?country ?label
		 WHERE {
				 <%s> <http://purl.org/dc/terms/spatial> ?country .
				 ?country rdfs:label ?label .
			}
			ORDER BY ?country OFFSET 0 LIMIT 100`, personID)
	return r.repo.Query(query)
}

// GetPersonSubtypesInfo retrieves a list of types for a person
func (r *SparqlReader) GetPersonSubtypesInfo(personID string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?subtype
		 WHERE {
				 <%s> a ?subtype .
				 FILTER ( ?subtype NOT IN (<http://xmlns.com/foaf/0.1/Person>, <http://xmlns.com/foaf/0.1/Agent>))
			}
			ORDER BY ?subtype OFFSET 0 LIMIT 100`, personID)
	return r.repo.Query(query)
}

// GetAuthorInfo retrieves a list of authors the given publication subject is part of
func (r *SparqlReader) GetAuthorInfo(publication string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?id ?label
		 WHERE {
				 <%s> <http://vivoweb.org/ontology/core#relatedBy> ?authorship .
				 ?authorship a <http://vivoweb.org/ontology/core#Authorship> .
				 ?authorship <http://vivoweb.org/ontology/core#relates> ?id .
				 ?id <http://www.w3.org/2004/02/skos/core#prefLabel> ?label .
			}
			ORDER BY ?author OFFSET 0 LIMIT 100`, publication)
	return r.repo.Query(query)
}

// GetConceptInfo retrieves a list of concepts for a subject
func (r *SparqlReader) GetConceptInfo(subject string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?id ?label
		 WHERE {
				 <%s> <http://purl.org/dc/terms/subject> ?id .
				 ?id a <http://www.w3.org/2004/02/skos/core#Concept> .
				 ?id <http://purl.org/dc/terms/subject> ?label .
			}
			ORDER BY ?id OFFSET 0 LIMIT 100`, subject)
	return r.repo.Query(query)
}

// GetGrantInfo retrieves a list of grants for a publication
func (r *SparqlReader) GetGrantInfo(publication string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?id ?label
		 WHERE {
				 <%s> <http://vivoweb.org/ontology/core#hasFundingVehicle> ?id .
				 ?id a <http://vivoweb.org/ontology/core#Grant> .
				 ?id <http://www.w3.org/2004/02/skos/core#prefLabel> ?label .
			}
			ORDER BY ?id OFFSET 0 LIMIT 100`, publication)
	return r.repo.Query(query)
}

// GetIdentifierInfo retrieves a list of identifiers for a publication
func (r *SparqlReader) GetIdentifierInfo(publication string) (*sparql.Results, error) {
	query := fmt.Sprintf(`SELECT ?id
		 WHERE {
				 <%s> <http://purl.org/ontology/bibo/identifier> ?id .
			}
			ORDER BY ?id OFFSET 0 LIMIT 100`, publication)
	return r.repo.Query(query)
}

func (r *SparqlReader) queryOrganizations(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?subtype ?name ?parent
			WHERE {
			  ?id a <http://xmlns.com/foaf/0.1/Organization> .
				?id a ?type .
				?id <http://www.w3.org/2004/02/skos/core#prefLabel> ?name .
				%s
				FILTER ( ?type = <http://xmlns.com/foaf/0.1/Organization>)
			  OPTIONAL {
					?id a ?subtype .
					FILTER ( ! bound(?subtype) || ?subtype NOT IN (<http://xmlns.com/foaf/0.1/Organization>, <http://xmlns.com/foaf/0.1/Agent>))
				}
				OPTIONAL {
					?id <http://purl.obolibrary.org/obo/BFO_0000050> ?parent .
			  }
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

func (r *SparqlReader) queryGrants(f func(*sparql.Results) error, ids ...string) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?id ?type ?name ?pi ?pi_label ?assigned ?assigned_label
			WHERE {
			  ?id a <http://vivoweb.org/ontology/core#Grant> .
				?id a ?type .
			  ?id <http://www.w3.org/2004/02/skos/core#prefLabel>|<http://www.w3.org/2000/01/rdf-schema#label> ?name .
				?id <http://vivoweb.org/ontology/core#relates> ?pi_role .
				?pi_role a <http://vivoweb.org/ontology/core#PrincipalInvestigatorRole> .
				?pi_role <http://purl.obolibrary.org/obo/RO_0000052> ?pi .
				?pi <http://www.w3.org/2004/02/skos/core#prefLabel> ?pi_label .
				?id <http://vivoweb.org/ontology/core#assignedBy> ?assigned .
				?assigned <http://www.w3.org/2004/02/skos/core#prefLabel> ?assigned_label .
				%s
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

//vivo:relates vivo:PrincipalInvestigatorRole obo:RO_0000052 foaf:Person

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
			  ?id <http://purl.org/dc/terms/subject> ?label .
				%s
			}
			ORDER BY ?id OFFSET %v LIMIT %v`, r.filter(ids), offset, tripleLimit)
		}, f)
}

// 	Publication resource types yet to map:
// 	"cites":          Predicates["bibo"]["cites"],
// 	"link":           Predicates["bibo"]["uri"],
// 	"journalIssue":   Predicates["dct"]["hasPart"],
// 	"subject":        Predicates["dct"]["subject"],
// 	"alternateTitle": Predicates["dct"]["alternate"],
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
			return fmt.Sprintf(`SELECT ?id ?type ?subtype ?title ?abstract ?doi
				?identifier ?publisher ?description ?created
			WHERE {
			  ?id a <http://purl.org/ontology/bibo/Document> .
				?id a ?type .
			  ?id <http://purl.org/dc/terms/title> ?title .
				?id <http://purl.org/dc/terms/created> ?created .
				%s
				FILTER ( ?type = <http://purl.org/ontology/bibo/Document>)
				OPTIONAL {
					?id a ?subtype .
					FILTER ( ?subtype != <http://purl.org/ontology/bibo/Document> )
				}
				OPTIONAL {
					?id <http://purl.org/ontology/bibo/abstract> ?abstract .
				}
				OPTIONAL {
					?id <http://purl.org/ontology/bibo/doi> ?doi .
				}
				OPTIONAL {
					?id <http://vivoweb.org/ontology/core#publisher> ?publisher .
				}
				OPTIONAL {
					?id <http://vivoweb.org/ontology/core#description> ?description .
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

// QueryByIDs returns a list of SPARQL results for all the triples that match the
// provided subjects the datastore. If there are no results for a particular subject,
// then it will be removed from the list.  Therefore, the input list may be
// longer than the output list.
func (r *SparqlReader) QueryByIDs(ids []string) ([]*sparql.Results, error) {
	results := []*sparql.Results{}
	for _, id := range ids {
		// First step is to find out what type this object has:
		doctype, err := r.queryTypeForID(id)
		if err != nil {
			return nil, err
		}
		if doctype == "" {
			// No registered type was found
			continue
		}
		result, err := r.QueryByTypeAndID(doctype, id)
		if err != nil {
			return nil, err
		}
		// If result is nil don't append it.
		if result == nil {
			// Non-Stanford Organizations aren't expected to have a more specific type.
			// (e.g. Grant funders or Publishers)
			if doctype != "http://xmlns.com/foaf/0.1/Organization" {
				log.Printf("No results found for %s and %s", doctype, id)
			}
			continue
		}

		results = append(results, result)
	}
	return results, nil
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
	solutions := res.Solutions()
	if len(solutions) < 1 {
		log.Printf("No type found for %s", id)
		return "", err
	}
	t := solutions[0]["type"].String()
	log.Printf("%s is a %s", id, t)
	return t, nil
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
		return nil, fmt.Errorf("No registered type '%s' (%v)", doctype, id)
	}
	return retval, nil
}
