package mysql

import (
	"errors"
	"fmt"
	"gloader/data"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

// CharType is a type for char.
type CharType struct {
	data.BaseValueType
	value string
}

// Parse parses the value and stores it in the receiver.
func (t *CharType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = fmt.Sprintf("%s", v)
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *CharType) GetTypeKind() data.Kind {
	return data.KindString
}

// GetTypeName returns the name of the type.
func (t *CharType) GetTypeName() string {
	return "CHAR"
}

// GetTypeSize returns the size of the value in bytes.
func (t *CharType) GetTypeSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *CharType) GetValue() any {
	return t.value
}

// SmallIntType is a type for smallint.
type SmallIntType struct {
	data.BaseValueType
	value int16
}

// Parse parses the value and stores it in the receiver.
func (t *SmallIntType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case int8:
		t.value = int16(v.(int8))
		return nil
	case int16:
		t.value = v.(int16)
		return nil
	case int32:
		t.value = int16(v.(int32))
		return nil
	case int64:
		t.value = int16(v.(int64))
		return nil
	case uint8:
		t.value = int16(v.(uint8))
		return nil
	case uint16:
		t.value = int16(v.(uint16))
		return nil
	case uint32:
		t.value = int16(v.(uint32))
		return nil
	case uint64:
		t.value = int16(v.(uint64))
		return nil
	case []byte:
		v, err := strconv.ParseInt(string(v.([]byte)), 10, 16)
		if err != nil {
			return err
		}
		t.value = int16(v)
		return nil

	default:
		return fmt.Errorf("%v: expected int16, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *SmallIntType) GetTypeKind() data.Kind {
	return data.KindInt16
}

// GetTypeName returns the name of the type.
func (t *SmallIntType) GetTypeName() string {
	return "SMALLINT"
}

// GetTypeSize returns the size of the value in bytes.
func (t *SmallIntType) GetTypeSize() uint64 {
	return 2
}

// GetValue returns the value stored in the receiver.
func (t *SmallIntType) GetValue() any {
	return t.value
}

// BigIntType is a type for bigint.
type BigIntType struct {
	data.BaseValueType
	value int64
}

// Parse parses the value and stores it in the receiver.
func (t *BigIntType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
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
	case uint8:
		t.value = int64(v.(uint8))
		return nil
	case uint16:
		t.value = int64(v.(uint16))
		return nil
	case uint32:
		t.value = int64(v.(uint32))
		return nil
	case uint64:
		t.value = int64(v.(uint64))
		return nil
	case []byte:
		v, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
		if err != nil {
			return err
		}
		t.value = v
		return nil
	default:
		return fmt.Errorf("%v: expected int64, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *BigIntType) GetTypeKind() data.Kind {
	return data.KindInt64
}

// GetTypeName returns the name of the type.
func (t *BigIntType) GetTypeName() string {
	return "BIGINT"
}

// GetTypeSize returns the size of the value in bytes.
func (t *BigIntType) GetTypeSize() uint64 {
	return 8
}

// GetValue returns the value stored in the receiver.
func (t *BigIntType) GetValue() any {
	return t.value
}

// LongBlobType is a type for longblob.
type LongBlobType struct {
	data.BaseValueType
	value []byte
}

// Parse parses the value and stores it in the receiver.
func (t *LongBlobType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = []byte(fmt.Sprintf("%s", v))
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *LongBlobType) GetTypeKind() data.Kind {
	return data.KindBytes
}

// GetTypeName returns the name of the type.
func (t *LongBlobType) GetTypeName() string {
	return "LONGBLOB"
}

// GetTypeSize returns the size of the value in bytes.
func (t *LongBlobType) GetTypeSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *LongBlobType) GetValue() any {
	return t.value
}

// DateTimeType is a type for datetime.
type DateTimeType struct {
	data.BaseValueType
	value time.Time
}

// Parse parses the value and stores it in the receiver.
func (t *DateTimeType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case time.Time:
		t.value = v.(time.Time)
		return nil
	case []byte:
		v, err := time.Parse("2006-01-02 15:04:05", string(v.([]byte)))
		if err != nil {
			return err
		}
		t.value = v
		return nil
	default:
		return fmt.Errorf("%v: expected time.Time, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *DateTimeType) GetTypeKind() data.Kind {
	return data.KindDateTime
}

// GetTypeName returns the name of the type.
func (t *DateTimeType) GetTypeName() string {
	return "DATETIME"
}

// GetTypeSize returns the size of the value in bytes.
func (t *DateTimeType) GetTypeSize() uint64 {
	return 8
}

// GetValue returns the value stored in the receiver.
func (t *DateTimeType) GetValue() any {
	return t.value
}

var ErrTypeNotFound = errors.New("type not found") // ErrTypeNotFound is returned when a type is not found.
// GetTypeFromName returns a type from its name.
func GetTypeFromName(name string) (data.Type, error) {
	switch {
	case mustMatchString("(?i)char", name):
		return &CharType{}, nil
	case mustMatchString("(?i)smallint", name):
		return &SmallIntType{}, nil
	case mustMatchString("(?i)bigint", name):
		return &BigIntType{}, nil
	case mustMatchString("(?i)longblob", name) || mustMatchString("(?i)blob", name):
		return &LongBlobType{}, nil
	case mustMatchString("(?i)datetime", name):
		return &DateTimeType{}, nil
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
