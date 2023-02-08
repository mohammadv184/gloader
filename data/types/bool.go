package types

import (
	"fmt"
	"gloader/data"
)

type BoolType struct {
	data.BaseValueType
	value bool
}

func (t *BoolType) Parse(p any) error {
	switch p.(type) {
	case bool:
		t.value = p.(bool)
		return nil
	default:
		return fmt.Errorf("%v: expected bool, got %T", data.ErrInvalidValue, p)
	}
}

func (t *BoolType) GetTypeKind() data.Kind {
	return data.KindBool
}

func (t *BoolType) GetTypeName() string {
	return "bool"
}

func (t *BoolType) GetTypeSize() uint64 {
	return 1
}

func (t *BoolType) GetValue() any {
	if t.value {
		return []byte{1}
	}
	return []byte{0}
}

func NewBoolType(value any) (data.Type, error) {
	boolType := &BoolType{}
	return boolType, boolType.Parse(value)
}
