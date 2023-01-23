package cockroach

import (
	"database/sql"
	"fmt"
	"gloader/data"
	"gloader/driver"
	"strings"

	"github.com/lib/pq"
)

type Connection struct {
	conn   *sql.DB
	dbName string
}

func (m *Connection) Close() error {
	return m.conn.Close()
}

func (m *Connection) GetDetails() (*driver.DataBaseDetails, error) {
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
		columns, err := m.conn.Query("SHOW COLUMNS FROM " + table.Name)
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
func (m *Connection) Write(table string, dataBatch *data.Batch) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn(table, dataBatch.Get(0).GetKeys()...))
	if err != nil {
		return err
	}

	for _, dataSet := range *dataBatch {
		values := make([]interface{}, dataSet.GetLength())
		for i, key := range strings.Split(dataSet.String(", "), ", ") {
			values[i] = key
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		for _, dataSet := range *dataBatch {
			fmt.Println(dataSet.String(", "))
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
