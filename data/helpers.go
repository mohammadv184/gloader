package data

import "reflect"

func IsCompatibleWithType(t1, t2 ValueType) bool {
	if _, err := t1.To(t2); err != nil {
		return false
	}
	return true
}

func GetNewType(t Type) Type {
	// dereference pointer
	if reflect.TypeOf(t).Kind() == reflect.Ptr {
		t = reflect.TypeOf(t).Elem().(Type)
	}

	return reflect.New(reflect.TypeOf(t)).Interface().(Type)
}
