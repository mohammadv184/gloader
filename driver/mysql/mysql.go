package mysql

import (
	"database/sql"
	"gloader/driver"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct{}

func init() {
	err := driver.Register(&MySQL{})
	if err != nil {
		log.Println(err)
	}
}

func (m *MySQL) GetDriverName() string {
	return "mysql"
}
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
