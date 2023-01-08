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
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Connection{conn: conn}, nil
}
