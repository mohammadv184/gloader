package types

import (
	"fmt"
	"gloader/data"
	"unsafe"
)

type ArrayType struct {
	data.BaseType
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
	return data.GetKindFromName(data.KindArray.String())
}
func (t *ArrayType) GetTypeName() string {
	return "array"
}
func (t *ArrayType) GetTypeSize() int {
	return int(unsafe.Sizeof(t.value))
}
func (t *ArrayType) GetValue() any {
	return t.value
}
