package data

import (
	"bytes"
	"fmt"
)

// Batch is a collection of data sets.
type Batch []*Set

func (b *Batch) Add(set ...*Set) {
	*b = append(*b, set...)
}
func (b *Batch) Get(index int) *Set {
	if index >= b.GetLength() {
		return nil
	}
	return (*b)[index]
}

func (b *Batch) Pop() *Set {
	if b.GetLength() == 0 {
		return nil
	}
	set := b.Get(0)
	*b = (*b)[1:]
	return set
}

// GetSize returns the size of the batch in bytes.
func (b *Batch) GetSize() uint64 {
	var size uint64
	for _, set := range *b {
		size += set.GetSize()
	}
	return size
}

func (b *Batch) GetLength() int {
	return len(*b)
}

func (b *Batch) Remove(index int) {
	*b = append((*b)[:index], (*b)[index+1:]...)
}
func (b *Batch) Clear() {
	*b = nil
}

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
			valueKind := data.GetValue().GetTypeKind()
			var value string
			switch valueKind {
			case KindBool:
				value = fmt.Sprintf("%t", data.GetValue().GetValue())
			case KindInt:
				value = fmt.Sprintf("%d", data.GetValue().GetValue())
			case KindFloat:
				value = fmt.Sprintf("%f", data.GetValue().GetValue())
			case KindBytes:
				value = fmt.Sprintf("%s", string(data.GetValue().GetValue().([]byte)))
			default:
				value = fmt.Sprintf("%v", data.GetValue().GetValue())
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
func NewDataBatch() *Batch {
	return &Batch{}
}
