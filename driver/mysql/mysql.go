package mysql

import (
	"database/sql"
	"log"

	"github.com/mohammadv184/gloader/driver"

	// Import the mysql driver.
	_ "github.com/go-sql-driver/mysql"
)

// MySQL is a driver for MySQL.
type MySQL struct{}

func init() {
	err := driver.Register(&MySQL{})
	if err != nil {
		log.Println(err)
	}
}

// GetDriverName returns the name of the driver.
func (m *MySQL) GetDriverName() string {
	return "mysql"
}

// Open opens a connection to the database.
func (m *MySQL) Open(name string) (driver.Connection, error) {
	config, err := parseConfig(name)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("mysql", config.String())
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}
