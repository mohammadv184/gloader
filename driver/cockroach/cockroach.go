package cockroach

import (
	"database/sql"
	"gloader/driver"
	"log"

	// Import the postgres driver.
	_ "github.com/lib/pq"
)

// Cockroach is a driver for CockroachDB.
type Cockroach struct{}

func init() {
	err := driver.Register(&Cockroach{})
	if err != nil {
		log.Println(err)
	}
}

// GetDriverName returns the name of the driver.
func (m *Cockroach) GetDriverName() string {
	return "cockroach"
}

// Open opens a connection to the database.
func (m *Cockroach) Open(dsn string) (driver.Connection, error) {
	config, err := parseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("postgres", config.String())
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}
