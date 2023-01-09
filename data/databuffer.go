package data

import "sync"

type Buffer struct {
	data *Batch

	locker sync.RWMutex

	maxSize int

	sizeLocker sync.Mutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		data: NewDataBatch(),
	}
}

func (b *Buffer) Write(data *Set) {
	b.checkMaxSize()
	b.locker.Lock()
	defer b.locker.Unlock()
	b.data.Add(data)
}

func (b *Buffer) Read() *Set {
	b.locker.RLock()
	defer b.locker.RUnlock()
	if b.data.GetSize() == 0 {
		return nil
	}
	data := b.data.Pop()
	return data
}

func (b *Buffer) Clear() {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.data = NewDataBatch()
}

func (b *Buffer) GetSize() int {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.data.GetSize()
}

func (b *Buffer) IsEmpty() bool {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.data.GetSize() == 0
}

func (b *Buffer) Clone() *Buffer {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return &Buffer{
		data: b.data.Clone(),
	}
}

func (b *Buffer) checkMaxSize() {
	if b.maxSize > b.GetSize() {
		for {
			if b.maxSize < b.GetSize() {
				return
			}
		}
	}
}
