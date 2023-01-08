package gloader

import (
	"gloader/driver"
)

type GLoader struct {
	source      driver.ReadableDriver
	destination driver.WritableDriver
	sDSN        string
	dDSN        string
}

func New() *GLoader {
	return &GLoader{}
}

func (g *GLoader) Source(source string, dsn string) error {
	src, err := driver.GetDriver(source)
	if err != nil {
		return err
	}
	if !driver.IsReadableConnection(src) {
		return driver.ErrDriverNotReadable
	}
	g.setSource(src.(driver.ReadableDriver))
	return nil
}

func (g *GLoader) Destination(destination string) error {
	dst, err := driver.GetDriver(destination)
	if err != nil {
		return err
	}
	if !driver.IsWritableConnection(dst) {
		return driver.ErrDriverNotWritable
	}
	g.setDestination(dst.(driver.WritableDriver))
	return nil
}
