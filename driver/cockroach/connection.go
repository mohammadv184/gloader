package cockroach

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"

	"github.com/lib/pq"
)

// Connection is a connection to a CockroachDB database.
type Connection struct {
	conn   *sql.DB
	dbName string
	config *Config
}

// Close closes the connection to the database.
func (m *Connection) Close() error {
	return m.conn.Close()
}

// GetDetails returns the details of the database.
func (m *Connection) GetDetails(_ context.Context) (*driver.DataBaseDetails, error) {
	databaseInfo := driver.DataBaseDetails{
		Name:            m.dbName,
		DataCollections: make([]driver.DataCollectionDetails, 0),
	}

	tables, err := m.conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		err = tables.Scan(&tableName)
		if err != nil {
			return nil, err
		}

		databaseInfo.DataCollections = append(databaseInfo.DataCollections, driver.DataCollectionDetails{
			Name:         tableName,
			DataMap:      make(map[string]data.Type),
			DataSetCount: 0,
		})
	}

	for i, table := range databaseInfo.DataCollections {
		columns, err := m.conn.Query("SHOW COLUMNS FROM $1", table.Name)
		if err != nil {
			return nil, err
		}

		for columns.Next() {
			var columnName, columnType string
			err = columns.Scan(&columnName, &columnType)
			if err != nil {
				return nil, err
			}
			t, err := GetTypeFromName(columnType)
			if err != nil {
				return nil, err
			}
			databaseInfo.DataCollections[i].DataMap[columnName] = t
		}
	}

	return nil, nil
}

// Write writes a batch of data to the database.
func (m *Connection) Write(_ context.Context, table string, dataBatch *data.Batch) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn(table, dataBatch.Get(0).GetKeys()...))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, dataSet := range *dataBatch {
		values := make([]interface{}, dataSet.GetLength())
		for i, key := range dataSet.GetStringValues() {
			values[i] = key
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" { // 23505 is the unique_violation error code
				fmt.Println("Unique violation detected: ", err.Detail, table)
			}
			fmt.Println("Error executing statement: ", err.Code, err.Message, err.Detail, table)
		}
		fmt.Println("Error executing statement")
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return tx.Commit()
}
