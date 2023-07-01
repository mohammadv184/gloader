package data

// Map is a map of data types. It is used to determine schema of dataCollection.
// The key is the name of the data type and the value is the data type.
type Map map[string]Type

// Get returns the data type with the given name.
// It returns nil if the data type does not exist.
func (m Map) Get(key string) Type {
	return m[key]
}

// Set sets the data type with the given name.
func (m Map) Set(key string, value Type) {
	m[key] = value
}

// Delete deletes the data type with the given name.
func (m Map) Delete(key string) {
	delete(m, key)
}

// Has returns true if the data type with the given name exists.
func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Len returns the number of data types in the map.
func (m Map) Len() int {
	return len(m)
}

// Keys returns the keys of the map.
func (m Map) Keys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Types returns the values of the map.
func (m Map) Types() []Type {
	types := make([]Type, 0, len(m))
	for _, value := range m {
		types = append(types, value)
	}
	return types
}
