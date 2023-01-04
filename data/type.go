package data

import (
	"reflect"
	"unsafe"
)

type Type interface {
	Parse(v any) error
	GetTypeName() string
	GetTypeKind() Kind
	GetTypeSize() int
	GetValue() any
	To(t Type) (Type, error)
}

type BaseType struct{}

var _ Type = &BaseType{}

func (_ *BaseType) Parse(_ any) error {
	return ErrParseFuncNotImplemented
}

func (_ *BaseType) GetTypeKind() Kind {
	return GetKindFromName(KindUnknown.String())
}
func (b *BaseType) GetTypeName() string {
	return reflect.TypeOf(b).String()
}
func (b *BaseType) GetTypeSize() int {
	return int(unsafe.Sizeof(b))
}
func (b *BaseType) GetValue() any {
	return nil
}
func (b *BaseType) To(t Type) (Type, error) {
	if b.GetTypeKind() != t.GetTypeKind() {
		return nil, ErrDataTypeKindNotMatch
	}
	if b.GetTypeName() == t.GetTypeName() {
		return t, nil
	}
	err := t.Parse(b.GetValue())
	if err != nil {
		return nil, err
	}
	return t, nil
}
