package models

import (
	"time"
	//	"github.com/go-pg/pg/orm"
)

// BeforeInsert is a pg db hook that auto-initializes CreatedAt/UpdatedAt timestamps
// func (p *Person) BeforeInsert(db orm.DB) error {
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
// func (p *Person) BeforeUpdate(db orm.DB) error {
// 	updated, err := time.Parse(RailsTimestampFormat, time.Now().String())
// 	if err != nil {
// 		panic(err)
// 	}
// 	p.UpdatedAt = updated
// 	return nil
// }

// Person models data stored in a database
type Person struct {
	ID           int       `sql:",pk"`
	CreatedAt    time.Time `sql:"type:timestamptz, default:now(), null"`
	UpdatedAt    time.Time `sql:"type:timestamptz, default:now(), null"`
	URI          string    `sql:",unique"`
	Metadata     PersonMetadata
	Publications []Publication `pg:"many2many:people_publications"`
}

// PersonMetadata models jsonb data for Persons
type PersonMetadata struct {
	Department               string
	InstitutionalAffiliation string
	Name                     string
}
