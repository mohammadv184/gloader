package gloader

import (
	"fmt"
	"gloader/data"
	"gloader/driver"
	"log"
	"sync"
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
		rowPerBatch:    DefaultRowsPerBatch,
		workers:        DefaultWorkers,
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
	defer func() {
		fmt.Println("closing buffer")
		err := r.buffer.Close()
		if err != nil {
			log.Println(err)
		}
	}()
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
		ch := make(chan *data.Batch)
		readChs[i] = ch

		go func(readCh chan<- *data.Batch, startOffset, endOffset, rowPerBatch uint64) {
			defer func() {
				fmt.Println("closing readCh")
				close(readCh)
			}()

			for i := startOffset; i < endOffset; i += rowPerBatch {
				if i+rowPerBatch > endOffset {
					rowPerBatch = endOffset - i
					if rowPerBatch == 0 {
						break
					}
				}
			retryRead:
				batch, err := sConn.Read(r.dataCollection, startOffset, endOffset)
				if err != nil {
					log.Println(err)
					goto retryRead
				}
				if batch.GetLength() == 0 {
					continue
				}
				readCh <- batch
			}
		}(ch, startOffset, endOffset, r.rowPerBatch)

	}
	fmt.Println(len(readChs))

	r.fanIn(readChs)

	return r.connectionP.Close()
}

func (r *Reader) fanIn(readChs []<-chan *data.Batch) {
	wg := sync.WaitGroup{}
	wg.Add(len(readChs))
	// fmt.Println(len(readChs))
	for _, readCh := range readChs {
		go func(rCh <-chan *data.Batch) {
			for {
				select {
				case batch, ok := <-rCh:
					fmt.Println("readCh received a data batch")
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
