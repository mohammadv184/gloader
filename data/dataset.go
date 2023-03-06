package data

import (
	"fmt"
	"strings"
	"time"
)

// Set is a collection of related data. It`s used to collect and organize data.
// also in relational database, it`s called a row.
type Set []*Data

// Add adds the given data to the set.
// It does not check for duplicate data.
func (d *Set) Add(data *Data) {
	*d = append(*d, data)
}

// Get returns the data with the given key.
// It returns nil if the data does not exist.
func (d *Set) Get(key string) *Data {
	for _, data := range *d {
		if data.GetKey() == key {
			return data
		}
	}
	return nil
}

// GetByIndex returns the data at the given index.
// It returns nil if the index is out of range.
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

// GetLength returns the length of the data set.
func (d *Set) GetLength() int {
	return len(*d)
}

// Remove removes the data with the given key.
// if the key is not found, it does nothing.
func (d *Set) Remove(key string) {
	for i, data := range *d {
		if data.GetKey() == key {
			*d = append((*d)[:i], (*d)[i+1:]...)
			break
		}
	}
}

// RemoveByIndex removes the data at the given index.
// It does nothing if the index is out of range.
func (d *Set) RemoveByIndex(index int) {
	*d = append((*d)[:index], (*d)[index+1:]...)
}

// Set sets the value of the data with the given key.
// If the data does not exist, it creates a new data with the given key and value.
func (d *Set) Set(key string, value ValueType) {
	for _, data := range *d {
		if data.GetKey() == key {
			data.SetValue(value)
			return
		}
	}
	*d = append(*d, NewData(key, value))
}

// SetByIndex sets the value of the data at the given index.
// It does nothing if the index is out of range.
func (d *Set) SetByIndex(index int, value ValueType) {
	(*d)[index].SetValue(value)
}

// Swap swaps the data at the given indexes.
func (d *Set) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

// String returns the string representation of the data set.
// It uses the given delimiter to separate the data.
// e.g. delimiter = ", " => "value1, value2".
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

// GetKeys returns the keys of the data set.
func (d *Set) GetKeys() []string {
	keys := make([]string, len(*d))
	for i, data := range *d {
		keys[i] = data.GetKey()
	}
	return keys
}

// GetValues returns the values of the data set.
func (d *Set) GetValues() []ValueType {
	values := make([]ValueType, len(*d))
	for i, data := range *d {
		values[i] = data.GetValue()
	}
	return values
}

// GetStringValues returns the string values of the data set.
func (d *Set) GetStringValues() []string {
	values := make([]string, len(*d))
	for i, data := range *d {
		switch data.GetValue().GetTypeKind() {
		case KindInt, KindInt8, KindInt16, KindInt32, KindInt64,
			KindUint, KindUint8, KindUint16, KindUint32, KindUint64,
			KindFloat32, KindFloat64:
			values[i] = fmt.Sprintf("%v", data.GetValue().GetValue())
		case KindBool:
			values[i] = fmt.Sprintf("%t", data.GetValue().GetValue())
		case KindDateTime:
			values[i] = fmt.Sprintf("%s", data.GetValue().GetValue().(time.Time).Format("2006-01-02 15:04:05.999999999"))
		default:
			values[i] = fmt.Sprintf("%s", data.GetValue().GetValue())
		}
	}
	return values
}

// Clone returns a clone of the data set.
func (d *Set) Clone() *Set {
	clone := NewDataSet()
	for _, data := range *d {
		clone.Add(data.Clone())
	}
	return clone
}

// NewDataSet returns a new data set.
func NewDataSet() *Set {
	return &Set{}
}
