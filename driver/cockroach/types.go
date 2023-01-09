package cockroach

import (
	"errors"
	"fmt"
	"gloader/data"
)

type ArrayType struct {
	data.BaseValueType
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
	return data.KindArray
}
func (t *ArrayType) GetTypeName() string {
	return "ARRAY"
}
func (t *ArrayType) GetTypeSize() int {
	return len(t.value)
}
func (t *ArrayType) GetValue() any {
	return t.value
}

type JsonBType struct {
	data.BaseValueType
	value []byte
}

func (t *JsonBType) Parse(v any) error {
	t.value = []byte(fmt.Sprintf("%v", v))
	return nil
}
func (t *JsonBType) GetTypeKind() data.Kind {
	return data.KindBytes
}
func (t *JsonBType) GetTypeName() string {
	return "JSONB"
}
func (t *JsonBType) GetTypeSize() int {
	return len(t.value)
}
func (t *JsonBType) GetValue() any {
	return t.value
}

var ErrTypeNotFound = errors.New("type not found")

var typeNamesMap = map[string]data.Type{
	"ARRAY": &ArrayType{},
	"JSONB": &JsonBType{},
	// TODO: add more types
}

func GetTypeFromName(name string) (data.Type, error) {
	t, ok := typeNamesMap[name]
	if !ok {
		return nil, ErrTypeNotFound
	}
	return t, nil
}
