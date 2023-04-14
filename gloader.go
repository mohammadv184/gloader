package gloader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

const (
	DefaultRowsPerBatch = 100
	DefaultWorkers      = 3
)

var (
	ErrBufferNotSet                 = errors.New("buffer not set")
	ErrConnectionPoolNotSet         = errors.New("connection pool not set")
	ErrDataMapNotSet                = errors.New("data map not set")
	ErrEndOffsetLessThanStartOffset = errors.New("end offset less than start offset")
	ErrEndOffsetRequired            = errors.New("end offset required")
	ErrSrcConnectionIsRequired      = errors.New("source connection is required")
)

var (
	CCCauseStopFuncCalled = errors.New("stop func called")
)

type GLoader struct {
	srcConnector              *driver.Connector
	destConnector             *driver.Connector
	reader                    *Reader
	writer                    *Writer
	dataCollectionEndOffset   map[string]uint64
	dataCollectionStartOffset map[string]uint64
	includedDataCollections   []string
	excludedDataCollections   []string
	rowsPerBatch              uint64
	workers                   uint
	ctx                       context.Context
	ctxCancelFunc             context.CancelCauseFunc
}

func NewGLoader() *GLoader {
	return &GLoader{
		rowsPerBatch:              DefaultRowsPerBatch,
		workers:                   DefaultWorkers,
		dataCollectionEndOffset:   make(map[string]uint64),
		dataCollectionStartOffset: make(map[string]uint64),
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

func (g *GLoader) Filter(dataCollection, key string, condition driver.Condition, value string) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.WhereCondition(dataCollection, condition, key, value)
	return g
}

func (g *GLoader) OrderBy(dataCollection, key string, direction driver.Direction) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.OrderBy(dataCollection, key, direction)
	return g
}
func (g *GLoader) FilterAll(key string, condition driver.Condition, value string) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.WhereRootCondition(condition, key, value)
	return g
}

func (g *GLoader) OrderByAll(key string, direction driver.Direction) *GLoader {
	if g.srcConnector == nil {
		panic(ErrSrcConnectionIsRequired)
	}
	g.srcConnector.OrderByRoot(key, direction)
	return g
}

func (g *GLoader) Include(dataCollections ...string) *GLoader {
	g.includedDataCollections = dataCollections
	return g
}

func (g *GLoader) Exclude(dataCollections ...string) *GLoader {
	g.excludedDataCollections = dataCollections
	return g
}
func (g *GLoader) SetEndOffset(dataCollection string, offset uint64) *GLoader {
	g.dataCollectionEndOffset[dataCollection] = offset
	return g
}
func (g *GLoader) SetStartOffset(dataCollection string, offset uint64) *GLoader {
	g.dataCollectionStartOffset[dataCollection] = offset
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
	return g.StartWithContext(context.Background())
}

func (g *GLoader) StartWithContext(ctx context.Context) error {
	c, cCancelCauseFunc := context.WithCancelCause(ctx)
	g.ctx = c
	g.ctxCancelFunc = cCancelCauseFunc

	srcConn, err := g.srcConnector.Connect(c)
	if err != nil {
		return err
	}
	defer srcConn.Close()

	sDetails, err := srcConn.GetDetails(c)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, dc := range sDetails.DataCollections {
		if dc.DataSetCount == 0 {
			continue
		}

		if len(g.includedDataCollections) > 0 {
			var found bool
			for _, included := range g.includedDataCollections {
				if dc.Name == included {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(g.excludedDataCollections) > 0 {
			var found bool
			for _, excluded := range g.excludedDataCollections {
				if dc.Name == excluded {
					found = true
					break
				}
			}
			if found {
				continue
			}
		}

		wg.Add(2)
		fmt.Println("Starting to load", dc.Name, "from", 0, "to", dc.DataSetCount)
		buffer := data.NewBuffer(c)

		rConnectionPool := driver.NewConnectionPool(g.srcConnector)
		wConnectionPool := driver.NewConnectionPool(g.destConnector)

		reader := NewReader(c, dc.Name, buffer, &dc.DataMap, rConnectionPool)

		if offset, ok := g.dataCollectionStartOffset[dc.Name]; ok {
			reader.SetStartOffset(offset)
		}

		if offset, ok := g.dataCollectionEndOffset[dc.Name]; ok {
			reader.SetEndOffset(offset)
		} else {
			reader.SetEndOffset(uint64(dc.DataSetCount))
		}

		reader.SetRowsPerBatch(g.rowsPerBatch)
		reader.SetWorkers(g.workers)

		writer := NewWriter(c, dc.Name, buffer, wConnectionPool)
		writer.SetRowsPerBatch(g.rowsPerBatch)
		writer.SetWorkers(g.workers)

		go func(reader *Reader, rcPool *driver.ConnectionPool) {
			err := reader.Start()
			if err != nil {
				log.Println(err)
			}
			wg.Done()
			err = rcPool.CloseAll()
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
			err = wcPool.CloseAll()
			if err != nil {
				log.Println(err)
			}
		}(writer, wConnectionPool)

	}
	wg.Wait()
	return nil

}

func (g *GLoader) Stop() {
	g.ctxCancelFunc(CCCauseStopFuncCalled)
}
