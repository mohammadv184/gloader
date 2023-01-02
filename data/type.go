package data

import (
	"errors"
	"reflect"
)

type Type interface {
	Parse(p []byte) error
	parseType(p []byte) error
	GetTypeName() string
	GetTypeKind() string
	GetTypeSize() int
	GetValue() []byte
	To(t Type) Type
}

type BaseType struct {
	size int
	kind reflect.Kind
}

func (t *BaseType) GetTypeName() string {
	return reflect.TypeOf(t).Name()
}
func (t *BaseType) GetTypeKind() string {
	return t.kind.String()
}
func (t *BaseType) GetTypeSize() int {
	return t.size
}
func (t *BaseType) Parse(p []byte) error {
	if reflect.TypeOf(t).Name() != "BaseType" {
		return errors.New("BaseType.Parse() can only be called on BaseType")
	}

	return reflect.
		ValueOf(t).
		MethodByName("parseType").
		Call([]reflect.Value{reflect.ValueOf(p)})[0].Interface().(error)

}
