package driver

import (
	"fmt"
	"gloader/data"
)

// Driver is a driver for a database.
type Driver interface {
	// GetDriverName returns the name of the driver.
	GetDriverName() string
	// Open opens a connection to the database.
	Open(dsn string) (Connection, error)
}

// Connection is a connection to a database.
type Connection interface {
	// Close closes the connection.
	Close() error
	// GetDetails returns the details of the database.
	GetDetails() (*DataBaseDetails, error)
}

// WritableConnection is a connection to a database that can write data.
type WritableConnection interface {
	Connection // Embeds Connection
	// Write writes a batch of data to the database.
	Write(dataCollection string, dataBatch *data.Batch) error
}

// ReadableConnection is a connection to a database that can read data.
type ReadableConnection interface {
	Connection // Embeds Connection
	// Read reads a batch of data from the database.
	Read(dataCollection string, startOffset, endOffset uint64) (*data.Batch, error)
}

// DataBaseDetails is the details of a data batch.
type DataBaseDetails struct {
	Name            string
	DataCollections []DataCollectionDetails
}

// DataCollectionDetails is the details of a data collection.
type DataCollectionDetails struct {
	Name         string
	DataMap      data.Map
	DataSetCount int
}

// Connector is a connector to a database.
type Connector struct {
	DefaultSortBuilder
	DefaultFilterBuilder
	dsn    string
	driver Driver
}

// Connect connects to the database.
func (c *Connector) Connect() (Connection, error) {
	conn, err := c.driver.Open(c.dsn)
	if err != nil {
		return nil, err
	}

	if fConn, ok := conn.(FilterableConnection); ok {
		for _, filter := range c.GetFilters() {
			fConn.(FilterableConnection).WhereCondition(filter.GetCondition(), filter.GetKey(), filter.GetValue())
		}
	}

	if sConn, ok := conn.(SortableConnection); ok {
		for _, sort := range c.GetSorts() {
			sConn.(SortableConnection).OrderBy(sort.GetKey(), sort.GetDirection())
		}
	}
	return conn, nil
}

// IsWritable returns true if the connection is writable.
func (c *Connector) IsWritable() bool {
	conn, err := c.Connect()
	if err != nil {
		panic(err)
	}
	_, ok := conn.(WritableConnection)
	return ok
}

// IsReadable returns true if the connection is readable.
func (c *Connector) IsReadable() bool {
	conn, err := c.Connect()
	if err != nil {
		panic(err)
	}
	_, ok := conn.(ReadableConnection)
	return ok
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
func (cp *ConnectionPool) Connect() (Connection, uint, error) {
	conn, err := cp.connector.Connect()
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

// Close closes all connections in the pool.
func (cp *ConnectionPool) Close() error {
	for _, conn := range cp.connections {
		if err := conn.Close(); err != nil {
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
	err = conn.Close()
	if err != nil {
		return err
	}
	if index == cp.GetConnectionLength()-1 {
		cp.connections = cp.connections[:index]
	} else {
		cp.connections = append(cp.connections[:index], cp.connections[index+1:]...)
	}
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

// GetConnectionLength returns the length of the connections.
func (cp *ConnectionPool) GetConnectionLength() uint {
	return uint(len(cp.connections))
}

// NewConnectionPool returns a new connection pool.
func NewConnectionPool(connector *Connector) *ConnectionPool {
	return &ConnectionPool{
		connector:   connector,
		connections: make([]Connection, 0),
	}
}
