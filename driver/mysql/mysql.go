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
func (m *MySQL) Open(dsn string) (driver.Connection, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Connection{conn: conn}, nil
}
