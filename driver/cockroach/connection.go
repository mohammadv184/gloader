package cockroach

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

// Connection is a connection to a CockroachDB database.
type Connection struct {
	conn     *pgx.Conn
	isClosed bool
	dbName   string
	config   *Config
}

// Close closes the connection to the database.
func (m *Connection) Close() error {
	if m.isClosed {
		return driver.ErrConnectionIsClosed
	}

	err := m.conn.Close(context.Background())
	if err != nil {
		return err
	}
	m.isClosed = true
	return err
}

func (m *Connection) Ping() error {
	if m.isClosed {
		return driver.ErrConnectionIsClosed
	}
	return m.conn.Ping(context.Background())
}

// IsClosed returns the status of the connection.
func (m *Connection) IsClosed() bool {
	return m.isClosed
}

// GetDetails returns the details of the database.
func (m *Connection) GetDetails(ctx context.Context) (driver.DatabaseDetail, error) {
	if m.isClosed {
		return driver.DatabaseDetail{}, driver.ErrConnectionIsClosed
	}

	databaseInfo := driver.DatabaseDetail{
		Name:            m.dbName,
		DataCollections: make([]driver.DataCollectionDetail, 0),
	}

	tables, err := m.conn.Query(ctx, "SHOW TABLES")
	if err != nil {
		return driver.DatabaseDetail{}, err
	}

	// the client will automatically close the rows when all the rows are read
	for tables.Next() {
		var tableName string
		err = tables.Scan(nil, &tableName, nil, nil, nil, nil)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}

		databaseInfo.DataCollections = append(databaseInfo.DataCollections, driver.DataCollectionDetail{
			Name:         tableName,
			DataMap:      new(data.Map),
			DataSetCount: 0,
		})
	}

	if err = tables.Err(); err != nil {
		return driver.DatabaseDetail{}, err
	}

	for i, table := range databaseInfo.DataCollections {
		columns, err := m.conn.Query(ctx, "SHOW COLUMNS FROM $1", table.Name)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}

		// the client will automatically close the rows when all the rows are read
		for columns.Next() {
			var columnName, columnType string
			var columnNullable bool
			err = columns.Scan(&columnName, &columnType, &columnNullable, nil, nil, nil, nil)
			if err != nil {
				return driver.DatabaseDetail{}, err
			}
			t, err := GetTypeFromName(columnType)
			if err != nil {
				return driver.DatabaseDetail{}, err
			}
			databaseInfo.DataCollections[i].DataMap.Set(columnName, t, columnNullable)
		}
		if err = columns.Err(); err != nil {
			return driver.DatabaseDetail{}, err
		}

		var count int

		c := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM $1", table.Name)

		err = c.Scan(&count)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}
		databaseInfo.DataCollections[i].DataSetCount = count
		columns.Close()
	}

	return databaseInfo, nil
}

// Write writes a batch of data to the database.
func (m *Connection) Write(ctx context.Context, table string, dataBatch *data.Batch) error {
	tx, err := m.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows := make([][]any, dataBatch.GetLength())
	for _, dataSet := range *dataBatch {
		row := make([]any, dataSet.GetLength())
		for i, key := range dataSet.GetValues() {
			row[i] = key
		}
		rows = append(rows, row)
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{table},
		dataBatch.Get(0).GetKeys(),
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is the unique_violation error code
				fmt.Println("Unique violation detected: ", pgErr.Detail, table)
			}
			fmt.Println("Error executing statement: ", pgErr.Code, pgErr.Message, pgErr.Detail, table)
		}
		fmt.Println("Error executing statement")
		return err
	}

	return tx.Commit(ctx)
}
