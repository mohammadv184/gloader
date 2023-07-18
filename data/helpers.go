package data

import "reflect"

func IsCompatibleWithType(t1, t2 ValueType) bool {
	if _, err := t1.To(t2); err != nil {
		return false
	}
	return true
}

func GetNewValueTypeAs(t Type) ValueType {
	// dereference pointer
	if reflect.TypeOf(t).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(t).Elem().Type()).Interface().(ValueType)
	}

	return reflect.New(reflect.TypeOf(t)).Interface().(ValueType)
}
