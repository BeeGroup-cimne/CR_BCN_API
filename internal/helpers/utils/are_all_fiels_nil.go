package utils

import "reflect"

func AreAllFieldsNil[T any](input T) bool {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Struct {
		panic("areAllFieldsNil: input must be a struct")
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			return false
		}
	}
	return true
}
