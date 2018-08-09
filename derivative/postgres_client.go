package derivative

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// PostgresClient represents the functions on the Postgres database
type PostgresClient struct {
	db *pg.DB
}

// NewPostgresClient returns a new PostGresClient instance
func NewPostgresClient(host string, database string, port string, username string, password string) *PostgresClient {
	db := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     fmt.Sprintf("%v:%v", host, port),
	})
	return &PostgresClient{
		db: db,
	}
}

// RemoveResourcesOfType clears the database of all the data with the matching model
func (d *PostgresClient) RemoveResourcesOfType(resourceType string) error {
	model := models.ModelFromResourceType(resourceType)
	if model == nil {
		return nil
	}
	_, err := d.db.Model(&model).Delete()
	return err
}

// RemoveAll clears known models out of the database
func (d *PostgresClient) RemoveAll() error {
	for _, model := range models.AllModels() {
		_, err := d.db.Model(&model).Delete()
		if err != nil {
			log.Printf("Could not remove model %v in RemoveAll", model)
		}
	}
	return nil
}

// Add inserts or updates rows in Postgres
func (d *PostgresClient) Add(rows []interface{}, model string) error {
	for _, row := range rows {

	}
	return nil
}
