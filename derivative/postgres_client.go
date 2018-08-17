package derivative

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
)

// PostgresClient represents the functions on the Postgres derivative tables
type PostgresClient struct {
	db        *sql.DB
	canonical repository.Repository
}

// NewPostgresClient returns a new PostgresClient instance
func NewPostgresClient(config *PostgresConfig, repo repository.Repository) *PostgresClient {
	db, err := sql.Open("postgres", config.toConnString())
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresClient{
		db:        db,
		canonical: repo,
	}
}

// RemoveResourcesOfType clears the index of all the data with the matching type
// func (d *PostgresClient) RemoveResourcesOfType(resourceType string) error {
// 	query := fmt.Sprintf("%s:%s", typeField, resourceType)
// 	data := map[string]interface{}{"query": query}
// 	_, err := d.si.Delete(data, nil)
// 	return err
// }

// RemoveAll clears the index of all the data
func (d *PostgresClient) RemoveAll() error {
	_, err := d.db.Exec(`TRUNCATE TABLE people_publications`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`TRUNCATE TABLE people`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`TRUNCATE TABLE publications`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`TRUNCATE TABLE organizations`)
	return err
}

// Add puts a bunch of documents in Solr
func (d *PostgresClient) Add(resourceList []models.Resource) error {
	for _, resource := range resourceList {
		d.addOne(resource)
	}
	return nil
}

func (d *PostgresClient) addOne(resource models.Resource) error {
	if resource.IsPerson() {
		d.addPerson(resource)
	} else if resource.IsOrganization() {
		d.addOrganization(resource)
	} else {
		return fmt.Errorf("Unrecognized resource type: %v", resource)
	}

	return nil
}

func (d *PostgresClient) retrieveOnePerson(subject string) (string, error) {
	return d.retrieveOneRecord("people", subject)
}

func (d *PostgresClient) retrieveOneOrganization(subject string) (string, error) {
	return d.retrieveOneRecord("organizations", subject)
}

func (d *PostgresClient) retrieveOneRecord(table string, subject string) (string, error) {
	query := fmt.Sprintf("SELECT metadata FROM %v WHERE uri = $1", table)
	rows, err := d.db.Query(query, subject)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var name string
	rows.Next()
	rows.Scan(&name)
	return name, nil
}

func (d *PostgresClient) addPerson(resource models.Resource) error {
	return d.addResource("people", resource.Subject(), d.personMetadata(resource))
}

func (d *PostgresClient) addOrganization(resource models.Resource) error {
	return d.addResource("organizations", resource.Subject(), d.organizationMetadata(resource))
}

func (d *PostgresClient) addResource(table string, subject string, data string) error {
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "metadata", "created_at", "updated_at") VALUES ($1, $2, $3, $4) RETURNING "id"`, table)
	_, err := d.db.Exec(sql, subject, data, time.Now(), time.Now())
	return err
}

// return the Person resource as a JSON string.  Must include the following properties:
//   name (string)
//   department (URI)
//   institutionalAffiliation (URI)
func (d *PostgresClient) personMetadata(resource models.Resource) string {
	return fmt.Sprintf(`{"name": "%s"}`, d.retrieveAssociatedName(resource))
}

// return the Organization resource as a JSON string.  Must include the following properties:
//   name (string)
//   type (URI) the most specific type (e.g. Department or University)
func (d *PostgresClient) organizationMetadata(resource models.Resource) string {
	return fmt.Sprintf(`{"name": "%s"}`, resource.ValueOf("orgName")[0])
}

// TODO: This method is copied from PersonIndexer.  In order to be more efficient,
// we should lookup names before passing to the postgres/solr writers.
func (d *PostgresClient) retrieveAssociatedName(resource models.Resource) string {
	nameURI := resource.ValueOf("name")
	if len(nameURI) == 0 {
		log.Printf("No name URI found for %s", resource.Subject())
		return ""
	}

	nameResource, err := d.canonical.SubjectToResource(nameURI[0].String())

	if err != nil {
		panic(err)
	}
	givenName := nameResource.ValueOf("given-name")
	familyName := nameResource.ValueOf("family-name")

	if len(givenName) == 0 || len(familyName) == 0 {
		return ""
	}
	return fmt.Sprintf("%v %v", givenName[0], familyName[0])
}
