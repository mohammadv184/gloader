package cockroach

import (
	"database/sql"
	"gloader/driver"
	"log"

	_ "github.com/lib/pq"
)

type Cockroach struct{}

func init() {
	err := driver.Register(&Cockroach{})
	if err != nil {
		log.Println(err)
	}
}

func (m *Cockroach) GetDriverName() string {
	return "cockroach"
}

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
