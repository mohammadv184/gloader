package data

import (
	"fmt"
	"strings"
	"time"
)

// Set is a collection of related data. It`s used to collect and organize data.
// also in relational database, it`s called a row.
type Set []*Data

func (d *Set) Add(data *Data) {
	*d = append(*d, data)
}
func (d *Set) Get(key string) *Data {
	for _, data := range *d {
		if data.GetKey() == key {
			return data
		}
	}
	return nil
}
func (d *Set) GetByIndex(index int) *Data {
	return (*d)[index]
}

// GetSize returns the size of the data set in bytes.
func (d *Set) GetSize() uint64 {
	var size uint64
	for _, data := range *d {
		size += data.GetSize()
	}
	return size
}

func (d *Set) GetLength() int {
	return len(*d)
}

func (d *Set) Remove(key string) {
	for i, data := range *d {
		if data.GetKey() == key {
			*d = append((*d)[:i], (*d)[i+1:]...)
			break
		}
	}
}
func (d *Set) RemoveByIndex(index int) {
	*d = append((*d)[:index], (*d)[index+1:]...)
}
func (d *Set) Set(key string, value ValueType) {
	for _, data := range *d {
		if data.GetKey() == key {
			data.SetValue(value)
			return
		}
	}
	*d = append(*d, NewData(key, value))
}
func (d *Set) SetByIndex(index int, value ValueType) {
	(*d)[index].SetValue(value)
}
func (d *Set) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func (d *Set) String(delimiter string) string {
	var s strings.Builder
	for i, data := range *d {
		if i > 0 {
			s.WriteString(delimiter)
		}
		switch data.GetValue().GetTypeKind() {
		case KindInt, KindInt8, KindInt16, KindInt32, KindInt64,
			KindUint, KindUint8, KindUint16, KindUint32, KindUint64,
			KindFloat32, KindFloat64:
			s.WriteString(fmt.Sprintf("%v", data.GetValue().GetValue()))
		case KindBool:
			s.WriteString(fmt.Sprintf("%t", data.GetValue().GetValue()))
		case KindDateTime:
			s.WriteString(fmt.Sprintf("%s", data.GetValue().GetValue().(time.Time).Format("2006-01-02 15:04:05")))
		default:
			s.WriteString(fmt.Sprintf("%s", data.GetValue().GetValue()))
		}

	}
	return s.String()
}

func (d *Set) GetKeys() []string {
	keys := make([]string, len(*d))
	for i, data := range *d {
		keys[i] = data.GetKey()
	}
	return keys
}
func (d *Set) Clone() *Set {
	clone := NewDataSet()
	for _, data := range *d {
		clone.Add(data.Clone())
	}
	return clone
}

func NewDataSet() *Set {
	return &Set{}
}
