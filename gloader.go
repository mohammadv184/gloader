package gloader

import (
	"gloader/driver"
	"reflect"
)

type GLoader struct {
	source      driver.ReadableDriver
	destination driver.WritableDriver
}

func New() *GLoader {
	return &GLoader{}
}

func (g *GLoader) Source(source string) error {
	src, err := driver.GetDriver(source)
	if err != nil {
		return err
	}
	if !driver.IsReadableDriver(src) {
		return driver.ErrDriverNotReadable
	}
	g.setSource(src.(driver.ReadableDriver))
	return nil
}
func (g *GLoader) setSource(source driver.ReadableDriver) {
	g.source = source
}

func (g *GLoader) Destination(destination string) error {
	dst, err := driver.GetDriver(destination)
	if err != nil {
		return err
	}
	if !driver.IsWritableDriver(dst) {
		return driver.ErrDriverNotWritable
	}
	g.setDestination(dst.(driver.WritableDriver))
	return nil
}
func (g *GLoader) setDestination(destination driver.WritableDriver) {
	g.destination = destination
	reflect.FuncOf().
}
