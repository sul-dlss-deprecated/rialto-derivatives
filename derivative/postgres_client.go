package derivative

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
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

// Add puts a bunch of documents in Solr
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

func (d *PostgresClient) retrieveOnePerson(subject string) (string, error) {
	return d.retrieveOneRecord("people", subject)
}

func (d *PostgresClient) retrieveOneOrganization(subject string) (string, error) {
	return d.retrieveOneRecord("organizations", subject)
}

func (d *PostgresClient) retrieveOnePublication(subject string) (string, error) {
	return d.retrieveOneRecord("publications", subject)
}

func (d *PostgresClient) retrievePeoplePublicationRelationship(subject string) (*[]string, error) {
	return d.retrieveRelationship("people_publications", "people", "person_id", "publications", "publication_id", subject)
}

func (d *PostgresClient) retrieveOneRecord(table string, subject string) (string, error) {
	query := fmt.Sprintf("SELECT metadata FROM %v WHERE uri = $1", table)
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

func (d *PostgresClient) retrieveOneKey(table string, uri string) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %v WHERE uri = $1", table)
	rows, err := d.DB.Query(query, uri)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	ret := rows.Next()
	// Make sure at least one row
	if ret == false {
		return 0, fmt.Errorf("No matches for %s in %s", uri, table)
	}
	var key int
	rows.Scan(&key)
	ret = rows.Next()
	// Make sure not more than one row
	if ret == true {
		return 0, fmt.Errorf("More than one match for %s in %s", uri, table)
	}

	return key, nil
}

func (d *PostgresClient) retrieveRelationship(manyTable string, selectTable string,
	selectTableJoinField string, whereTable string, whereTableJoinField, whereURI string) (*[]string, error) {

	query := fmt.Sprintf(`SELECT st.uri FROM "%v" st
		INNER JOIN "%v" mt ON st.id=mt.%v
		INNER JOIN "%v" wt ON mt.%v=wt.id
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
	return d.addResource("people", resource.Subject(), d.personSerializer.Serialize(resource))
}

func (d *PostgresClient) addOrganization(resource *models.Organization) error {
	return d.addResource("organizations", resource.Subject(), d.organizationSerializer.Serialize(resource))
}

func (d *PostgresClient) addPublication(resource *models.Publication) error {
	err := d.addResource("publications", resource.Subject(), d.publicationSerializer.Serialize(resource))
	if err != nil {
		return err
	}
	peopleURIs := make([]string, len(resource.Authors))
	for i, author := range resource.Authors {
		peopleURIs[i] = author.URI
	}
	return d.addRelationship("people_publications", "publication_id", "publications", resource.URI, "person_id", "people", peopleURIs)
}

func (d *PostgresClient) addResource(table string, subject string, data string) error {
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uri) DO UPDATE SET metadata=$2, updated_at=$4 WHERE %v.uri=$1`, table, table)
	_, err := d.DB.Exec(sql, subject, data, time.Now(), time.Now())
	return err
}

func (d *PostgresClient) addRelationship(table string, oneField string, oneTable string, one string, manyField string, manyTable string, many []string) error {
	// First, find the keys of the one and manies.
	oneKey, err := d.retrieveOneKey(oneTable, one)
	if err != nil {
		return err
	}
	manyKeys := make([]int, len(many))
	var key int
	for i, manyValue := range many {
		key, err = d.retrieveOneKey(manyTable, manyValue)
		if err != nil {
			return err
		}
		manyKeys[i] = key
	}
	deleteSQL := fmt.Sprintf(`DELETE FROM "%s" WHERE %s=$1`, table, oneField)
	_, err = d.DB.Exec(deleteSQL, oneKey)
	if err != nil {
		return err
	}
	// Uggh, there's no bulk/batch insert in database/sql: https://github.com/golang/go/issues/5171
	insertSQL := fmt.Sprintf(`INSERT INTO "%s" ("%s", "%s") VALUES ($1, $2)`, table, oneField, manyField)
	insertStmt, err := d.DB.Prepare(insertSQL)
	if err != nil {
		return err
	}
	for _, manyKey := range manyKeys {
		_, err := insertStmt.Exec(oneKey, manyKey)
		if err != nil {
			return err
		}
	}
	return nil
}
