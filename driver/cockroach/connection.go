package cockroach

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

// Connection is a connection to a CockroachDB database.
type Connection struct {
	conn         *pgx.Conn
	isClosed     bool
	dbName       string
	tableDetails map[string]driver.DataCollectionDetail
	config       *Config

	driver.DefaultFilterBuilder
	driver.DefaultSortBuilder
}

// Close closes the connection to the database.
func (c *Connection) Close() error {
	if c.isClosed {
		return driver.ErrConnectionIsClosed
	}

	err := c.conn.Close(context.Background())
	if err != nil {
		return err
	}
	c.isClosed = true
	return err
}

func (c *Connection) Ping() error {
	if c.isClosed {
		return driver.ErrConnectionIsClosed
	}
	return c.conn.Ping(context.Background())
}

// IsClosed returns the status of the connection.
func (c *Connection) IsClosed() bool {
	return c.isClosed
}

// GetDetails returns the details of the database.
func (c *Connection) GetDetails(ctx context.Context) (driver.DatabaseDetail, error) {
	if c.isClosed {
		return driver.DatabaseDetail{}, driver.ErrConnectionIsClosed
	}

	databaseInfo := driver.DatabaseDetail{
		Name:            c.dbName,
		DataCollections: make([]driver.DataCollectionDetail, 0),
	}

	tables, err := c.conn.Query(ctx, "SHOW TABLES")
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
		dc, err := c.getTableDetails(ctx, table.Name)
		if err != nil {
			return driver.DatabaseDetail{}, err
		}
		databaseInfo.DataCollections[i] = dc
	}

	return databaseInfo, nil
}

// Write writes a batch of data to the database.
func (c *Connection) Write(ctx context.Context, table string, dataBatch *data.Batch) error {
	if c.isClosed {
		return driver.ErrConnectionIsClosed
	}

	if dataBatch.GetLength() == 0 {
		return nil
	}

	tDetails, err := c.getTableDetails(ctx, table)
	if err != nil {
		return fmt.Errorf("cockroach: failed to get table details: %w", err)
	}

	rows := make(map[string][][]any)
	rowsKeys := make(map[string][]string)

	for _, dataSet := range *dataBatch {
		row := make([]any, 0)
		rowKeys := make([]string, 0)

		for i, tKey := range tDetails.DataMap.Keys() {
			t := tDetails.GetDataMap().GetIndex(i)
			if !dataSet.Has(tKey) || dataSet.Get(tKey).GetValueType().GetValue() == nil {
				if tDetails.DataMap.HasDefaultValue(tKey) {
					continue
				}

				if tDetails.DataMap.IsNullable(tKey) {
					row = append(row, nil)
					rowKeys = append(rowKeys, tKey)
					continue
				}

				return fmt.Errorf("cockroach: column %s is not nullable or does not have a default value but is not set", tKey)
			}
			d := dataSet.Get(tKey)

			if !t.GetTypeKind().IsCompatibleWith(d.GetValueType().GetTypeKind()) {
				return fmt.Errorf(
					"cockroach: column %s data type kind mismatch ( %s != %s )",
					tKey,
					t.GetTypeKind().String(),
					d.GetValueType().GetTypeKind().String(),
				)
			}

			dValueType, err := d.GetValueType().To(t)
			if err != nil {
				return fmt.Errorf("cockroach: failed to convert data type: %w", err)
			}
			row = append(row, dValueType.GetValue())
			rowKeys = append(rowKeys, tKey)
		}
		rKeysCopy := make([]string, len(rowKeys))
		copy(rKeysCopy, rowKeys)
		sort.Slice(rKeysCopy, func(i, j int) bool { return rKeysCopy[i] < rKeysCopy[j] })
		k := strings.Join(rKeysCopy, ",")
		rows[k] = append(rows[k], row)
		if _, ok := rowsKeys[k]; !ok {
			rowsKeys[k] = rowKeys
		}
	}

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// insert rows with default values
	for k, rows := range rows {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{table},
			rowsKeys[k],
			pgx.CopyFromRows(rows),
		)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" { // 23505 is the unique_violation error code
					fmt.Println("Unique violation detected: ", pgErr.Detail, table)
					// TODO: handle unique violation
					//dupRow := regexp.MustCompile(`\((.*)\)=\((.*)\)`).FindStringSubmatch(pgErr.Detail)
					//if len(dupRow) != 3 {
					//	fmt.Println("Error parsing duplicate row")
					//	return err
					//}
					//
					//for i, set := range *dataBatch {
					//	if set.Has(dupRow[1]) && set.Get(dupRow[1]).GetValueType().GetValue() == dupRow[2] {
					//
					//	}
					//}
				}
				fmt.Println("Error executing statement: ", pgErr.Code, pgErr.Message, pgErr.Detail, table)
			}
			fmt.Println("Error executing statement")
			return err
		}

	}

	return tx.Commit(ctx)
}

func (c *Connection) getTableDetails(ctx context.Context, table string) (driver.DataCollectionDetail, error) {
	if c.tableDetails == nil {
		c.tableDetails = make(map[string]driver.DataCollectionDetail)
	}

	if t, isExists := c.tableDetails[table]; isExists {
		return t, nil
	}

	dc := driver.DataCollectionDetail{
		Name:    table,
		DataMap: new(data.Map),
	}

	columns, err := c.conn.Query(ctx, fmt.Sprintf("SHOW COLUMNS FROM %s", table))
	if err != nil {
		return driver.DataCollectionDetail{}, err
	}

	// the client will automatically close the rows when all the rows are read
	for columns.Next() {
		var columnName, columnType string
		var columnDefault []byte
		var columnNullable bool
		err = columns.Scan(&columnName, &columnType, &columnNullable, &columnDefault, nil, nil, nil)
		if err != nil {
			return driver.DataCollectionDetail{}, err
		}
		t, err := GetTypeFromName(columnType)
		if err != nil {
			return driver.DataCollectionDetail{}, err
		}
		dc.DataMap.Set(columnName, t, columnNullable, len(columnDefault) != 0)
	}
	if err = columns.Err(); err != nil {
		return driver.DataCollectionDetail{}, err
	}

	var count int

	res := c.conn.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, c.BuildFilterSQL(table)))

	err = res.Scan(&count)
	if err != nil {
		return driver.DataCollectionDetail{}, err
	}
	dc.DataSetCount = count

	c.tableDetails[table] = dc

	return dc, nil
}
