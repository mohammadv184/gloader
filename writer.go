package gloader

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/mohammadv184/gloader/data"
	"github.com/mohammadv184/gloader/driver"
)

type Writer struct {
	buffer         *data.Buffer
	connectionP    *driver.ConnectionPool
	dataCollection string
	workers        uint
	rowPerBatch    uint64
	ctx            context.Context
}

func NewWriter(ctx context.Context, dataCollection string, buffer *data.Buffer, connectionP *driver.ConnectionPool) *Writer {
	return &Writer{
		buffer:         buffer,
		connectionP:    connectionP,
		dataCollection: dataCollection,
		workers:        DefaultWorkers,
		rowPerBatch:    DefaultRowsPerBatch,
		ctx:            ctx,
	}
}

func (w *Writer) SetWorkers(workers uint) {
	w.workers = workers
}

func (w *Writer) SetRowsPerBatch(rowsPerBatch uint64) {
	w.rowPerBatch = rowsPerBatch
}

func (w *Writer) Start() error {
	if w.buffer == nil {
		return ErrBufferNotSet
	}
	if w.connectionP == nil {
		return ErrConnectionPoolNotSet
	}

	wg := &sync.WaitGroup{}
	wg.Add(int(w.workers))

	for i := uint(0); i < w.workers; i++ {
		go func() {
			defer wg.Done()
			w.RunWorker()
		}()
	}

	wg.Wait()
	return nil
}

func (w *Writer) RunWorker() {
	conn, cIndex, err := w.connectionP.Connect(w.ctx)
	if err != nil {
		panic(err)
	}
	wConn := conn.(driver.WritableConnection)

	batch := data.NewDataBatch()
	for {
		batch.Clear()

		for i := uint64(0); i < w.rowPerBatch; i++ {
			dSet, err := w.buffer.Read()
			if err != nil {
				if errors.Is(err, data.ErrBufferIsClosed) {
					batch.Add(dSet)
					break
				}
				log.Println(err)
				return
			}

			// fmt.Println("Length: ", w.dataCollection, batch.GetLength())
			batch.Add(dSet)
		}

		if batch.GetLength() > 0 {
			if err := wConn.Write(w.ctx, w.dataCollection, batch); err != nil {
				log.Printf("error on writing data to %s: %s", w.dataCollection, err)
				panic(err)
			}
		}

		if w.buffer.IsClosed() && w.buffer.IsEmpty() {
			log.Printf("%s Writer worker is closed", w.dataCollection)
			err := w.connectionP.CloseConnection(cIndex)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
}
