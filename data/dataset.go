package data

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
func (d *Set) GetSize() int {
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
func (d *Set) Set(key string, value Type) {
	for _, data := range *d {
		if data.GetKey() == key {
			data.SetValue(value)
			return
		}
	}
	*d = append(*d, NewData(key, value))
}
func (d *Set) SetByIndex(index int, value Type) {
	(*d)[index].SetValue(value)
}
func NewDataSet() *Set {
	return &Set{}
}
