package derivative

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/repository"
	"github.com/sul-dlss/rialto-derivatives/transform"
)

// PostgresClient represents the functions on the Postgres derivative tables
type PostgresClient struct {
	DB                     *sql.DB
	personSerializer       *transform.PersonSerializer
	organizationSerializer *transform.OrganizationSerializer
	publicationSerializer  *transform.PublicationSerializer
}

// NewPostgresClient returns a new PostgresClient instance
func NewPostgresClient(config *PostgresConfig, repo repository.Repository) *PostgresClient {
	db, err := sql.Open("postgres", config.ToConnString())
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresClient{
		DB:                     db,
		personSerializer:       transform.NewPersonSerializer(repo),
		organizationSerializer: &transform.OrganizationSerializer{},
		publicationSerializer:  &transform.PublicationSerializer{},
	}
}

// RemoveAll clears the index of all the data
func (d *PostgresClient) RemoveAll() error {
	_, err := d.DB.Exec(`TRUNCATE TABLE people_publications`)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(`TRUNCATE TABLE people`)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(`TRUNCATE TABLE publications`)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(`TRUNCATE TABLE organizations`)
	return err
}

// Add puts a bunch of documents into Postgres
func (d *PostgresClient) Add(resourceList []models.Resource) error {
	for _, resource := range resourceList {
		err := d.addOne(resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *PostgresClient) addOne(resource models.Resource) error {
	switch v := resource.(type) {
	case *models.Person:
		return d.addPerson(v)
	case *models.Organization:
		return d.addOrganization(v)
	case *models.Publication:
		return d.addPublication(v)
	default:
		return fmt.Errorf("Unrecognized resource type: %v", resource)
	}
}

func (d *PostgresClient) retrieveOnePerson(subject string) (string, string, error) {
	return d.retrieveOneRecord("people", subject)
}

func (d *PostgresClient) retrieveOneOrganization(subject string) (string, string, error) {
	return d.retrieveOneRecord("organizations", subject)
}

func (d *PostgresClient) retrieveOnePublication(subject string) (string, error) {
	query := fmt.Sprintf("SELECT metadata FROM publications WHERE uri = $1")
	rows, err := d.DB.Query(query, subject)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var metadata string
	rows.Next()
	rows.Scan(&metadata)
	return metadata, nil
}

func (d *PostgresClient) retrievePeoplePublicationRelationship(subject string) (*[]string, error) {
	return d.retrieveRelationship("people_publications", "people", "person_uri", "publications", "publication_uri", subject)
}

func (d *PostgresClient) retrieveOneRecord(table string, subject string) (string, string, error) {
	query := fmt.Sprintf("SELECT name, metadata FROM %v WHERE uri = $1", table)
	rows, err := d.DB.Query(query, subject)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	var name string
	var metadata string
	rows.Next()
	rows.Scan(&name, &metadata)
	return name, metadata, nil
}

func (d *PostgresClient) retrieveRelationship(manyTable string, selectTable string,
	selectTableJoinField string, whereTable string, whereTableJoinField, whereURI string) (*[]string, error) {

	query := fmt.Sprintf(`SELECT st.uri FROM "%v" st
		INNER JOIN "%v" mt ON st.uri=mt.%v
		INNER JOIN "%v" wt ON mt.%v=wt.uri
		WHERE wt.uri=$1`, selectTable, manyTable, selectTableJoinField, whereTable, whereTableJoinField)
	rows, err := d.DB.Query(query, whereURI)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	uris := make([]string, 0)
	for rows.Next() {
		var uri string
		if err := rows.Scan(&uri); err != nil {
			return nil, err
		}
		uris = append(uris, uri)
	}

	return &uris, nil
}

func (d *PostgresClient) addPerson(resource *models.Person) error {
	return d.addResource(d.personSerializer.SQLForInsert(resource))
}

func (d *PostgresClient) addOrganization(resource *models.Organization) error {
	return d.addResource(d.organizationSerializer.SQLForInsert(resource))
}

func (d *PostgresClient) addPublication(resource *models.Publication) error {
	err := d.addResource(d.publicationSerializer.SQLForInsert(resource))
	if err != nil {
		return err
	}
	peopleURIs := make([]string, len(resource.Authors))
	for i, author := range resource.Authors {
		peopleURIs[i] = author.URI
	}
	return d.addRelationship("people_publications", "publication_uri", resource.URI, "person_uri", peopleURIs)
}

func (d *PostgresClient) addResource(sql string, vals []interface{}) error {
	_, err := d.DB.Exec(sql, vals...)
	return err
}

func (d *PostgresClient) addRelationship(table string, oneField string, one string, manyField string, many []string) error {
	deleteSQL := fmt.Sprintf(`DELETE FROM "%s" WHERE %s=$1`, table, oneField)
	_, err := d.DB.Exec(deleteSQL, one)
	if err != nil {
		return err
	}
	// Uggh, there's no bulk/batch insert in database/sql: https://github.com/golang/go/issues/5171
	insertSQL := fmt.Sprintf(`INSERT INTO "%s" ("%s", "%s") VALUES ($1, $2)`, table, oneField, manyField)
	insertStmt, err := d.DB.Prepare(insertSQL)
	if err != nil {
		return err
	}
	for _, manyItem := range many {
		_, err := insertStmt.Exec(one, manyItem)
		if err != nil {
			return err
		}
	}
	return nil
}
