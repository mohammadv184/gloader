// Package types is a package that contains some generic data types.
package types

import (
	"fmt"
	"unsafe"

	"github.com/mohammadv184/gloader/data"
)

type StringType struct {
	data.BaseValueType
	value string
}

var _ data.Type = &StringType{}

func (t *StringType) Parse(p any) error {
	switch p.(type) {
	case string:
		t.value = p.(string)
		return nil
	default:
		return fmt.Errorf("%v: expected string, got %T", data.ErrInvalidValue, p)
	}
}

func (t *StringType) GetTypeKind() data.Kind {
	return data.KindString
}

func (t *StringType) GetTypeName() string {
	return "string"
}

func (t *StringType) GetValueSize() uint64 {
	return uint64(unsafe.Sizeof(t.value))
}

func (t *StringType) GetValue() any {
	return t.value
}

func NewStringType(value any) (data.Type, error) {
	stringType := &StringType{}
	return stringType, stringType.Parse(value)
}
