package gloader

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
	"github.com/mohammadv184/gloader/pkg/stats"
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
	ErrDestConnectionIsRequired     = errors.New("destination connection is required")
)

var ErrCCCauseStopFuncCalled = errors.New("stop func called")

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
	stats                     *stats.Stats
}

func NewGLoader() *GLoader {
	return &GLoader{
		rowsPerBatch:              DefaultRowsPerBatch,
		workers:                   DefaultWorkers,
		dataCollectionEndOffset:   make(map[string]uint64),
		dataCollectionStartOffset: make(map[string]uint64),
		stats:                     NewStats(),
	}
}

func (g *GLoader) Src(name, dsn string) error {
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

func (g *GLoader) GetSrcDetails(ctx context.Context) (driver.DatabaseDetail, error) {
	if g.srcConnector == nil {
		return driver.DatabaseDetail{}, ErrSrcConnectionIsRequired
	}

	conn, err := g.srcConnector.Connect(ctx)
	if err != nil {
		return driver.DatabaseDetail{}, err
	}

	return conn.GetDetails(ctx)
}

func (g *GLoader) GetDestDetails(ctx context.Context) (driver.DatabaseDetail, error) {
	if g.destConnector == nil {
		return driver.DatabaseDetail{}, ErrDestConnectionIsRequired
	}

	conn, err := g.destConnector.Connect(ctx)
	if err != nil {
		return driver.DatabaseDetail{}, err
	}

	return conn.GetDetails(ctx)
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

func (g *GLoader) Stats() *stats.Stats {
	return g.stats
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

	destConn, err := g.destConnector.Connect(c)
	if err != nil {
		return err
	}
	defer destConn.Close()

	dDetails, err := destConn.GetDetails(c)
	if err != nil {
		return err
	}

	var DCs []driver.DataCollectionDetail
	DCs = sDetails.DataCollections
	if len(g.includedDataCollections) > 0 {
		DCs = sDetails.OnlyDataCollections(g.includedDataCollections...)
	}

	if len(g.excludedDataCollections) > 0 {
		DCs = sDetails.AllDataCollectionsExcept(g.excludedDataCollections...)
	}

	wg := sync.WaitGroup{}
	for i, dc := range DCs {
		if dc.DataSetCount == 0 {
			continue
		}

		dDC, err := dDetails.GetDataCollection(dc.Name)
		if err != nil {
			return fmt.Errorf("GLoader: failed to get destination data collection details for %s", dc.Name)
		}

		var srcSchema strings.Builder
		srcSchema.WriteString("Migration Source Schema:\n")
		var destSchema strings.Builder
		destSchema.WriteString("Migration Destination Schema:\n")
		for k, v := range dc.GetDataMap().GetTypeMap() {
			srcSchema.WriteString(fmt.Sprintf("%d \t %s: \t %s \t IS NULLABLE: %t\n", i, k, v.GetTypeName(), dc.GetDataMap().IsNullable(k)))

			dDT := dDC.GetDataMap().Get(k)
			if dDT == nil {
				continue
			}

			destSchema.WriteString(fmt.Sprintf("%d \t %s: \t %s \t IS NULLABLE: %t\n", i, k, dDT.GetTypeName(), dDC.GetDataMap().IsNullable(k)))
		}

		fmt.Println(srcSchema.String())
		fmt.Println(destSchema.String())

		wg.Add(2)

		buffer := data.NewBuffer(c).
			WithObserver(NewBufferObserverAdapter(g.stats, dc.Name))

		rConnectionPool := driver.NewConnectionPool(g.srcConnector)
		wConnectionPool := driver.NewConnectionPool(g.destConnector)

		reader := NewReader(c, dc.Name, buffer, dc.DataMap, rConnectionPool)

		if offset, ok := g.dataCollectionStartOffset[dc.Name]; ok {
			reader.SetStartOffset(offset)
		}

		if offset, ok := g.dataCollectionEndOffset[dc.Name]; ok {
			reader.SetEndOffset(offset)
		} else {
			reader.SetEndOffset(uint64(dc.DataSetCount))
		}
		fmt.Println("Starting to load", dc.Name, "from", reader.startOffset, "to", reader.endOffset)

		reader.SetRowsPerBatch(g.rowsPerBatch)
		reader.SetWorkers(g.workers)

		writer := NewWriter(c, dc.Name, buffer, wConnectionPool)
		writer.SetRowsPerBatch(g.rowsPerBatch)
		writer.SetWorkers(g.workers)

		go func(reader *Reader, rcPool *driver.ConnectionPool) {
			err := reader.Start()
			if err != nil {
				panic(err)
			}
			wg.Done()
			err = rcPool.CloseAll()
			if err != nil {
				// TODO: logging system
				//log.Println(err)
			}
		}(reader, rConnectionPool)

		go func(writer *Writer, wcPool *driver.ConnectionPool) {
			err := writer.Start()
			if err != nil {
				panic(err)
			}
			wg.Done()
			err = wcPool.CloseAll()
			if err != nil {
				// TODO: logging system
				//log.Println(err)
			}
		}(writer, wConnectionPool)

	}
	wg.Wait()
	return nil
}

func (g *GLoader) Stop() {
	g.ctxCancelFunc(ErrCCCauseStopFuncCalled)
}
