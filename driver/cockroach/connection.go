package cockroach

import (
	"database/sql"
	"gloader/data"
	"gloader/driver"

	"github.com/lib/pq"
)

type Connection struct {
	conn   *sql.DB
	dbName string
}

func (m *Connection) Close() error {
	return m.conn.Close()
}

func (m *Connection) GetDataBaseDetails() (driver.DataBaseDetails, error) {
	databaseInfo := driver.DataBaseDetails{
		Name:            m.dbName,
		DataCollections: make([]driver.DataCollectionDetails, 0),
	}

	tables, err := m.conn.Query("SHOW TABLES")
	if err != nil {
		return databaseInfo, err
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		err = tables.Scan(&tableName)
		if err != nil {
			return databaseInfo, err
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
			return databaseInfo, err
		}

		for columns.Next() {
			var columnName, columnType string
			err = columns.Scan(&columnName, &columnType)
			if err != nil {
				return databaseInfo, err
			}
			t, err := GetTypeFromName(columnType)
			if err != nil {
				return databaseInfo, err
			}
			databaseInfo.DataCollections[i].DataMap[columnName] = t
		}
	}

	return databaseInfo, nil
}
func (m *Connection) Write(table string, dataBatch data.Batch) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn(table, dataBatch[0].GetKeys()...))
	if err != nil {
		return err
	}

	for _, dataSet := range dataBatch {
		_, err = stmt.Exec(dataSet.String(string(rune(9))))
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return tx.Commit()
}
