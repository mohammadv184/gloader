package driver

type Driver interface {
	GetDriverName() string
	Open() error
}

type WritableDriver interface {
	Driver
	Write(database, table string, row []byte) error
}

type ReadableDriver interface {
	Driver
	StartReader(database, table string, workers int) []<-chan map[string][]byte
}
