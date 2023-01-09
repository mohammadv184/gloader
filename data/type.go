package data

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type Type interface {
	GetTypeName() string
	GetTypeKind() Kind
	GetTypeSize() int
}

type ValueType interface {
	Type
	Parse(v any) error
	GetValue() any
	To(t ValueType) (ValueType, error)
	AssignTo(t any) error
}

type BaseValueType struct{}

var _ Type = &BaseValueType{}

func (_ *BaseValueType) Parse(_ any) error {
	return ErrParseFuncNotImplemented
}

func (_ *BaseValueType) GetTypeKind() Kind {
	return KindUnknown
}
func (b *BaseValueType) GetTypeName() string {
	return reflect.TypeOf(b).String()
}
func (b *BaseValueType) GetTypeSize() int {
	return int(unsafe.Sizeof(b))
}
func (b *BaseValueType) GetValue() any {
	return nil
}
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
	//TODO: add more types
	return errors.New("not implemented")
}
