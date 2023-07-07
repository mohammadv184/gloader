package cockroach

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unsafe"

	"github.com/mohammadv184/gloader/data"
)

// ref: https://www.cockroachlabs.com/docs/stable/data-types.html

// ArrayType is a type for array.
type ArrayType struct {
	data.BaseValueType
	value    []any
	hasValue bool
}

var _ data.ValueType = &ArrayType{}

// Parse parses the value and stores it in the receiver.
func (t *ArrayType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}
	switch v.(type) {
	case []any:
		t.value = v.([]any)
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected []any, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *ArrayType) GetTypeKind() data.Kind {
	return data.KindSlice
}

// GetTypeName returns the name of the type.
func (t *ArrayType) GetTypeName() string {
	return "ARRAY"
}

func (t *ArrayType) GetTypeSize() uint64 {
	return 0
}

// GetValueSize returns the size of the value in bytes.
func (t *ArrayType) GetValueSize() uint64 {
	// Calculating the size of reference types is difficult for now.
	// So, we return the size with unsafe package.
	// TODO: calculate the size of reference types.
	return uint64(unsafe.Sizeof(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *ArrayType) GetValue() any {
	if !t.hasValue {
		return nil
	}

	return t.value
}

type BitType struct {
	data.BaseValueType
	value    []bool
	hasValue bool
}

var _ data.ValueType = &BitType{}

func (t *BitType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch tv := v.(type) {
	case []bool:
		t.value = tv
		t.hasValue = true
		return nil
	case []int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64:
		t.value = make([]bool, 0)
		for _, i := range v.([]int) {
			t.value = append(t.value, i == 1)
		}
		t.hasValue = true
		return nil
	case []string:
		t.value = make([]bool, 0)
		for _, s := range tv {
			t.value = append(t.value, s == "1")
		}
		t.hasValue = true
		return nil
	case []any:
		t.value = make([]bool, 0)
		for _, i := range tv {
			switch i.(type) {
			case int:
				t.value = append(t.value, i.(int) == 1)
			case int8:
				t.value = append(t.value, i.(int8) == 1)
			case int16:
				t.value = append(t.value, i.(int16) == 1)
			case int32:
				t.value = append(t.value, i.(int32) == 1)
			case int64:
				t.value = append(t.value, i.(int64) == 1)
			case uint:
				t.value = append(t.value, i.(uint) == 1)
			case uint8:
				t.value = append(t.value, i.(uint8) == 1)
			case uint16:
				t.value = append(t.value, i.(uint16) == 1)
			case uint32:
				t.value = append(t.value, i.(uint32) == 1)
			case uint64:
				t.value = append(t.value, i.(uint64) == 1)
			case float32:
				t.value = append(t.value, i.(float32) == 1)
			case float64:
				t.value = append(t.value, i.(float64) == 1)
			case string:
				t.value = append(t.value, i.(string) == "1")
			default:

			}

		}
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected []any, got %T", data.ErrInvalidValue, v)
	}
}

func (t *BitType) GetTypeKind() data.Kind {
	return data.KindSlice
}

func (t *BitType) GetTypeName() string {
	return "BIT"
}

func (t *BitType) GetTypeSize() uint64 {
	return 1
}

func (t *BitType) GetValueSize() uint64 {
	return t.GetTypeSize() * uint64(len(t.value))
}

func (t *BitType) GetValue() any {
	if !t.hasValue {
		return nil
	}

	return t.value
}

type BoolType struct {
	data.BaseValueType
	value    bool
	hasValue bool
}

var _ data.ValueType = &BoolType{}

func (t *BoolType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case bool:
		t.value = v.(bool)
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected bool, got %T", data.ErrInvalidValue, v)
	}
}
func (t *BoolType) GetTypeKind() data.Kind {
	return data.KindBool
}

func (t *BoolType) GetTypeName() string {
	return "BOOL"
}
func (t *BoolType) GetTypeSize() uint64 {
	return 1
}
func (t *BoolType) GetValueSize() uint64 {
	return t.GetTypeSize()
}
func (t *BoolType) GetValue() any {
	if !t.hasValue {
		return nil
	}

	return t.value
}

type BytesType struct {
	data.BaseValueType
	value    []byte
	hasValue bool
}

var _ data.ValueType = &BytesType{}

func (t *BytesType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch tv := v.(type) {
	case []byte:
		t.value = tv
		t.hasValue = true
		return nil
	case string:
		t.value = []byte(tv)
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected []byte, got %T", data.ErrInvalidValue, v)
	}
}
func (t *BytesType) GetTypeKind() data.Kind {
	return data.KindBytes
}
func (t *BytesType) GetTypeName() string {
	return "BYTES"
}
func (t *BytesType) GetValueSize() uint64 {
	return t.GetTypeSize() * uint64(len(t.value))
}
func (t *BytesType) GetValue() any {
	if !t.hasValue {
		return nil
	}

	return t.value
}

type DateType struct {
	data.BaseValueType
	value    time.Time
	hasValue bool
}

var _ data.ValueType = &DateType{}

func (t *DateType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch tv := v.(type) {
	case time.Time:
		t.value = tv
		t.hasValue = true
		return nil
	case string:
		t.value, _ = time.Parse("2006-01-02", tv)
		t.hasValue = true
		return nil
	case []byte:
		t.value, _ = time.Parse("2006-01-02", string(tv))
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected time.Time, got %T", data.ErrInvalidValue, v)
	}
}
func (t *DateType) GetTypeKind() data.Kind {
	return data.KindTime
}
func (t *DateType) GetTypeName() string {
	return "DATE"
}
func (t *DateType) GetValueSize() uint64 {
	return t.GetTypeSize()
}
func (t *DateType) GetValue() any {
	if !t.hasValue {
		return nil
	}

	return t.value
}

// JSONBType is a type for jsonb.
type JSONBType struct {
	data.BaseValueType
	value    any
	hasValue bool
}

// Parse parses the value and stores it in the receiver.
func (t *JSONBType) Parse(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch tv := v.(type) {
	case []byte:
		err := json.Unmarshal(tv, &t.value)
		if err != nil {
			return err
		}
		t.hasValue = true
		return nil
	case string:
		err := json.Unmarshal([]byte(tv), &t.value)
		if err != nil {
			return err
		}
		t.hasValue = true
		return nil
	default:
		return fmt.Errorf("%v: expected []byte, got %T", data.ErrInvalidValue, v)
	}
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
func (t *JSONBType) GetValueSize() uint64 {
	bv, _ := json.Marshal(t.value)
	return t.GetTypeSize() * uint64(len(bv))
}

// GetValue returns the value stored in the receiver.
func (t *JSONBType) GetValue() any {
	return t.value
}

var ErrTypeNotFound = errors.New("type not found") // ErrTypeNotFound is returned when a type is not found.
// GetTypeFromName returns a type from its name.
func GetTypeFromName(name string) (data.Type, error) {
	switch {
	case mustMatchString("(?i)array", name):
		return &ArrayType{}, nil
	case mustMatchString("(?i)bit", name) || mustMatchString("(?i)varbit", name):
		return &BitType{}, nil
	case mustMatchString("(?i)bool", name) || mustMatchString("(?i)boolean", name):
		return &BoolType{}, nil
	case mustMatchString("(?i)bytes", name) || mustMatchString("(?i)blob", name) || mustMatchString("(?i)bytea", name):
		return &BytesType{}, nil
	case mustMatchString("(?i)date", name):
		return &DateType{}, nil
	case mustMatchString("(?i)jsonb", name) || mustMatchString("(?i)json", name):
		return &JSONBType{}, nil

	default:
		return nil, fmt.Errorf("%v: %s", ErrTypeNotFound, name)
	}
}

func mustMatchString(pattern, str string) bool {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		panic(err)
	}
	return matched
}
