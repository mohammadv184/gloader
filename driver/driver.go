package driver

import (
	"context"
	"fmt"

	"github.com/mohammadv184/gloader/data"
)

// Driver is a driver for a database.
type Driver interface {
	// GetDriverName returns the name of the driver.
	GetDriverName() string
	// IsWritable returns true if the driver is writable.
	IsWritable() bool
	// IsReadable returns true if the driver is readable.
	IsReadable() bool
	// Open opens a connection to the database.
	Open(ctx context.Context, dsn string) (Connection, error)
}

// Connection is a connection to a database.
type Connection interface {
	// Close closes the connection.
	Close() error
	// GetDetails returns the details of the database.
	GetDetails(ctx context.Context) (*DataBaseDetails, error)
}

// WritableConnection is a connection to a database that can write data.
type WritableConnection interface {
	Connection // Embeds Connection
	// Write writes a batch of data to the database.
	Write(ctx context.Context, dataCollection string, dataBatch *data.Batch) error
}

// ReadableConnection is a connection to a database that can read data.
type ReadableConnection interface {
	Connection // Embeds Connection
	// Read reads a batch of data from the database.
	Read(ctx context.Context, dataCollection string, startOffset, endOffset uint64) (*data.Batch, error)
}

// DataBaseDetails is the details of a data batch.
type DataBaseDetails struct {
	Name            string
	DataCollections []DataCollectionDetails
}

// DataCollectionDetails is the details of a data collection.
type DataCollectionDetails struct {
	DataMap      data.Map
	Name         string
	DataSetCount int
}

// Connector is a connector to a database.
type Connector struct {
	driver Driver
	dsn    string
	DefaultSortBuilder
	DefaultFilterBuilder
}

// Connect connects to the database.
func (c *Connector) Connect(ctx context.Context) (Connection, error) {
	conn, err := c.driver.Open(ctx, c.dsn)
	if err != nil {
		return nil, err
	}

	if fConn, ok := conn.(FilterableConnection); ok {
		for dc, filters := range c.GetAllFilters() {
			for _, filter := range filters {
				fConn.(FilterableConnection).WhereCondition(dc, filter.GetCondition(), filter.GetKey(), filter.GetValue())
			}
		}

		for _, filter := range c.GetRootFilters() {
			fConn.(FilterableConnection).WhereRootCondition(filter.GetCondition(), filter.GetKey(), filter.GetValue())
		}
	}

	if sConn, ok := conn.(SortableConnection); ok {
		for dc, sorts := range c.GetAllSorts() {
			for _, sort := range sorts {
				sConn.(SortableConnection).OrderBy(dc, sort.GetKey(), sort.GetDirection())
			}
		}

		for _, sort := range c.GetRootSorts() {
			sConn.(SortableConnection).OrderByRoot(sort.GetKey(), sort.GetDirection())
		}
	}
	return conn, nil
}

// IsWritable returns true if the connection is writable.
func (c *Connector) IsWritable() bool {
	return c.driver.IsWritable()
}

// IsReadable returns true if the connection is readable.
func (c *Connector) IsReadable() bool {
	return c.driver.IsReadable()
}

// GetDriver returns the driver.
func (c *Connector) GetDriver() Driver {
	return c.driver
}

// NewConnector returns a new connector.
func NewConnector(driver Driver, dsn string) *Connector {
	return &Connector{
		dsn:    dsn,
		driver: driver,
	}
}

// ConnectionPool is a pool of connections.
type ConnectionPool struct {
	connector   *Connector
	connections []Connection
}

// Connect connects to the database.
func (cp *ConnectionPool) Connect(ctx context.Context) (Connection, uint, error) {
	conn, err := cp.connector.Connect(ctx)
	if err != nil {
		return nil, 0, err
	}

	cp.connections = append(cp.connections, conn)
	return conn, uint(len(cp.connections) - 1), nil
}

// GetConnection returns a connection from the pool.
func (cp *ConnectionPool) GetConnection(index uint) (Connection, error) {
	if index >= uint(len(cp.connections)) {
		return nil, fmt.Errorf("%v: connection index [%d], with pool size: %d", ErrConnectionPoolOutOfIndex, index, len(cp.connections))
	}
	return cp.connections[index], nil
}

// CloseAll closes all connections in the pool.
func (cp *ConnectionPool) CloseAll() error {
	for i, c := range cp.connections {
		if c == nil {
			continue
		}

		if err := cp.CloseConnection(uint(i)); err != nil {
			return err
		}
	}
	return nil
}

// CloseConnection closes a connection in the pool.
func (cp *ConnectionPool) CloseConnection(index uint) error {
	conn, err := cp.GetConnection(index)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}
	err = conn.Close()
	if err != nil {
		return err
	}

	cp.connections[index] = nil
	return nil
}

// GetConnector returns the connector.
func (cp *ConnectionPool) GetConnector() *Connector {
	return cp.connector
}

// GetConnections returns the connections.
func (cp *ConnectionPool) GetConnections() []Connection {
	return cp.connections
}

// GetConnectionsLength returns the length of the connections.
func (cp *ConnectionPool) GetConnectionsLength() uint {
	return uint(len(cp.connections))
}

// NewConnectionPool returns a new connection pool.
func NewConnectionPool(connector *Connector) *ConnectionPool {
	return &ConnectionPool{
		connector:   connector,
		connections: make([]Connection, 0),
	}
}
