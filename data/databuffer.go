package data

import (
	"errors"
	"fmt"
	"sync"
)

// DefaultMaxBufferSize is the default maximum size of the buffer in bytes.
const (
	DefaultMaxBufferSize = 1024 * 1024 * 256 // 256 MB
	// DefaultMaxBufferLength is the default maximum length of the buffer.
	DefaultMaxBufferLength = 500000
)

// Buffer is a thread-safe data buffer. It is used to buffer data sets before they are written to the database.
type Buffer struct {
	data *Batch

	locker *sync.RWMutex

	maxSize uint64

	maxLength uint64

	close chan any
}

// NewBuffer creates a new buffer with the given maximum size in bytes.
// If no size is given, the DefaultMaxBufferLength is used.
// size: the maximum size of the buffer in bytes. and cannot be 0.
func NewBuffer(size ...uint64) *Buffer {
	var bSize uint64 = DefaultMaxBufferSize
	if len(size) > 0 && size[0] > 0 {
		bSize = size[0]
	}

	return &Buffer{
		data:      NewDataBatch(),
		maxSize:   bSize,
		close:     make(chan any),
		locker:    &sync.RWMutex{},
		maxLength: DefaultMaxBufferLength,
	}
}

// Write writes the given data sets to the buffer.
// if the buffer exceeds the maximum conditions, it will be blocked until the buffer conditions are met.
// if the buffer is closed, it will return an error.
func (b *Buffer) Write(data ...*Set) error {
	if b.IsClosed() {
		return ErrBufferIsClosed
	}
	b.checkMaxSize()
	b.locker.Lock()
	defer b.locker.Unlock()
	b.data.Add(data...)
	return nil
}

// Read pops the first data set from the buffer. If the buffer is empty, it will be blocked until the next data set is written to the buffer.
// If the buffer is closed, and all data sets have been read, it will return nil.
func (b *Buffer) Read() (*Set, error) {
	for {
		data, err := b.popDataSet()

		if data != nil || err != nil {
			if errors.Is(err, ErrBufferIsClosed) {
				return nil, nil
			}
			return data, err
		}
	}
}

// popDataSet pops the first data set from the buffer. If the buffer is empty, it will be wait.
// If the buffer is closed, and all data sets have been read, it will return an error.
func (b *Buffer) popDataSet() (*Set, error) {
	for {
		if b.IsClosed() && b.IsEmpty() {
			fmt.Println("buffer is closed")
			return nil, ErrBufferIsClosed
		}
		if b.IsEmpty() {
			continue
		}
		break
	}
	b.locker.Lock()
	defer b.locker.Unlock()
	return b.data.Pop(), nil
}

// Clear clears the buffer.
func (b *Buffer) Clear() {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.data = NewDataBatch()
}

// GetSize returns the current size of the buffer in bytes.
func (b *Buffer) GetSize() uint64 {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.data.GetSize()
}

// GetLength returns the current length of the buffer.
func (b *Buffer) GetLength() int {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.data.GetLength()
}

// Size sets the maximum size of the buffer in bytes.
// If the buffer exceeds the maximum size, it will be blocked until the buffer conditions are met.
func (b *Buffer) Size(size uint64) {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.maxSize = size
}

// Length sets the maximum length of the buffer.
// If the buffer exceeds the maximum length, it will be blocked until the buffer conditions are met.
func (b *Buffer) Length(length uint64) {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.maxLength = length
}

// IsEmpty returns true if the buffer is empty.
func (b *Buffer) IsEmpty() bool {
	return b.data.GetLength() == 0
}

// Clone returns a copy of the buffer.
// The copy will not be affected by the original buffer.
// and Clone always returns a not closed buffer.
func (b *Buffer) Clone() *Buffer {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return &Buffer{
		data:      b.data.Clone(),
		locker:    &sync.RWMutex{},
		maxSize:   b.maxSize,
		close:     make(chan any),
		maxLength: b.maxLength,
	}
}

// IsClosed returns true if the buffer is closed.
func (b *Buffer) IsClosed() bool {
	select {
	case <-b.close:
		return true
	default:
		return false
	}
}

// Close will close the buffer.
// closed buffer will not accept new data sets.
// If the buffer is already closed, it will return an error.
func (b *Buffer) Close() error {
	if b.IsClosed() {
		return ErrBufferAlreadyIsClosed
	}
	close(b.close)
	return nil
}

// checkMaxSize checks the maximum size and length of the buffer.
// If the buffer exceeds the maximum conditions, it will be blocked until the buffer conditions are met.
func (b *Buffer) checkMaxSize() {
	checkConditionsF := func() bool {
		return b.maxSize > b.GetSize() && b.maxLength > uint64(b.GetLength())
	}

	for {
		switch {
		case checkConditionsF():
			return
		}
	}
}
