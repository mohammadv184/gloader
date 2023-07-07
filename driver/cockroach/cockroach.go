package cockroach

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/mohammadv184/gloader/driver"
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
func (*Cockroach) GetDriverName() string {
	return "cockroach"
}

func (*Cockroach) IsReadable() bool {
	return false
}

func (*Cockroach) IsWritable() bool {
	return true
}

// Open opens a connection to the database.
func (*Cockroach) Open(ctx context.Context, dsn string) (driver.Connection, error) {
	config, err := parseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(ctx, config.String())
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Connection{conn: conn, config: config}, nil
}
