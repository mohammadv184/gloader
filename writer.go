package gloader

import (
	"fmt"
	"gloader/data"
	"gloader/driver"
	"log"
	"sync"
)

type Writer struct {
	buffer         *data.Buffer
	connectionP    *driver.ConnectionPool
	dataCollection string
	workers        uint
	rowPerBatch    uint64
}

func NewWriter(dataCollection string, buffer *data.Buffer, connectionP *driver.ConnectionPool) *Writer {
	return &Writer{
		buffer:         buffer,
		connectionP:    connectionP,
		dataCollection: dataCollection,
		workers:        DefaultWorkers,
		rowPerBatch:    DefaultRowsPerBatch,
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

	wg := sync.WaitGroup{}
	wg.Add(int(w.workers))
	for i := uint(0); i < w.workers; i++ {
		c, _, err := w.connectionP.Connect()
		if err != nil {
			return err
		}

		conn := c.(driver.WritableConnection)

		go func(c driver.WritableConnection) {
			defer wg.Done()
			defer c.Close()
			for {
				batch := data.NewDataBatch()

				for i := uint64(0); i < w.rowPerBatch; i++ {
					dSet, err := w.buffer.Read()
					if err != nil {
						log.Println(err)
						return
					}
					if dSet == nil {
						if batch.GetLength() > 0 {
							fmt.Println("write")
							if err := c.Write(w.dataCollection, batch); err != nil {
								log.Println(err)
							}
						}

						return
					}
					// fmt.Println("Length: ", w.dataCollection, batch.GetLength())
					batch.Add(dSet)
				}
				rt := 10 // retry times
				for i := 0; i < rt; i++ {
					if err := c.Write(w.dataCollection, batch); err != nil {
						log.Println(err)
						continue
					}
					break
				}

			}
		}(conn)
	}
	wg.Wait()
	return nil
}
