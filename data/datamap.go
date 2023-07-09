package data

// Map is a map of data types. It is used to determine schema of dataCollection.
// The key is the name of the data type and the value is the data type.
type Map struct {
	typesMap      map[string]Type
	isNullableMap map[string]bool
	typesIndex    []string
}

// Get returns the data type with the given name.
// It returns nil if the data type does not exist.
func (m *Map) Get(key string) Type {
	return m.typesMap[key]
}

// GetIndex returns the data type with the given index.
func (m *Map) GetIndex(index int) Type {
	return m.typesMap[m.typesIndex[index]]
}

func (m *Map) IsNullable(key string) bool {
	if _, ok := m.isNullableMap[key]; !ok {
		return false
	}

	return m.isNullableMap[key]
}

// Set sets the data type with the given name.
func (m *Map) Set(key string, value Type, isNullable ...bool) {
	if len(isNullable) > 0 {
		m.isNullableMap[key] = isNullable[0]
	}

	m.typesMap[key] = value
	m.typesIndex = append(m.typesIndex, key)
}

// Delete deletes the data type with the given name.
func (m *Map) Delete(key string) {
	for i, k := range m.typesIndex {
		if k == key {
			if i == len(m.typesIndex)-1 {
				m.typesIndex = m.typesIndex[:i]
			} else {
				m.typesIndex = append(m.typesIndex[:i], m.typesIndex[i+1:]...)
			}
			break
		}
	}
	delete(m.typesMap, key)
	delete(m.isNullableMap, key)
}

// Has returns true if the data type with the given name exists.
func (m *Map) Has(key string) bool {
	_, ok := m.typesMap[key]
	return ok
}

func (m *Map) HasIndex(index int) bool {
	return index >= 0 && index < len(m.typesIndex)
}

// Len returns the number of data types in the map.
func (m *Map) Len() int {
	return len(m.typesMap)
}

// Keys returns the keys of the map.
func (m *Map) Keys() []string {
	keys := make([]string, 0, len(m.typesMap))
	for key := range m.typesMap {
		keys = append(keys, key)
	}
	return keys
}

// Types returns the values of the map.
func (m *Map) Types() []Type {
	types := make([]Type, 0, len(m.typesMap))
	for _, value := range m.typesMap {
		types = append(types, value)
	}
	return types
}
