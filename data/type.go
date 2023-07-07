package data

import (
	"reflect"
)

// Type is the interface that has the basic data type methods.
// It is provided GetTypeName, GetTypeKind methods.
// All driver types must implement this interface.
// A Type represents a type in the database.
type Type interface {
	// GetTypeName returns the name of the type.
	GetTypeName() string
	// GetTypeKind returns the kind of the type.
	GetTypeKind() Kind
	// GetTypeSize returns the size of 1 length value of the type in bytes.
	GetTypeSize() uint64
}

// ValueType is a Type that holds a value with same kind.
// It is provided Parse, GetValueSize, GetValue, To, AssignTo, Clone methods. For handling values.
type ValueType interface {
	Type // Type interface
	// Parse parses the value and stores it in the receiver.
	// It can parse any type that is compatible with the type kind and nil.
	Parse(v any) error
	// GetValueSize returns the size of the value in bytes.
	GetValueSize() uint64
	// GetValue returns the value stored in the receiver.
	// returned value kind is the same as the type kind which is accessible by GetTypeKind method.
	// note: the returned value could be nil.
	GetValue() any
	// To convert the value to the given type.
	To(t Type) (ValueType, error)
	// AssignTo assigns the value to the given pointer destination.
	// The t should be a pointer to a type that is same with the type kind.
	AssignTo(t any) error
	// Clone returns a copy of the receiver.
	Clone() ValueType
}

// BaseValueType implements ValueType interface.
// It can be embedded in other types to implement ValueType interface quickly.
type BaseValueType struct{}

var _ ValueType = &BaseValueType{} // BaseValueType implements Type interface.
// Parse parses the value and stores it in the receiver.
// It should be implemented by the parent type. Otherwise, it returns ErrParseFuncNotImplemented.
func (*BaseValueType) Parse(_ any) error {
	return ErrParseFuncNotImplemented
}

// GetTypeKind returns the kind of the type.
// if not implemented by the parent type, it returns KindUnknown.
func (*BaseValueType) GetTypeKind() Kind {
	return KindUnknown
}

// GetTypeName returns the name of the type.
func (b *BaseValueType) GetTypeName() string {
	return reflect.TypeOf(b).String()
}

// GetTypeSize returns the size of the type in bytes.
func (b *BaseValueType) GetTypeSize() uint64 {
	return uint64(GetBaseKindSize(b.GetTypeKind()))
}

// GetValueSize returns the size of the value in bytes.
func (b *BaseValueType) GetValueSize() uint64 {
	return b.GetTypeSize()
}

// GetValue returns the value stored in the receiver.
func (b *BaseValueType) GetValue() any {
	return nil
}

// Clone returns a copy of the receiver.
func (b *BaseValueType) Clone() ValueType {
	valueType := reflect.New(reflect.TypeOf(b).Elem()).Interface().(ValueType)
	_ = valueType.Parse(b.GetValue())
	return valueType
}

// To convert the value to the given type.
func (b *BaseValueType) To(t Type) (ValueType, error) {
	if b.GetTypeKind() != t.GetTypeKind() {
		return nil, ErrDataTypeKindNotMatch
	}

	vt := GetNewValueTypeAs(t)

	err := vt.Parse(b.GetValue())
	if err != nil {
		return nil, err
	}

	return vt, nil
}

// AssignTo assigns the value to the given destination.
// It returns ErrDestMustBePointer if the destination is not a pointer.
// It returns ErrDataTypeKindNotMatch if the destination type kind is not the same as the receiver type kind.
func (b *BaseValueType) AssignTo(dest any) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return ErrDestMustBePointer
	}

	if reflect.TypeOf(dest).Elem().Kind() != b.GetTypeKind().GetReflectKind() {
		return ErrDataTypeKindNotMatch
	}

	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(b.GetValue()))
	return nil
}
