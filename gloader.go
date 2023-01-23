package gloader

import (
	"gloader/data"
	"gloader/driver"
	"log"
	"sync"
)

type GLoader struct {
	srcConnector  *driver.Connector
	destConnector *driver.Connector
}

func NewGLoader() *GLoader {
	return &GLoader{}
}

func (g *GLoader) Source(name, dsn string) *GLoader {
	d, _ := driver.GetDriver(name)

	g.srcConnector = driver.NewConnector(d, dsn)
	return g
}

func (g *GLoader) Dest(name, dsn string) *GLoader {
	d, _ := driver.GetDriver(name)

	g.destConnector = driver.NewConnector(d, dsn)
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

		buffer := data.NewBuffer()

		rConnectionPool := driver.NewConnectionPool(g.srcConnector)
		wConnectionPool := driver.NewConnectionPool(g.destConnector)

		reader := NewReader(dc.Name, buffer, &dc.DataMap, rConnectionPool)
		reader.SetEndOffset(uint64(dc.DataSetCount))

		writer := NewWriter(dc.Name, buffer, wConnectionPool)

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
