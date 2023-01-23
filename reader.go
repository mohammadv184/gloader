package gloader

import (
	"errors"
	"fmt"
	"gloader/data"
	"gloader/driver"
	"log"
	"sync"
)

const (
	defaultRowsPerBatch = 50
	defaultWorkers      = 2
)

var (
	ErrBufferNotSet                 = errors.New("buffer not set")
	ErrConnectionPoolNotSet         = errors.New("connection pool not set")
	ErrDataMapNotSet                = errors.New("data map not set")
	ErrEndOffsetLessThanStartOffset = errors.New("end offset less than start offset")
	ErrEndOffsetRequired            = errors.New("end offset required")
)

type Reader struct {
	connectionP    *driver.ConnectionPool
	buffer         *data.Buffer
	dataCollection string
	dataMap        *data.Map
	rowPerBatch    uint64
	workers        uint
	startOffset    uint64
	endOffset      uint64
}

func NewReader(dataCollection string, buffer *data.Buffer, dataMap *data.Map, connectionP *driver.ConnectionPool) *Reader {
	return &Reader{
		connectionP:    connectionP,
		buffer:         buffer,
		dataCollection: dataCollection,
		dataMap:        dataMap,
		rowPerBatch:    defaultRowsPerBatch,
		workers:        defaultWorkers,
	}
}
func (r *Reader) SetRowsPerBatch(rowsPerBatch uint64) {
	r.rowPerBatch = rowsPerBatch
}
func (r *Reader) SetWorkers(workers uint) {
	r.workers = workers
}
func (r *Reader) SetStartOffset(startOffset uint64) {
	r.startOffset = startOffset
}
func (r *Reader) SetEndOffset(endOffset uint64) {
	r.endOffset = endOffset
}

func (r *Reader) Start() error {
	if r.buffer == nil {
		return ErrBufferNotSet
	}
	if r.connectionP == nil {
		return ErrConnectionPoolNotSet
	}
	if r.dataMap == nil {
		return ErrDataMapNotSet
	}
	if r.endOffset == 0 {
		return ErrEndOffsetRequired
	}

	if r.endOffset < r.startOffset {
		return ErrEndOffsetLessThanStartOffset
	}

	readChs := make([]<-chan *data.Batch, r.workers)

	for i := uint(0); i < r.workers; i++ {
		startOffset := r.startOffset + uint64(float64(i)*float64(r.endOffset-r.startOffset)/float64(r.workers))

		endOffset := r.startOffset + uint64(float64(i+1)*float64(r.endOffset-r.startOffset)/float64(r.workers))
		fmt.Println(startOffset, endOffset)
		conn, _, err := r.connectionP.Connect()
		if err != nil {
			return err
		}

		sConn := conn.(driver.ReadableConnection)
		readChs[i] = sConn.StartReader(r.dataCollection, *r.dataMap, startOffset, endOffset)
	}
	fmt.Println(len(readChs))

	r.fanIn(readChs)
	fmt.Println("closing buffer")
	err := r.buffer.Close()
	if err != nil {
		return err
	}

	return r.connectionP.Close()
}

func (r *Reader) fanIn(readChs []<-chan *data.Batch) {
	wg := sync.WaitGroup{}
	wg.Add(len(readChs))
	//fmt.Println(len(readChs))
	for _, readCh := range readChs {
		go func(readCh <-chan *data.Batch) {
			for {
				select {
				case batch, ok := <-readCh:
					fmt.Println("readCh recived a data batch")
					if !ok {
						fmt.Println("batch readCh closed", r.buffer.GetLength())
						wg.Done()
						return
					}
					err := r.buffer.Write(*batch...)
					if err != nil {
						log.Fatal(err)
						return
					}
					fmt.Println(r.buffer.GetLength())
				}
			}
		}(readCh)
	}
	wg.Wait()
}
