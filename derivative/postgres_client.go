package derivative

import (
	"fmt"

	"github.com/go-pg/pg"
)

// PostgresClient represents the functions on the Postgres database
type PostgresClient struct {
	Db *pg.DB
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
		Db: db,
	}
}
