package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

// BeforeInsert is a pg db hook that auto-initializes CreatedAt/UpdatedAt timestamps
func (o *Organization) BeforeInsert(db orm.DB) error {
	now := RailsTimestampNow()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now
	return nil
}

// BeforeUpdate is a pg db hook that auto-initializes UpdatedAt timestamps
func (o *Organization) BeforeUpdate(db orm.DB) error {
	o.UpdatedAt = RailsTimestampNow()
	return nil
}

// Organization models data stored in a database
type Organization struct {
	ID        int       `sql:",pk"`
	CreatedAt time.Time `sql:"type:timestamptz, default:now(), null"`
	UpdatedAt time.Time `sql:"type:timestamptz, default:now(), null"`
	URI       string    `sql:",unique"`
	Metadata  OrganizationMetadata
}

// OrganizationMetadata models jsonb data for Organizations
type OrganizationMetadata struct {
	Department string
	Country    string
	Name       string
	Type       string
}
