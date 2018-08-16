package models

import (
	"log"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
)

// railsTimestampLayout defines a Go time pattern matching how Rails db timestamps are formatted
const railsTimestampLayout = "2006-01-02 15:04:05.000000"

// RailsTimestampNow returns a time.Time object formatted per the Rails database layout for timestamps
func RailsTimestampNow() time.Time {
	now := time.Now().Format(railsTimestampLayout)
	timestamp, _ := time.Parse(railsTimestampLayout, now)
	return timestamp
}

// ModelWriter wraps a pg client with model-specific knowledge
type ModelWriter struct {
	Client *derivative.PostgresClient
}

// NewModelWriter returns a new ModelWriter instance
func NewModelWriter(client *derivative.PostgresClient) *ModelWriter {
	return &ModelWriter{
		Client: client,
	}
}

// allModels returns a list of models
func (w *ModelWriter) allModels() []interface{} {
	var person []Person
	var publication []Publication
	var organization []Organization

	return []interface{}{
		w.Client.Db.Model(&person),
		w.Client.Db.Model(&publication),
		w.Client.Db.Model(&organization),
	}
}

// AddOrganization adds an organization model instance
func (w *ModelWriter) AddOrganization(org *Organization) error {
	_, err := w.Client.Db.Model(org).Insert()
	return err
}

// GetOrganization retrieves an organization model instance
func (w *ModelWriter) GetOrganization(org *Organization) error {
	err := w.Client.Db.Model(org).Select()
	return err
}

// UpdateOrganization updates an organization model instance
func (w *ModelWriter) UpdateOrganization(org *Organization) error {
	orgForSelect := Organization{
		URI: org.URI,
	}
	err := w.GetOrganization(&orgForSelect)
	if err != nil {
		panic(err)
	}
	orgForUpdate :=
	_, err = w.Client.Db.
		Model(org).
		OnConflict("(uri) DO UPDATE").
		Set("metadata = EXCLUDED.metadata").
		Insert()
	return err
}

// RemoveAll removes all data from all models
func (w *ModelWriter) RemoveAll() error {
	for _, model := range w.allModels() {
		// Using Where("1=1") is a hack. WherePK() should work but does not.
		// With WherePK(), the generated query looks like: DELETE FROM "organizations" AS "organization" WHERE "organization"."id" = NULL
		// And this query causes a Postgres error: ``duplicate key value violates unique constraint "index_organizations_on_uri"``
		_, err := model.(*orm.Query).Where("1=1").Delete()
		if err != nil {
			log.Printf("Could not remove model %v in RemoveAll()", model)
		}
	}
	return nil
}
