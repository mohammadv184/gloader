package mysql

import (
	"database/sql"
	"fmt"
	"gloader/data"
	"gloader/driver"
)

type Connection struct {
	driver.DefaultFilterBuilder
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

		rows, err := m.conn.Query("SELECT COUNT(*) FROM " + table.Name)
		if err != nil {
			return databaseInfo, err
		}

		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				return databaseInfo, err
			}
			databaseInfo.DataCollections[i].DataSetCount = count
		}
		rows.Close()
		columns.Close()
	}
	return databaseInfo, nil
}

func (m *Connection) StartReader(dataCollection string, dataMap map[string]data.Type, startOffset, endOffset uint64) <-chan *data.Batch {
	readerCh := make(chan *data.Batch)
	// TODO: rowPerBatch should be configurable dynamically
	rowPerBatch := 50

	go func() {
		defer close(readerCh)

		for i := startOffset; i <= endOffset; i += uint64(rowPerBatch) {
			batch := data.NewDataBatch()
			rows, err := m.conn.Query("SELECT * FROM " + dataCollection + " LIMIT " + fmt.Sprintf("%d", i) + ", " + fmt.Sprintf("%d", rowPerBatch) + m.FiltersToSQL())

			if err != nil {
				panic(err)
			}

			for rows.Next() {
				row := make([]interface{}, len(dataMap))
				for i := range row {
					row[i] = new(interface{})
				}
				err = rows.Scan(row...)
				if err != nil {
					panic(err)
				}

				rowData := data.NewDataSet()
				i := 0
				for columnName, dataType := range dataMap {
					dType := dataType
					err := dType.Parse(row[i])
					if err != nil {
						panic(err)
					}

					rowData.Set(columnName, dType)
				}
				batch.Add(rowData)
				rows.Close()
			}
			readerCh <- batch
		}
	}()

	return readerCh
}
