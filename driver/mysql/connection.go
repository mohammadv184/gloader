package mysql

import (
	"database/sql"
	"fmt"
	"gloader/data"
	"gloader/driver"
	"reflect"
)

type Connection struct {
	driver.DefaultFilterBuilder
	driver.DefaultSortBuilder
	conn   *sql.DB
	config *Config
}

func (m *Connection) Close() error {
	return m.conn.Close()
}
func (m *Connection) GetDetails() (*driver.DataBaseDetails, error) {
	databaseInfo := &driver.DataBaseDetails{
		Name:            m.config.Database,
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
			var null any
			err = columns.Scan(&columnName, &columnType, &null, &null, &null, &null)
			if err != nil {
				return nil, err
			}
			t, err := GetTypeFromName(columnType)
			if err != nil {
				return nil, err
			}

			databaseInfo.DataCollections[i].DataMap[columnName] = t
		}

		rows, err := m.conn.Query("SELECT COUNT(*) FROM " + table.Name + m.BuildFilterSQL())
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				return nil, err
			}
			databaseInfo.DataCollections[i].DataSetCount = count
		}
		rows.Close()
		columns.Close()
	}
	return databaseInfo, nil
}

func (m *Connection) StartReader(dataCollection string, dataMap data.Map, startOffset, endOffset, rowPerBatch uint64) <-chan *data.Batch {
	readerCh := make(chan *data.Batch)

	go func() {
		defer close(readerCh)

		for i := startOffset; i < endOffset; i += rowPerBatch {
			fmt.Println("Reading from", i, "to", i+rowPerBatch)
			if i+uint64(rowPerBatch) > endOffset {
				rowPerBatch = endOffset - i
				if rowPerBatch == 0 {
					break
				}
			}

			batch := data.NewDataBatch()
			rows, err := m.conn.Query("SELECT * FROM " + dataCollection + m.BuildFilterSQL() + m.BuildSortSQL() + " LIMIT " + fmt.Sprint(i) + ", " + fmt.Sprint(rowPerBatch))

			if err != nil {
				panic(err)
			}

			columnNames, err := rows.Columns()
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
				for i, r := range row {
					dataType := reflect.New(reflect.TypeOf(dataMap[columnNames[i]]).Elem()).Interface().(data.ValueType)

					err := dataType.Parse(r)
					if err != nil {
						panic(err)
					}

					rowData.Set(columnNames[i], dataType)
				}
				batch.Add(rowData)
			}
			readerCh <- batch
			rows.Close()
		}
	}()

	return readerCh
}
