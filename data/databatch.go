package data

import (
	"bytes"
	"fmt"
)

// Batch is a collection of data sets.
type Batch []*Set

func (b *Batch) Add(set *Set) {
	*b = append(*b, set)
}
func (b *Batch) Get(index int) *Set {
	return (*b)[index]
}
func (b *Batch) GetSize() int {
	return len(*b)
}
func (b *Batch) Remove(index int) {
	*b = append((*b)[:index], (*b)[index+1:]...)
}

// ToCSV converts the data batch to a RFC4180-compliant CSV string.
func (b *Batch) ToCSV() []byte {
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
			case KindString:
				value = fmt.Sprintf("%s", data.GetValue().GetValue())
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
