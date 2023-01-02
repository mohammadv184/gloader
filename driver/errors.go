package driver

import (
	"errors"
)

var (
	ErrDriverNotFound    = errors.New("driver not found")
	ErrDriverNotReadable = errors.New("driver is not readable")
	ErrDriverNotWritable = errors.New("driver is not writable")
)
