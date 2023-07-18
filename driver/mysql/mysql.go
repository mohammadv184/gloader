package mysql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/mohammadv184/gloader/driver"

	// Import the mysql driver.
	_ "github.com/go-sql-driver/mysql"
)

// MySQL is a driver for MySQL.
type MySQL struct {
	connP map[string]*sql.DB
	mu    *sync.Mutex
}

func init() {
	err := driver.Register(&MySQL{
		connP: make(map[string]*sql.DB),
		mu:    &sync.Mutex{},
	})
	if err != nil {
		// TODO: logging system
		//log.Println(err)
	}
}

// GetDriverName returns the name of the driver.
func (m *MySQL) GetDriverName() string {
	return "mysql"
}

func (m *MySQL) IsReadable() bool {
	return true
}

func (m *MySQL) IsWritable() bool {
	return false
}

// Open opens a connection to the database.
func (m *MySQL) Open(ctx context.Context, name string) (driver.Connection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, err := parseConfig(name)
	if err != nil {
		return nil, err
	}

	if c, isExist := m.connP[config.String()]; !isExist || c != nil {
		connP, err := sql.Open("mysql", config.String())
		if err != nil {
			return nil, err
		}
		m.connP[config.String()] = connP
	}

	conn, err := m.connP[config.String()].Conn(ctx)
	if err != nil {
		return nil, err
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}
