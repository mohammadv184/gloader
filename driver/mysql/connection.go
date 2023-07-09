package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

// Connection is a connection to a MySQL database.
type Connection struct {
	conn     *sql.DB
	isClosed bool
	config   *Config
	driver.DefaultFilterBuilder
	driver.DefaultSortBuilder
}

// Close closes the connection to the database.
func (m *Connection) Close() error {
	if m.isClosed {
		return driver.ErrConnectionIsClosed
	}

	err := m.conn.Close()
	if err != nil {
		return err
	}
	m.isClosed = true
	return err
}

// IsClosed returns the status of the connection.
func (m *Connection) IsClosed() bool {
	return m.isClosed
}

// Ping pings the database.
func (m *Connection) Ping() error {
	if m.isClosed {
		return driver.ErrConnectionIsClosed
	}
	return m.conn.Ping()
}

// GetDetails returns the details of the database.
func (m *Connection) GetDetails(ctx context.Context) (driver.DatabaseDetail, error) {
	if m.isClosed {
		return driver.DatabaseDetail{}, driver.ErrConnectionIsClosed
	}
	databaseInfo := driver.DatabaseDetail{
		Name:            m.config.Database,
		DataCollections: make([]driver.DataCollectionDetail, 0),
	}

	tables, err := m.conn.QueryContext(ctx, "SHOW TABLES")
	if err != nil {
		return driver.DatabaseDetail{}, err
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		err = tables.Scan(&tableName)

		if err != nil {
			return driver.DatabaseDetail{}, err
		}

		databaseInfo.DataCollections = append(databaseInfo.DataCollections, driver.DataCollectionDetail{
			Name:         tableName,
			DataMap:      new(data.Map),
			DataSetCount: 0,
		})
	}

	for i, table := range databaseInfo.DataCollections {
		columns, err := m.conn.QueryContext(ctx, "SHOW COLUMNS FROM ?", table.Name)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}

		for columns.Next() {
			var columnName, columnType string
			var columnNullable bool
			var null any
			err = columns.Scan(&columnName, &columnType, &columnNullable, &null, &null, &null)
			if err != nil {
				return driver.DatabaseDetail{}, err
			}

			t, err := GetTypeFromName(columnType)
			if err != nil {
				return driver.DatabaseDetail{}, err
			}

			databaseInfo.DataCollections[i].DataMap.Set(columnName, t, columnNullable)
		}

		rows, err := m.conn.QueryContext(ctx, "SELECT COUNT(*) FROM ? "+m.BuildFilterSQL(table.Name), table.Name)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}

		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				return driver.DatabaseDetail{}, err
			}
			databaseInfo.DataCollections[i].DataSetCount = count
		}
		rows.Close()
		columns.Close()
	}
	return databaseInfo, nil
}

// Read reads data from the database.
func (m *Connection) Read(ctx context.Context, dataCollection string, startOffset, endOffset uint64) (*data.Batch, error) {
	if m.isClosed {
		return nil, driver.ErrConnectionIsClosed
	}
	fmt.Println("Reading from", startOffset, "to", endOffset)

	batch := data.NewDataBatch()

	rows, err := m.conn.QueryContext(
		ctx,
		"SELECT * FROM "+
			dataCollection+
			m.BuildFilterSQL(dataCollection)+
			m.BuildSortSQL(dataCollection)+
			" LIMIT "+
			fmt.Sprint(startOffset)+
			", "+
			fmt.Sprint(endOffset-startOffset),
	)
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
