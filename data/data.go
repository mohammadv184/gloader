package data

type Data struct {
	Key   string
	Value Type
}

func (d *Data) GetKey() string {
	return d.Key
}
func (d *Data) GetValue() Type {
	return d.Value
}
func (d *Data) SetKey(key string) {
	d.Key = key
}
func (d *Data) SetValue(value Type) {
	d.Value = value
}
func NewData(key string, value Type) *Data {
	return &Data{
		Key:   key,
		Value: value,
	}
}
