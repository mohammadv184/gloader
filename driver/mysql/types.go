package mysql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/mohammadv184/gloader/data"
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

// GetValueSize returns the size of the value in bytes.
func (t *CharType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *CharType) GetValue() any {
	return t.value
}

// VarCharType is a type for varchar.
type VarCharType struct {
	data.BaseValueType
	value string
}

// Parse parses the value and stores it in the receiver.
func (t *VarCharType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = fmt.Sprintf("%s", v)
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *VarCharType) GetTypeKind() data.Kind {
	return data.KindString
}

// GetTypeName returns the name of the type.
func (t *VarCharType) GetTypeName() string {
	return "VARCHAR"
}

// GetValueSize returns the size of the value in bytes.
func (t *VarCharType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *VarCharType) GetValue() any {
	return t.value
}

type TextType struct {
	data.BaseValueType
	value string
}

// Parse parses the value and stores it in the receiver.
func (t *TextType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = fmt.Sprintf("%s", v)
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *TextType) GetTypeKind() data.Kind {
	return data.KindString
}

// GetTypeName returns the name of the type.
func (t *TextType) GetTypeName() string {
	return "TEXT"
}

// GetValueSize returns the size of the value in bytes.
func (t *TextType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *TextType) GetValue() any {
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

// GetValueSize returns the size of the value in bytes.
func (t *SmallIntType) GetValueSize() uint64 {
	return 2
}

// GetValue returns the value stored in the receiver.
func (t *SmallIntType) GetValue() any {
	return t.value
}

type IntType struct {
	data.BaseValueType
	value int32
}

// Parse parses the value and stores it in the receiver.
func (t *IntType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case int8:
		t.value = int32(v.(int8))
		return nil
	case int16:
		t.value = int32(v.(int16))
		return nil
	case int32:
		t.value = v.(int32)
		return nil
	case int64:
		t.value = int32(v.(int64))
		return nil
	case uint8:
		t.value = int32(v.(uint8))
		return nil
	case uint16:
		t.value = int32(v.(uint16))
		return nil
	case uint32:
		t.value = int32(v.(uint32))
		return nil
	case uint64:
		t.value = int32(v.(uint64))
		return nil
	case []byte:
		v, err := strconv.ParseInt(string(v.([]byte)), 10, 32)
		if err != nil {
			return err
		}
		t.value = int32(v)
		return nil

	default:
		return fmt.Errorf("%v: expected int32, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *IntType) GetTypeKind() data.Kind {
	return data.KindInt32
}

// GetTypeName returns the name of the type.
func (t *IntType) GetTypeName() string {
	return "INT"
}

// GetValueSize returns the size of the value in bytes.
func (t *IntType) GetValueSize() uint64 {
	return 4
}

// GetValue returns the value stored in the receiver.
func (t *IntType) GetValue() any {
	return t.value
}

type TinyIntType struct {
	data.BaseValueType
	value bool
}

// Parse parses the value and stores it in the receiver.
func (t *TinyIntType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case bool:
		t.value = v.(bool)
		return nil
	case int8:
		t.value = v.(int8) != 0
		return nil
	case int16:
		t.value = v.(int16) != 0
		return nil
	case int32:
		t.value = v.(int32) != 0
		return nil
	case int64:
		t.value = v.(int64) != 0
		return nil
	case uint8:
		t.value = v.(uint8) != 0
		return nil
	case uint16:
		t.value = v.(uint16) != 0
		return nil
	case uint32:
		t.value = v.(uint32) != 0
		return nil
	case uint64:
		t.value = v.(uint64) != 0
		return nil
	case []byte:
		v, err := strconv.ParseInt(string(v.([]byte)), 10, 8)
		if err != nil {
			return err
		}
		t.value = v != 0
		return nil

	default:
		return fmt.Errorf("%v: expected bool, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *TinyIntType) GetTypeKind() data.Kind {
	return data.KindBool
}

// GetTypeName returns the name of the type.
func (t *TinyIntType) GetTypeName() string {
	return "TINYINT"
}

// GetValueSize returns the size of the value in bytes.
func (t *TinyIntType) GetValueSize() uint64 {
	return 1
}

// GetValue returns the value stored in the receiver.
func (t *TinyIntType) GetValue() any {
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

// GetValueSize returns the size of the value in bytes.
func (t *BigIntType) GetValueSize() uint64 {
	return 8
}

// GetValue returns the value stored in the receiver.
func (t *BigIntType) GetValue() any {
	return t.value
}

type DecimalType struct {
	data.BaseValueType
	value float64
}

// Parse parses the value and stores it in the receiver.
func (t *DecimalType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch v.(type) {
	case int8:
		t.value = float64(v.(int8))
		return nil
	case int16:
		t.value = float64(v.(int16))
		return nil
	case int32:
		t.value = float64(v.(int32))
		return nil
	case int64:
		t.value = float64(v.(int64))
		return nil
	case uint8:
		t.value = float64(v.(uint8))
		return nil
	case uint16:
		t.value = float64(v.(uint16))
		return nil
	case uint32:
		t.value = float64(v.(uint32))
		return nil
	case uint64:
		t.value = float64(v.(uint64))
		return nil
	case float32:
		t.value = float64(v.(float32))
		return nil
	case float64:
		t.value = v.(float64)
		return nil
	case []byte:
		v, err := strconv.ParseFloat(string(v.([]byte)), 64)
		if err != nil {
			return err
		}
		t.value = v
		return nil
	default:
		return fmt.Errorf("%v: expected float64, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *DecimalType) GetTypeKind() data.Kind {
	return data.KindFloat64
}

// GetTypeName returns the name of the type.
func (t *DecimalType) GetTypeName() string {
	return "DECIMAL"
}

// GetValueSize returns the size of the value in bytes.
func (t *DecimalType) GetValueSize() uint64 {
	return 8
}

// GetValue returns the value stored in the receiver.
func (t *DecimalType) GetValue() any {
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

// GetValueSize returns the size of the value in bytes.
func (t *LongBlobType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *LongBlobType) GetValue() any {
	return t.value
}

type MediumBlobType struct {
	data.BaseValueType
	value []byte
}

// Parse parses the value and stores it in the receiver.
func (t *MediumBlobType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = []byte(fmt.Sprintf("%s", v))
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *MediumBlobType) GetTypeKind() data.Kind {
	return data.KindBytes
}

// GetTypeName returns the name of the type.
func (t *MediumBlobType) GetTypeName() string {
	return "MEDIUMBLOB"
}

// GetValueSize returns the size of the value in bytes.
func (t *MediumBlobType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *MediumBlobType) GetValue() any {
	return t.value
}

type EnumType struct {
	data.BaseValueType
	value string
}

// Parse parses the value and stores it in the receiver.
func (t *EnumType) Parse(v any) error {
	if reflect.TypeOf(v).Kind() == reflect.Pointer {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	t.value = fmt.Sprintf("%s", v)
	return nil
}

// GetTypeKind returns the kind of the type.
func (t *EnumType) GetTypeKind() data.Kind {
	return data.KindString
}

// GetTypeName returns the name of the type.
func (t *EnumType) GetTypeName() string {
	return "ENUM"
}

// GetValueSize returns the size of the value in bytes.
func (t *EnumType) GetValueSize() uint64 {
	return uint64(len(t.value))
}

// GetValue returns the value stored in the receiver.
func (t *EnumType) GetValue() any {
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
		tm, err := time.Parse("2006-01-02 15:04:05.999999999", string(v.([]byte)))
		if err != nil {
			return err
		}
		t.value = tm
		return nil
	case string:
		tm, err := time.Parse("2006-01-02 15:04:05.999999999", v.(string))
		if err != nil {
			return err
		}
		t.value = tm
		return nil
	default:
		return fmt.Errorf("%v: expected time.Time, got %T", data.ErrInvalidValue, v)
	}
}

// GetTypeKind returns the kind of the type.
func (t *DateTimeType) GetTypeKind() data.Kind {
	return data.KindTime
}

// GetTypeName returns the name of the type.
func (t *DateTimeType) GetTypeName() string {
	return "DATETIME"
}

// GetValueSize returns the size of the value in bytes.
func (t *DateTimeType) GetValueSize() uint64 {
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
	case mustMatchString("(?i)enum", name):
		return &EnumType{}, nil
	case mustMatchString("(?i)mediumblob", name):
		return &MediumBlobType{}, nil
	case mustMatchString("(?i)decimal", name):
		return &DecimalType{}, nil
	case mustMatchString("(?i)int", name):
		return &IntType{}, nil
	case mustMatchString("(?i)varchar", name):
		return &VarCharType{}, nil
	case mustMatchString("(?i)text", name):
		return &TextType{}, nil
	case mustMatchString("(?i)tinyint", name):
		return &TinyIntType{}, nil

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
