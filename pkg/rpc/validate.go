package rpc

import (
	"fmt"
	"reflect"
	"strings"
)

var strict_check_tag = "strict_check"

func Validate(elem interface{}) error {
	t := reflect.TypeOf(elem)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}

		if strings.Contains(tag, strict_check_tag) {
			value := reflect.ValueOf(elem).Field(i)
			if isEmpty(value) {
				fieldName := field.Tag.Get("json")
				return fmt.Errorf("field '%s' is required but missing or null in JSON", strings.TrimPrefix(strings.Split(fieldName, ",")[0], "json:"))
			}
		}

	}
	return nil
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(v, reflect.Zero(v.Type()))
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.String:
		return v.Len() == 0
	}
	return false
}
