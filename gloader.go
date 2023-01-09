package gloader

import (
	"bytes"
	"gloader/driver"
)

type GLoader struct {
	sourceDriver driver.Driver
	destDriver   driver.Driver
	sourceDSN    string
	destDSN      string
}

func NewGLoader(sourceDriver, destDriver driver.Driver, sourceDSN, destDSN string) *GLoader {
	return &GLoader{
		sourceDriver: sourceDriver,
		destDriver:   destDriver,
		sourceDSN:    sourceDSN,
		destDSN:      destDSN,
	}
	bytes.Buffer{}
}
func Start() error {

}
