package data

// Data is a key-value pair. The key is a string and the value is a Type.
// The value can be any type that implements the Type interface.
// Data is the smallest unit of data in gloader.
type Data struct {
	Key   string
	Value Type
}

// GetKey returns the key of the data.
func (d *Data) GetKey() string {
	return d.Key
}

// GetValue returns the value of the data.
func (d *Data) GetValue() Type {
	return d.Value
}

// SetKey sets the key of the data.
func (d *Data) SetKey(key string) {
	d.Key = key
}

// SetValue sets the value of the data.
func (d *Data) SetValue(value Type) {
	d.Value = value
}

// NewData creates a new data with the given key and value.
func NewData(key string, value Type) *Data {
	return &Data{
		Key:   key,
		Value: value,
	}
}
