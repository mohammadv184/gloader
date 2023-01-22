package driver

import (
	"fmt"
	"gloader/data"
)

type Driver interface {
	GetDriverName() string
	Open(dsn string) (Connection, error)
}

type Connection interface {
	Close() error

	GetDetails() (*DataBaseDetails, error)
}

type WritableConnection interface {
	Connection
	Write(dataCollection string, dataBatch *data.Batch) error
}

type ReadableConnection interface {
	Connection
	StartReader(dataCollection string, dataMap data.Map, startOffset, endOffset uint64) <-chan *data.Batch
}

type DataBaseDetails struct {
	Name            string
	DataCollections []DataCollectionDetails
}

type DataCollectionDetails struct {
	Name         string
	DataMap      data.Map
	DataSetCount int
}

type Connector struct {
	dsn    string
	driver Driver
}

func (c *Connector) Connect() (Connection, error) {
	return c.driver.Open(c.dsn)
}
func (c *Connector) GetDriver() Driver {
	return c.driver
}

func NewConnector(driver Driver, dsn string) *Connector {
	return &Connector{
		dsn:    dsn,
		driver: driver,
	}
}

type ConnectionPool struct {
	connector   *Connector
	connections []Connection
}

func (cp *ConnectionPool) Connect() (Connection, uint, error) {
	conn, err := cp.connector.Connect()
	if err != nil {
		return nil, 0, err
	}

	cp.connections = append(cp.connections, conn)
	return conn, uint(len(cp.connections) - 1), nil
}
func (cp *ConnectionPool) GetConnection(index uint) (Connection, error) {
	if index >= uint(len(cp.connections)) {
		return nil, fmt.Errorf("%v: connection index [%d], with pool size: %d", ErrConnectionPoolOutOfIndex, index, len(cp.connections))
	}
	return cp.connections[index], nil
}

func (cp *ConnectionPool) Close() error {
	for _, conn := range cp.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

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

func (cp *ConnectionPool) GetConnector() *Connector {
	return cp.connector
}
func (cp *ConnectionPool) GetConnections() []Connection {
	return cp.connections
}
func (cp *ConnectionPool) GetConnectionLength() uint {
	return uint(len(cp.connections))
}

func NewConnectionPool(connector *Connector) *ConnectionPool {
	return &ConnectionPool{
		connector:   connector,
		connections: make([]Connection, 0),
	}
}
