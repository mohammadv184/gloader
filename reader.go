package gloader

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

type Reader struct {
	connectionP    *driver.ConnectionPool
	buffer         *data.Buffer
	dataMap        *data.Map
	dataCollection string
	rowPerBatch    uint64
	workers        uint
	startOffset    uint64
	endOffset      uint64
	ctx            context.Context
}

func NewReader(ctx context.Context, dataCollection string, buffer *data.Buffer, dataMap *data.Map, connectionP *driver.ConnectionPool) *Reader {
	return &Reader{
		connectionP:    connectionP,
		buffer:         buffer,
		dataCollection: dataCollection,
		dataMap:        dataMap,
		rowPerBatch:    DefaultRowsPerBatch,
		workers:        DefaultWorkers,
		ctx:            ctx,
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

	defer func() {
		fmt.Println("closing buffer")
		err := r.buffer.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(int(r.workers))

	for i := uint(0); i < r.workers; i++ {
		startOffset := r.startOffset + uint64(float64(i)*float64(r.endOffset-r.startOffset)/float64(r.workers))

		endOffset := r.startOffset + uint64(float64(i+1)*float64(r.endOffset-r.startOffset)/float64(r.workers))
		fmt.Println(startOffset, endOffset)

		go func(startOffset, endOffset uint64) {
			defer wg.Done()
			r.RunWorker(startOffset, endOffset, r.rowPerBatch)
		}(startOffset, endOffset)
	}

	wg.Wait()
	return nil
}

func (r *Reader) RunWorker(startOffset, endOffset, rowPerBatch uint64) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	ch := make(chan *data.Batch)
	conn, cIndex, err := r.connectionP.Connect(r.ctx)
	if err != nil {
		panic(err)
	}
	rConn := conn.(driver.ReadableConnection)

	go func() {
		for i := startOffset; i < endOffset; i += rowPerBatch {
			if i+rowPerBatch > endOffset {
				rowPerBatch = endOffset - i
				if rowPerBatch == 0 {
					break
				}
			}
		retryRead:
			batch, err := rConn.Read(r.ctx, r.dataCollection, i, i+rowPerBatch)
			if err != nil {
				log.Printf("error on reading data batch from %d,%d err: %s\n", i, i+rowPerBatch, err)
				goto retryRead

			}
			if batch.GetLength() == 0 {
				continue
			}

			select {
			case <-r.ctx.Done():
				log.Printf("context canceled: raeder worker %s:%d,%d stoped\n", r.dataCollection, startOffset, endOffset)
				goto stopWorker
			case ch <- batch:
				continue
			}

		}
	stopWorker:
		err := r.connectionP.CloseConnection(cIndex)
		if err != nil {
			log.Println("failed on closing database connection err:", err)
		}
		wg.Done()
		close(ch)
	}()

	go func() {
		for {
			select {
			case batch, ok := <-ch:
				if !ok {
					log.Printf("batch readCh closed %s:%d,%d\n", r.dataCollection, startOffset, endOffset)
					wg.Done()
					return
				}
				err := r.buffer.Write(*batch...)
				if err != nil {
					panic(err)
					return
				}
			}
		}
	}()

	wg.Wait()
}
