// Package types is a package that contains some generic data types.
package types

import (
	"fmt"
	"unsafe"

	"github.com/mohammadv184/gloader/data"
)

// ArrayType is a generic array type that can hold any array data type.
// It implements the data.ValueType interface.
type ArrayType struct {
	data.BaseValueType
	value []any
}

// Parse parses the given value and sets it to the ArrayType.
// It returns an error if the given value is not of type []any.
func (t *ArrayType) Parse(v any) error {
	switch v.(type) {
	case []any:
		t.value = v.([]any)
		return nil
	default:
		return fmt.Errorf("%v: expected []any, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the data.Kind of the ArrayType.
func (t *ArrayType) GetTypeKind() data.Kind {
	return data.KindArray
}

func (t *ArrayType) GetTypeName() string {
	return "array"
}

func (t *ArrayType) GetTypeSize() uint64 {
	return uint64(unsafe.Sizeof(t))
}

func (t *ArrayType) GetValueSize() uint64 {
	return uint64(unsafe.Sizeof(t.value))
}

func (t *ArrayType) GetValue() any {
	return t.value
}
