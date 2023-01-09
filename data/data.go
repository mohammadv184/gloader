package data

import "reflect"

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

func (d *Data) Clone() *Data {
	return NewData(d.Key, reflect.ValueOf(d.Value).Elem().Interface().(ValueType))
}

// NewData creates a new data with the given key and value.
func NewData(key string, value ValueType) *Data {
	return &Data{
		Key:   key,
		Value: value,
	}
}
