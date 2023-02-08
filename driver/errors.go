package driver

import (
	"errors"
)

var ErrDriverNotFound = errors.New("driver not found")

var (
	ErrConnectionNotReadable = errors.New("connection is not readable")
	ErrConnectionNotWritable = errors.New("connection not writable")
)

var ErrConnectionPoolOutOfIndex = errors.New("connection pool out of index")
