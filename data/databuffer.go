package data

import (
	"errors"
	"fmt"
	"sync"
)

// DefaultMaxBufferSize is the default maximum size of the buffer in bytes.
const DefaultMaxBufferSize = 1024 * 1024 * 256 // 256 MB

const DefaultMaxBufferLength = 500000

type Buffer struct {
	data *Batch

	locker *sync.RWMutex

	maxSize uint64

	maxLength uint64

	close chan any
}

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

func (b *Buffer) Clear() {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.data = NewDataBatch()
}

func (b *Buffer) GetSize() uint64 {
	return b.data.GetSize()
}

func (b *Buffer) GetLength() int {
	return b.data.GetLength()
}

func (b *Buffer) IsEmpty() bool {
	return b.data.GetLength() == 0
}

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

func (b *Buffer) IsClosed() bool {
	select {
	case <-b.close:
		fmt.Println("buffer is closed isClosed")
		return true
	default:
		return false
	}
}

func (b *Buffer) Close() error {
	fmt.Println("close buffer called")
	if b.IsClosed() {
		return ErrBufferAlreadyIsClosed
	}
	close(b.close)
	return nil
}

func (b *Buffer) checkMaxSize() {
	if b.maxSize < b.GetSize() {
		for {
			if b.maxSize > b.GetSize() {
				return
			}
			fmt.Println("checkMaxSize loop")
		}
	}

	if b.maxLength < uint64(b.GetLength()) {
		for {
			if b.maxLength > uint64(b.GetLength()) {
				return
			}
			//fmt.Println("checkMaxLength loop")
		}
	}
}
