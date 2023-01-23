package mysql

import (
	"database/sql"
	"gloader/driver"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct{}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

var config = Config{
	Host:     "localhost",
	Port:     3306,
	Username: "root",
	Database: "default",
}

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
	dbName := name[strings.LastIndex(name, "/")+1:]
	conn, err := sql.Open("mysql", name)
	if err != nil {
		return nil, err
	}

	return &Connection{conn: conn, dbName: dbName}, nil
}
