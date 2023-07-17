package data

// Map is a map of data types. It is used to determine schema of dataCollection.
// The key is the name of the data type and the value is the data type.
type Map struct {
	typesMap        map[string]Type
	isNullableMap   map[string]bool
	hasDefaultValue map[string]bool
	typesIndex      []string
}

// Get returns the data type with the given name.
// It returns nil if the data type does not exist.
func (m *Map) Get(key string) Type {
	m.init()
	return m.typesMap[key]
}

// GetIndex returns the data type with the given index.
func (m *Map) GetIndex(index int) Type {
	m.init()
	return m.typesMap[m.typesIndex[index]]
}

func (m *Map) GetTypeMap() map[string]Type {
	m.init()
	return m.typesMap
}

func (m *Map) GetNullableMap() map[string]bool {
	m.init()
	return m.isNullableMap
}

func (m *Map) IsNullable(key string) bool {
	m.init()
	if _, ok := m.isNullableMap[key]; !ok {
		return false
	}

	return m.isNullableMap[key]
}

func (m *Map) HasDefaultValue(key string) bool {
	m.init()
	if _, ok := m.hasDefaultValue[key]; !ok {
		return false
	}

	return m.hasDefaultValue[key]
}

// Set sets the data type with the given name.
// Options[0] isNullable
// Options[1] hasDefaultValue.
func (m *Map) Set(key string, value Type, options ...bool) {
	m.init()
	for i, v := range options {
		switch i {
		case 0:
			m.isNullableMap[key] = v
		case 1:
			m.hasDefaultValue[key] = v
		}
	}

	m.typesMap[key] = value
	m.typesIndex = append(m.typesIndex, key)
}

// Delete deletes the data type with the given name.
func (m *Map) Delete(key string) {
	m.init()
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
	m.init()
	_, ok := m.typesMap[key]
	return ok
}

func (m *Map) HasIndex(index int) bool {
	m.init()
	return index >= 0 && index < len(m.typesIndex)
}

// Len returns the number of data types in the map.
func (m *Map) Len() int {
	m.init()
	return len(m.typesMap)
}

// Keys returns the keys of the map.
func (m *Map) Keys() []string {
	m.init()
	return m.typesIndex
}

func (m *Map) KeysExcept(keys []string) []string {
	m.init()
	result := make([]string, 0, len(m.typesIndex))
	for _, k := range m.typesIndex {
		for _, key := range keys {
			if k == key {
				continue
			}
			result = append(result, k)
		}
	}
	return result
}

func (m *Map) NotNullableKeys() []string {
	m.init()
	keys := make([]string, 0, len(m.typesMap))
	for i, k := range m.typesIndex {
		if !m.isNullableMap[k] {
			keys = append(keys, m.typesIndex[i])
		}
	}
	return keys
}

// Types returns the values of the map.
func (m *Map) Types() []Type {
	m.init()
	types := make([]Type, 0, len(m.typesMap))
	for _, value := range m.typesMap {
		types = append(types, value)
	}
	return types
}

func (m *Map) init() {
	if m.typesMap == nil {
		m.typesMap = make(map[string]Type)
	}
	if m.isNullableMap == nil {
		m.isNullableMap = make(map[string]bool)
	}
	if m.hasDefaultValue == nil {
		m.hasDefaultValue = make(map[string]bool)
	}
}
