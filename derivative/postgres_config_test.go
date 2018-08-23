package derivative

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConfig(t *testing.T) {
	conf := NewPostgresConfig().
		WithHost("123.45.67.89").
		WithPort("5479").
		WithDbname("my_instance").
		WithUser("rialto").
		WithPassword("seckret")

	assert.Equal(t, conf.toConnString(), "user=rialto host=123.45.67.89 port=5479 dbname=my_instance password=seckret ")
}

func TestPostgresConfigWithSSLDisabled(t *testing.T) {
	conf := NewPostgresConfig().
		WithSSL(false).
		WithDbname("my_instance")

	assert.Equal(t, conf.toConnString(), "dbname=my_instance sslmode=disable")
}
