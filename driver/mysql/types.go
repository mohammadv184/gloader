package mysql

import (
	"errors"
	"fmt"
	"gloader/data"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CharType struct {
	data.BaseValueType
	value string
}

func (t *CharType) Parse(v any) error {

	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = fmt.Sprintf("%s", v)
	return nil
}
func (t *CharType) GetTypeKind() data.Kind {
	return data.KindString
}
func (t *CharType) GetTypeName() string {
	return "CHAR"
}
func (t *CharType) GetTypeSize() uint64 {
	return uint64(len(t.value))
}
func (t *CharType) GetValue() any {
	return t.value
}

type SmallIntType struct {
	data.BaseValueType
	value int16
}

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
func (t *SmallIntType) GetTypeKind() data.Kind {
	return data.KindInt16
}
func (t *SmallIntType) GetTypeName() string {
	return "SMALLINT"
}
func (t *SmallIntType) GetTypeSize() uint64 {
	return 2
}
func (t *SmallIntType) GetValue() any {
	return t.value
}

type BigIntType struct {
	data.BaseValueType
	value int64
}

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
func (t *BigIntType) GetTypeKind() data.Kind {
	return data.KindInt64
}
func (t *BigIntType) GetTypeName() string {
	return "BIGINT"
}
func (t *BigIntType) GetTypeSize() uint64 {
	return 8
}
func (t *BigIntType) GetValue() any {
	return t.value
}

type LongBlobType struct {
	data.BaseValueType
	value []byte
}

func (t *LongBlobType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = []byte(fmt.Sprintf("%s", v))
	return nil
}
func (t *LongBlobType) GetTypeKind() data.Kind {
	return data.KindBytes
}
func (t *LongBlobType) GetTypeName() string {
	return "LONGBLOB"
}
func (t *LongBlobType) GetTypeSize() uint64 {
	return uint64(len(t.value))
}
func (t *LongBlobType) GetValue() any {
	return t.value
}

type DateTimeType struct {
	data.BaseValueType
	value time.Time
}

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
func (t *DateTimeType) GetTypeKind() data.Kind {
	return data.KindDateTime
}
func (t *DateTimeType) GetTypeName() string {
	return "DATETIME"
}
func (t *DateTimeType) GetTypeSize() uint64 {
	return 8
}
func (t *DateTimeType) GetValue() any {
	return t.value
}

var ErrTypeNotFound = errors.New("type not found")

func GetTypeFromName(name string) (data.Type, error) {
	name = strings.ToUpper(regexp.MustCompile(`\(.*\).*`).ReplaceAllString(name, ""))

	switch name {
	case "CHAR":
		return &CharType{}, nil
	case "SMALLINT":
		return &SmallIntType{}, nil
	case "BIGINT":
		return &BigIntType{}, nil
	case "LONGBLOB":
		return &LongBlobType{}, nil
	case "DATETIME":
		return &DateTimeType{}, nil
	default:
		return nil, ErrTypeNotFound
	}
}
