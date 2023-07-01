package cockroach

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/gloader/data"
)

// ArrayType is a type for array.
type ArrayType struct {
	data.BaseValueType
	value []any
}

// Parse parses the value and stores it in the receiver.
func (t *ArrayType) Parse(v any) error {
	switch v.(type) {
	case []any:
		t.value = v.([]any)
		return nil
	default:
		return fmt.Errorf("%v: expected []any, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *ArrayType) GetTypeKind() data.Kind {
	return data.KindArray
}

// GetTypeName returns the name of the type.
func (t *ArrayType) GetTypeName() string {
	return "ARRAY"
}

// GetValueSize returns the size of the value in bytes.
func (t *ArrayType) GetValueSize() int {
	return len(t.value)
}

// GetValue returns the value stored in the receiver.
func (t *ArrayType) GetValue() any {
	return t.value
}

// JSONBType is a type for jsonb.
type JSONBType struct {
	data.BaseValueType
	value []byte
}

// Parse parses the value and stores it in the receiver.
func (t *JSONBType) Parse(v any) error {
	t.value = []byte(fmt.Sprintf("%v", v))
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *JSONBType) GetTypeKind() data.Kind {
	return data.KindBytes
}

// GetTypeName returns the name of the type.
func (t *JSONBType) GetTypeName() string {
	return "JSONB"
}

// GetValueSize returns the size of the value in bytes.
func (t *JSONBType) GetValueSize() int {
	return len(t.value)
}

// GetValue returns the value stored in the receiver.
func (t *JSONBType) GetValue() any {
	return t.value
}

var ErrTypeNotFound = errors.New("type not found") // ErrTypeNotFound is returned when the type is not found.
// GetTypeFromName returns the type from the given name.
func GetTypeFromName(name string) (data.Type, error) {
	switch name {
	case "ARRAY":
		return &ArrayType{}, nil
	case "JSONB":
		return &JSONBType{}, nil
	default:
		return nil, ErrTypeNotFound
	}
}
