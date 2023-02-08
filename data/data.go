// Package data contains the all data related types and functionalities.
// such as data.Data, data.ValueType, data.Kind, etc.
package data

import "unsafe"

// Data is a key-value pair. The key is a string and the value is a Type.
// The value can be any type that implements the Type interface.
// Data is the smallest unit of data in gloader.
type Data struct {
	Key   string
	Value ValueType
}

// GetKey returns the key of the data.
func (d *Data) GetKey() string {
	return d.Key
}

// GetValue returns the value of the data.
func (d *Data) GetValue() ValueType {
	return d.Value
}

// SetKey sets the key of the data.
func (d *Data) SetKey(key string) {
	d.Key = key
}

// SetValue sets the value of the data.
func (d *Data) SetValue(value ValueType) {
	d.Value = value
}

// Clone returns a copy of the data.
func (d *Data) Clone() *Data {
	return NewData(d.Key, d.Value.Clone())
}

// GetSize returns the size of the data in bytes.
func (d *Data) GetSize() uint64 {
	return d.Value.GetTypeSize() + uint64(unsafe.Sizeof(d.Key))
}

// NewData creates a new data with the given key and value.
// The value must implement the ValueType interface.
func NewData(key string, value ValueType) *Data {
	return &Data{
		Key:   key,
		Value: value,
	}
}
