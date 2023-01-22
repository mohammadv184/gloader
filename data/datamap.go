package data

type Map map[string]Type

func (m Map) Get(key string) Type {
	return m[key]
}
func (m Map) Set(key string, value Type) {
	m[key] = value
}
func (m Map) Delete(key string) {
	delete(m, key)
}
func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}
func (m Map) Len() int {
	return len(m)
}
func (m Map) Clone() Map {
	clone := Map{}
	for key, value := range m {
		clone[key] = value
	}
	return clone
}

func (m Map) Keys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
func (m Map) Types() []Type {
	types := make([]Type, 0, len(m))
	for _, value := range m {
		types = append(types, value)
	}
	return types
}
