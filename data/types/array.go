package types

import (
	"fmt"
	"gloader/data"
	"unsafe"
)

type ArrayType struct {
	data.BaseValueType
	value []any
}

func (t *ArrayType) Parse(v any) error {
	switch v.(type) {
	case []any:
		t.value = v.([]any)
		return nil
	default:
		return fmt.Errorf("%v: expected []any, got %T", data.ErrInvalidValue, v)
	}
}

func (t *ArrayType) GetTypeKind() data.Kind {
	return data.KindArray
}
func (t *ArrayType) GetTypeName() string {
	return "array"
}
func (t *ArrayType) GetTypeSize() uint64 {
	return uint64(unsafe.Sizeof(t.value))
}
func (t *ArrayType) GetValue() any {
	return t.value
}
