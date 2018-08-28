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
	db                     *sql.DB
	personSerializer       *transform.PersonSerializer
	organizationSerializer *transform.OrganizationSerializer
}

// NewPostgresClient returns a new PostgresClient instance
func NewPostgresClient(config *PostgresConfig, repo repository.Repository) *PostgresClient {
	db, err := sql.Open("postgres", config.toConnString())
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresClient{
		db:                     db,
		personSerializer:       transform.NewPersonSerializer(repo),
		organizationSerializer: &transform.OrganizationSerializer{},
	}
}

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
	switch v := resource.(type) {
	case *models.Person:
		d.addPerson(v)
	case *models.Organization:
		d.addOrganization(v)
	default:
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

func (d *PostgresClient) addPerson(resource *models.Person) error {
	return d.addResource("people", resource.Subject(), d.personSerializer.Serialize(resource))
}

func (d *PostgresClient) addOrganization(resource *models.Organization) error {
	return d.addResource("organizations", resource.Subject(), d.organizationSerializer.Serialize(resource))
}

func (d *PostgresClient) addResource(table string, subject string, data string) error {
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "metadata", "created_at", "updated_at") VALUES ($1, $2, $3, $4) RETURNING "id"`, table)
	_, err := d.db.Exec(sql, subject, data, time.Now(), time.Now())
	return err
}
