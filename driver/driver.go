package driver

import "gloader/data"

type Driver interface {
	GetDriverName() string
	Open(dsn string) (Connection, error)
}

type Connection interface {
	Close() error

	GetDataBaseDetails() (DataBaseDetails, error)
}

type WritableConnection interface {
	Connection
	Write(table string, dataBatch data.Batch) error
}

type ReadableConnection interface {
	Connection
	StartReader(dataCollection string, dataMap map[string]data.Type, startOffset, endOffset uint64) <-chan *data.Batch
}

type DataBaseDetails struct {
	Name            string
	DataCollections []DataCollectionDetails
}

type DataCollectionDetails struct {
	Name         string
	DataMap      map[string]data.Type
	DataSetCount int
}
