package data

import (
	"bytes"
	"fmt"
)

// Batch is a collection of data sets.
type Batch []*Set

// Add adds the given data sets to the batch.
// It does not if the given set is nil.
// It does not check for duplicate sets.
func (b *Batch) Add(set ...*Set) {
	// filter out nil sets
	for _, s := range set {
		if s == nil {
			continue
		}
		*b = append(*b, s)
	}
}

// Get returns the data set at the given index.
// It returns nil if the index is out of range.
func (b *Batch) Get(index int) *Set {
	if index >= b.GetLength() {
		return nil
	}
	return (*b)[index]
}

// Pop returns the first data set in the batch and removes it from the batch.
// It returns nil if the batch is empty.
func (b *Batch) Pop() *Set {
	if b.GetLength() == 0 {
		return nil
	}

	set := b.Get(0)
	*b = (*b)[1:]
	return set
}

// GetSize returns approximate size of the batch in bytes.
// It returns 0 if the batch is nil.
func (b *Batch) GetSize() uint64 {
	var size uint64
	for _, set := range *b {
		size += set.GetSize()
	}
	return size
}

// GetLength returns the length of the batch.
func (b *Batch) GetLength() int {
	return len(*b)
}

// Remove removes the data set by the given index.
// It does nothing if the index is out of range.
func (b *Batch) Remove(index int) {
	switch {
	case index < 0 || index >= b.GetLength():
		return
	case index == 0:
		*b = (*b)[1:]
	case index == b.GetLength()-1:
		*b = (*b)[:index]
	default:
		*b = append((*b)[:index], (*b)[index+1:]...)
	}
}

// Clear removes all data sets from the batch.
func (b *Batch) Clear() {
	*b = *NewDataBatch()
}

// Clone returns a copy of the batch.
func (b *Batch) Clone() *Batch {
	clone := NewDataBatch()
	for _, set := range *b {
		clone.Add(set.Clone())
	}
	return clone
}

// ToCSV converts the data batch to a RFC4180-compliant CSV string.
func (b *Batch) ToCSV() []byte {
	if b == nil {
		return nil
	}

	var buf bytes.Buffer
	for i, set := range *b {
		if i == 0 {
			for j, data := range *set {
				if j == 0 {
					buf.WriteString(data.GetKey())
				} else {
					buf.WriteString("," + data.GetKey())
				}
			}
			buf.WriteString("\n")
		}
		for j, data := range *set {
			valueKind := data.GetValueType().GetTypeKind()
			var value string
			switch valueKind {
			case KindBool:
				value = fmt.Sprintf("%t", data.GetValueType().GetValue())
			case KindInt:
				value = fmt.Sprintf("%d", data.GetValueType().GetValue())
			case KindFloat:
				value = fmt.Sprintf("%f", data.GetValueType().GetValue())
			case KindBytes:
				value = fmt.Sprintf("%s", string(data.GetValueType().GetValue().([]byte)))
			default:
				value = fmt.Sprintf("%v", data.GetValueType().GetValue())
			}
			if j == 0 {
				buf.WriteString(value)
			} else {
				buf.WriteString("," + value)
			}
		}
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

// NewDataBatch returns a new data batch.
func NewDataBatch() *Batch {
	return &Batch{}
}
