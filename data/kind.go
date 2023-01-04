package data

type Kind uint8

const (
	KindUnknown Kind = iota
	KindString
	KindInt
	KindFloat
	KindBool
	KindTime
	KindDate
	KindDateTime
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
	KindDate:      "date",
	KindDateTime:  "datetime",
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

func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "unknown"
}
func GetKindFromName(name string) Kind {
	for i, v := range kindNames {
		if v == name {
			return Kind(i)
		}
	}
	return KindUnknown
}
