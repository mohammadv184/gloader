package types

import (
	"fmt"
	"gloader/data"
)

type IntegerType struct {
	data.BaseValueType
	value int64
}

var _ data.Type = &IntegerType{}

func (t *IntegerType) Parse(v any) error {
	switch v.(type) {
	case int:
		t.value = int64(v.(int))
		return nil
	case int8:
		t.value = int64(v.(int8))
		return nil
	case int16:
		t.value = int64(v.(int16))
		return nil
	case int32:
		t.value = int64(v.(int32))
		return nil
	case int64:
		t.value = v.(int64)
		return nil
	default:
		return fmt.Errorf("%v: expected int, got %T", data.ErrInvalidValue, v)
	}
}

func (t *IntegerType) GetTypeKind() data.Kind {
	return data.KindInt
}

func (t *IntegerType) GetTypeName() string {
	return "int"
}

func (t *IntegerType) GetTypeSize() uint64 {
	return 8
}

func (t *IntegerType) GetValue() any {
	return t.value
}

func NewIntegerType(value any) (data.Type, error) {
	integerType := &IntegerType{}
	return integerType, integerType.Parse(value)
}
