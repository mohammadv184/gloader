package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

// Connection is a connection to a MySQL database.
type Connection struct {
	conn   *sql.DB
	config *Config
	driver.DefaultFilterBuilder
	driver.DefaultSortBuilder
}

// Close closes the connection to the database.
func (m *Connection) Close() error {
	return m.conn.Close()
}

// GetDetails returns the details of the database.
func (m *Connection) GetDetails(_ context.Context) (*driver.DataBaseDetails, error) {
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
		columns, err := m.conn.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table.Name))
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
			columnType = regexp.MustCompile("(\\(\\d+\\)).*").ReplaceAllString(columnType, "")

			t, err := GetTypeFromName(columnType)
			if err != nil {
				return nil, err
			}

			databaseInfo.DataCollections[i].DataMap[columnName] = t
		}

		rows, err := m.conn.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table.Name, m.BuildFilterSQL(table.Name)))
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

// Read reads data from the database.
func (m *Connection) Read(_ context.Context, dataCollection string, startOffset, endOffset uint64) (*data.Batch, error) {
	fmt.Println("Reading from", startOffset, "to", endOffset)

	batch := data.NewDataBatch()

	rows, err := m.conn.Query("SELECT * FROM " + dataCollection + m.BuildFilterSQL(dataCollection) + m.BuildSortSQL(dataCollection) + " LIMIT " + fmt.Sprint(startOffset) + ", " + fmt.Sprint(endOffset-startOffset))
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
