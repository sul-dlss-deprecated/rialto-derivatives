package models

import (
	"time"
	//	"github.com/go-pg/pg/orm"
)

// BeforeInsert is a pg db hook that auto-initializes CreatedAt/UpdatedAt timestamps
// func (p *Publication) BeforeInsert(db orm.DB) error {
// 	now := time.Now().String()
// 	if p.CreatedAt.IsZero() {
// 		created, err := time.Parse(RailsTimestampFormat, now)
// 		if err != nil {
// 			panic(err)
// 		}
// 		p.CreatedAt = created
// 	}

// 	if p.UpdatedAt.IsZero() {
// 		updated, err := time.Parse(RailsTimestampFormat, now)
// 		if err != nil {
// 			panic(err)
// 		}
// 		p.UpdatedAt = updated
// 	}

// 	return nil
// }

// BeforeUpdate is a pg db hook that auto-initializes UpdatedAt timestamps
// func (p *Publication) BeforeUpdate(db orm.DB) error {
// 	updated, err := time.Parse(RailsTimestampFormat, time.Now().String())
// 	if err != nil {
// 		panic(err)
// 	}
// 	p.UpdatedAt = updated
// 	return nil
// }

// Publication models data stored in a database
type Publication struct {
	ID        int       `sql:",pk"`
	CreatedAt time.Time `sql:"type:timestamptz, default:now(), null"`
	UpdatedAt time.Time `sql:"type:timestamptz, default:now(), null"`
	URI       string    `sql:",unique"`
	Metadata  PublicationMetadata
	Authors   []Person `pg:"many2many:people_publications"`
}

// PublicationMetadata models jsonb data for Publications
type PublicationMetadata struct {
}
