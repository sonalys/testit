package testit

import (
	"errors"
	"reflect"
)

var (
	ErrNoMatchType     = errors.New("no match type")
	ErrNoPointer       = errors.New("must be interface")
	ErrInvalidArgument = errors.New("invalid arguments")
)

func deepCopy(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Interface:
		value := src.Elem()
		if !value.IsValid() {
			return
		}
		newValue := reflect.New(value.Type()).Elem()
		deepCopy(newValue, value)
		dst.Set(newValue)
	case reflect.Ptr:
		value := src.Elem()
		if !value.IsValid() {
			return
		}
		dst.Set(reflect.New(value.Type()))
		deepCopy(dst.Elem(), value)
	case reflect.Map:
		dst.Set(reflect.MakeMap(src.Type()))
		keys := src.MapKeys()
		for _, key := range keys {
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			deepCopy(newValue, value)
			dst.SetMapIndex(key, newValue)
		}
	case reflect.Slice:
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			deepCopy(dst.Index(i), src.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			deepCopy(dst.Field(i), src.Field(i))
		}
	default:
		setValue(dst, src)
	}
}

// DeepClone is a function to deep clone a value.
// It can clone public and private values, but not private strings, as they are immutable.
func DeepClone[T any](v T) T {
	dst := reflect.New(reflect.TypeOf(v)).Elem()
	deepCopy(dst, reflect.ValueOf(v))
	return dst.Interface().(T)
}
