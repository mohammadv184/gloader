package data

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
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
}

// ValueType is Type that holds a value.
// It is provided Parse, GetTypeSize, GetValue, To, AssignTo, Clone methods. for handling values.
type ValueType interface {
	Type // Type interface
	// Parse parses the value and stores it in the receiver.
	Parse(v any) error
	// GetTypeSize returns the size of the value in bytes.
	GetTypeSize() uint64
	// GetValue returns the value stored in the receiver.
	GetValue() any
	// To converts the value to the given type.
	To(t ValueType) (ValueType, error)
	// AssignTo assigns the value to the given destination.
	AssignTo(t any) error
	// Clone returns a copy of the receiver.
	Clone() ValueType
}

// BaseValueType implements ValueType interface.
// It can be embedded in other types to implement ValueType interface quickly.
type BaseValueType struct{}

var _ Type = &BaseValueType{} // BaseValueType implements Type interface.
// Parse parses the value and stores it in the receiver.
// It should be implemented by the parent type. otherwise it returns ErrParseFuncNotImplemented.
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

// GetTypeSize returns the size of the value in bytes.
func (b *BaseValueType) GetTypeSize() uint64 {
	return uint64(unsafe.Sizeof(b))
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

// To converts the value to the given type.
func (b *BaseValueType) To(t ValueType) (ValueType, error) {
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

// AssignTo assigns the value to the given destination.
// It returns ErrDestMustBePointer if the destination is not a pointer.
func (b *BaseValueType) AssignTo(dest any) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return ErrDestMustBePointer
	}

	switch b.GetTypeKind() {
	case KindBool:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = b.GetValue().(bool)
			return nil
		case *int:
			*dest.(*int) = 0
			if b.GetValue().(bool) {
				*dest.(*int) = 1
			}
			return nil
		case *int8:
			*dest.(*int8) = 0
			if b.GetValue().(bool) {
				*dest.(*int8) = 1
			}
			return nil
		case *int16:
			*dest.(*int16) = 0
			if b.GetValue().(bool) {
				*dest.(*int16) = 1
			}
			return nil
		case *int32:
			*dest.(*int32) = 0
			if b.GetValue().(bool) {
				*dest.(*int32) = 1
			}
			return nil
		case *int64:
			*dest.(*int64) = 0
			if b.GetValue().(bool) {
				*dest.(*int64) = 1
			}
			return nil
		case *uint:
			*dest.(*uint) = 0
			if b.GetValue().(bool) {
				*dest.(*uint) = 1
			}
			return nil
		case *uint8:
			*dest.(*uint8) = 0
			if b.GetValue().(bool) {
				*dest.(*uint8) = 1
			}
			return nil
		case *uint16:
			*dest.(*uint16) = 0
			if b.GetValue().(bool) {
				*dest.(*uint16) = 1
			}
			return nil
		case *uint32:
			*dest.(*uint32) = 0
			if b.GetValue().(bool) {
				*dest.(*uint32) = 1
			}
			return nil
		case *uint64:
			*dest.(*uint64) = 0
			if b.GetValue().(bool) {
				*dest.(*uint64) = 1
			}
			return nil
		case *float32:
			*dest.(*float32) = 0
			if b.GetValue().(bool) {
				*dest.(*float32) = 1
			}
			return nil
		case *float64:
			*dest.(*float64) = 0
			if b.GetValue().(bool) {
				*dest.(*float64) = 1
			}
			return nil
		case *string:
			*dest.(*string) = fmt.Sprintf("%v", b.GetValue())
			return nil
		default:
			return fmt.Errorf("from bool to %T: %w", dest, ErrDestNotAssignable)
		}
	case KindInt:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = b.GetValue().(int)
			return nil
		case *int8:
			*dest.(*int8) = int8(b.GetValue().(int))
			return nil
		case *int16:
			*dest.(*int16) = int16(b.GetValue().(int))
			return nil
		case *int32:
			*dest.(*int32) = int32(b.GetValue().(int))
			return nil
		case *int64:
			*dest.(*int64) = int64(b.GetValue().(int))
			return nil
		default:
			return fmt.Errorf("from int to %T: %w", dest, ErrDestNotAssignable)
		}

	case KindInt8:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int8) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = int(b.GetValue().(int8))
			return nil
		case *int8:
			*dest.(*int8) = b.GetValue().(int8)
			return nil
		case *int16:
			*dest.(*int16) = int16(b.GetValue().(int8))
			return nil
		case *int32:
			*dest.(*int32) = int32(b.GetValue().(int8))
			return nil
		case *int64:
			*dest.(*int64) = int64(b.GetValue().(int8))
			return nil
		default:
			return fmt.Errorf("from int8 to %T: %w", dest, ErrDestNotAssignable)
		}

	case KindInt16:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int16) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = int(b.GetValue().(int16))
			return nil
		case *int8:
			*dest.(*int8) = int8(b.GetValue().(int16))
			return nil
		case *int16:
			*dest.(*int16) = b.GetValue().(int16)
			return nil
		case *int32:
			*dest.(*int32) = int32(b.GetValue().(int16))
			return nil
		case *int64:
			*dest.(*int64) = int64(b.GetValue().(int16))
			return nil
		default:
			return fmt.Errorf("from int16 to %T: %w", dest, ErrDestNotAssignable)
		}

	case KindInt32:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int32) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = int(b.GetValue().(int32))
			return nil
		case *int8:
			*dest.(*int8) = int8(b.GetValue().(int32))
			return nil
		case *int16:
			*dest.(*int16) = int16(b.GetValue().(int32))
			return nil
		case *int32:
			*dest.(*int32) = b.GetValue().(int32)
			return nil
		case *int64:
			*dest.(*int64) = int64(b.GetValue().(int32))
			return nil
		default:
			return fmt.Errorf("from int32 to %T: %w", dest, ErrDestNotAssignable)
		}

	case KindInt64:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int64) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = int(b.GetValue().(int64))
			return nil
		case *int8:
			*dest.(*int8) = int8(b.GetValue().(int64))
			return nil
		case *int16:
			*dest.(*int16) = int16(b.GetValue().(int64))
			return nil
		case *int32:
			*dest.(*int32) = int32(b.GetValue().(int64))
			return nil
		case *int64:
			*dest.(*int64) = b.GetValue().(int64)
			return nil
		default:
			return fmt.Errorf("from int64 to %T: %w", dest, ErrDestNotAssignable)
		}
	case KindUint:
		switch dest.(type) {
		case *bool:
			*dest.(*bool) = false
			if b.GetValue().(int) != 0 {
				*dest.(*bool) = true
			}
			return nil
		case *int:
			*dest.(*int) = b.GetValue().(int)
			return nil
		case *int8:
			*dest.(*int8) = int8(b.GetValue().(int))
			return nil
		case *int16:
			*dest.(*int16) = int16(b.GetValue().(int))
			return nil
		case *int32:
			*dest.(*int32) = int32(b.GetValue().(int))
			return nil
		case *int64:
			*dest.(*int64) = int64(b.GetValue().(int))
			return nil
		}
		return fmt.Errorf("from int to %T: %w", dest, ErrDestNotAssignable)
	}
	// TODO: add more types
	return errors.New("not implemented")
}
