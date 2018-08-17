package derivative

import "bytes"

// PostgresConfig provides configuration for the database connection
type PostgresConfig struct {
	User     *string
	Password *string
	Dbname   *string
	Host     *string
	Port     *string
	SSL      bool
}

// NewPostgresConfig creates a new instance of the config
func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{SSL: true}
}

func (c *PostgresConfig) toConnString() string {
	var b bytes.Buffer
	if c.User != nil {
		b.WriteString("user=")
		b.WriteString(*c.User)
		b.WriteString(" ")
	}
	if c.Host != nil {
		b.WriteString("host=")
		b.WriteString(*c.Host)
		b.WriteString(" ")
	}
	if c.Port != nil {
		b.WriteString("port=")
		b.WriteString(*c.Port)
		b.WriteString(" ")
	}
	if c.Dbname != nil {
		b.WriteString("dbname=")
		b.WriteString(*c.Dbname)
		b.WriteString(" ")
	}
	if c.Password != nil {
		b.WriteString("password=")
		b.WriteString(*c.Password)
		b.WriteString(" ")
	}

	if !c.SSL {
		b.WriteString("sslmode=disable")
	}
	return b.String() // "dbname=rialto_development sslmode=disable"
}

// WithUser sets the user on the config
func (c *PostgresConfig) WithUser(user string) *PostgresConfig {
	if len(user) > 0 {
		c.User = &user
	}
	return c
}

// WithPassword sets the password on the config
func (c *PostgresConfig) WithPassword(password string) *PostgresConfig {
	if len(password) > 0 {
		c.Password = &password
	}
	return c
}

// WithDbname sets the dbname on the config
func (c *PostgresConfig) WithDbname(dbname string) *PostgresConfig {
	if len(dbname) > 0 {
		c.Dbname = &dbname
	}
	return c
}

// WithHost sets the host on the config
func (c *PostgresConfig) WithHost(host string) *PostgresConfig {
	if len(host) > 0 {
		c.Host = &host
	}
	return c
}

// WithPort sets the port on the config
func (c *PostgresConfig) WithPort(port string) *PostgresConfig {
	if len(port) > 0 {
		c.Port = &port
	}
	return c
}

// WithSSL sets the ssl mode on the config
func (c *PostgresConfig) WithSSL(ssl bool) *PostgresConfig {
	c.SSL = ssl
	return c
}
