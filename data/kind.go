package data

// Kind represents the kind of value.
// It is a subset of the reflect.Kind.
// The zero Kind is KindUnknown.
type Kind uint8

const (
	KindUnknown Kind = iota
	KindString
	KindInt
	KindFloat
	KindBool
	KindTime
	KindTimestamp
	KindDuration
	KindBytes
	KindArray
	KindMap
	KindStruct
	KindInterface
	KindPointer
	KindFunc
	KindChan
	KindSlice
	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64
	KindInt8
	KindInt16
	KindInt32
	KindInt64
	KindFloat32
	KindFloat64
)

var kindNames = [...]string{
	KindUnknown:   "unknown",
	KindString:    "string",
	KindInt:       "int",
	KindFloat:     "float",
	KindBool:      "bool",
	KindTime:      "time",
	KindTimestamp: "timestamp",
	KindDuration:  "duration",
	KindBytes:     "bytes",
	KindArray:     "array",
	KindMap:       "map",
	KindStruct:    "struct",
	KindInterface: "interface",
	KindPointer:   "pointer",
	KindFunc:      "func",
	KindChan:      "chan",
	KindSlice:     "slice",
	KindUint:      "uint",
	KindUint8:     "uint8",
	KindUint16:    "uint16",
	KindUint32:    "uint32",
	KindUint64:    "uint64",
	KindInt8:      "int8",
	KindInt16:     "int16",
	KindInt32:     "int32",
	KindInt64:     "int64",
	KindFloat32:   "float32",
	KindFloat64:   "float64",
}

var baseKindSizes = [...]int{
	KindUnknown:   0,
	KindString:    4, // considering max single utf8 char
	KindInt:       8,
	KindFloat:     8,
	KindBool:      1,
	KindTime:      8,
	KindTimestamp: 8,
	KindDuration:  8,
	KindBytes:     1,
	KindUint:      8,
	KindUint8:     1,
	KindUint16:    2,
	KindUint32:    4,
	KindUint64:    8,
	KindInt8:      1,
	KindInt16:     2,
	KindInt32:     4,
	KindInt64:     8,
	KindFloat32:   4,
	KindFloat64:   8,
}

// GetBaseKindSize returns the size of the base kinds such as int, float, bool, etc. in bytes.
// It returns 0 if the kind is not a base kind.
func GetBaseKindSize(k Kind) int {
	if int(k) < len(baseKindSizes) {
		return baseKindSizes[k]
	}
	return 0
}

// String returns the name of the kind.
func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "unknown"
}

// GetKindFromName returns the kind from the given name.
// It returns KindUnknown if the name is not found.
func GetKindFromName(name string) Kind {
	for i, v := range kindNames {
		if v == name {
			return Kind(i)
		}
	}
	return KindUnknown
}
