package derivative

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// PostgresClient represents the functions on the Postgres derivative tables
type PostgresClient struct {
	db *sql.DB
}

// NewPostgresClient returns a new PostgresClient instance
func NewPostgresClient(config *PostgresConfig) *PostgresClient {
	db, err := sql.Open("postgres", config.toConnString())
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresClient{
		db: db,
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
	} else {
		return fmt.Errorf("Unrecognized resource type: %v", resource)
	}

	return nil
}

func (d *PostgresClient) addPerson(resource models.Resource) error {
	_, err := d.db.Exec(`INSERT INTO "people" ("uri", "metadata", "created_at", "updated_at") VALUES ($1, $2, $3, $4) RETURNING "id"`, "http://foo.com/123", "{\"title\":\"hello\"}", time.Now(), time.Now())
	return err
}
