package gloader

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

const (
	DefaultRowsPerBatch = 100
	DefaultWorkers      = 5
)

var (
	ErrBufferNotSet                 = errors.New("buffer not set")
	ErrConnectionPoolNotSet         = errors.New("connection pool not set")
	ErrDataMapNotSet                = errors.New("data map not set")
	ErrEndOffsetLessThanStartOffset = errors.New("end offset less than start offset")
	ErrEndOffsetRequired            = errors.New("end offset required")
	ErrSrcConnectionIsRequired      = errors.New("source connection is required")
)

type GLoader struct {
	srcConnector  *driver.Connector
	destConnector *driver.Connector
	reader        *Reader
	writer        *Writer
	rowsPerBatch  uint64
	workers       uint
}

func NewGLoader() *GLoader {
	return &GLoader{
		rowsPerBatch: DefaultRowsPerBatch,
		workers:      DefaultWorkers,
	}
}

func (g *GLoader) Source(name, dsn string) error {
	d, err := driver.GetDriver(name)
	if err != nil {
		return err
	}

	dc := driver.NewConnector(d, dsn)

	if !dc.IsReadable() {
		return driver.ErrConnectionNotReadable
	}

	g.srcConnector = dc
	return nil
}

func (g *GLoader) Dest(name, dsn string) error {
	d, err := driver.GetDriver(name)
	if err != nil {
		return err
	}
	dc := driver.NewConnector(d, dsn)

	if !dc.IsWritable() {
		return driver.ErrConnectionNotWritable
	}

	g.destConnector = dc
	return nil
}

func (g *GLoader) Filter(key string, condition driver.Condition, value string) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.WhereCondition(condition, key, value)
	return g
}

func (g *GLoader) OrderBy(key string, direction driver.Direction) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.OrderBy(key, direction)
	return g
}

func (g *GLoader) SetRowsPerBatch(rowsPerBatch uint64) *GLoader {
	g.rowsPerBatch = rowsPerBatch
	return g
}

func (g *GLoader) SetWorkers(workers uint) *GLoader {
	g.workers = workers
	return g
}

func (g *GLoader) Start() error {
	srcConn, err := g.srcConnector.Connect()
	if err != nil {
		return err
	}
	defer srcConn.Close()

	sDetails, err := srcConn.GetDetails()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, dc := range sDetails.DataCollections {
		wg.Add(2)
		fmt.Println("Starting to load", dc.Name, "from", 0, "to", dc.DataSetCount)
		buffer := data.NewBuffer()

		rConnectionPool := driver.NewConnectionPool(g.srcConnector)
		wConnectionPool := driver.NewConnectionPool(g.destConnector)

		reader := NewReader(dc.Name, buffer, &dc.DataMap, rConnectionPool)
		reader.SetEndOffset(uint64(dc.DataSetCount))
		reader.SetRowsPerBatch(g.rowsPerBatch)
		reader.SetWorkers(g.workers)

		writer := NewWriter(dc.Name, buffer, wConnectionPool)
		writer.SetRowsPerBatch(g.rowsPerBatch)
		writer.SetWorkers(g.workers)

		go func(reader *Reader, rcPool *driver.ConnectionPool) {
			err := reader.Start()
			if err != nil {
				log.Println(err)
			}
			wg.Done()
			err = rcPool.Close()
			if err != nil {
				log.Println(err)
			}
		}(reader, rConnectionPool)

		go func(writer *Writer, wcPool *driver.ConnectionPool) {
			err := writer.Start()
			if err != nil {
				log.Println(err)
			}
			wg.Done()
			err = wcPool.Close()
			if err != nil {
				log.Println(err)
			}
		}(writer, wConnectionPool)

	}
	wg.Wait()
	return nil
}
