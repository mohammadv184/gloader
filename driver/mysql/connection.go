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

func (m *Connection) Read(dataCollection string, startOffset, endOffset uint64) (*data.Batch, error) {
	fmt.Println("Reading from", startOffset, "to", endOffset)

	batch := data.NewDataBatch()
	rows, err := m.conn.Query("SELECT * FROM " + dataCollection + m.BuildFilterSQL() + m.BuildSortSQL() + " LIMIT " + fmt.Sprint(startOffset) + ", " + fmt.Sprint(endOffset-startOffset))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		row := make([]interface{}, len(columns))
		for i := range row {
			row[i] = new(interface{})
		}
		err = rows.Scan(row...)
		if err != nil {
			return nil, err
		}

		rowData := data.NewDataSet()
		for i, c := range row {
			dataType, err := GetTypeFromName(columns[i].DatabaseTypeName())
			if err != nil {
				return nil, err
			}

			vDataType, ok := dataType.(data.ValueType)
			if !ok {
				return nil, fmt.Errorf("Type %s is not a ValueType", reflect.TypeOf(dataType))
			}

			err = vDataType.Parse(c)
			if err != nil {
				return nil, err
			}

			rowData.Set(columns[i].Name(), vDataType)
		}
		batch.Add(rowData)
	}
	return batch, nil
}
