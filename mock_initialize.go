package testit

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

// initializeMocks receives a struct with mock fields.
// All fields that implements interface { AssertExpectations(t mock.TestingT) bool } will be initialized.
// This function can set private and public fields automatically, and assert expectations.
// Fields that do not implement the interface will not be affected.
// Example:
//
//	type dependencies struct {
//		login *mocks.Login
//	}
//	initializeMocks(&dependencies{})
//
// Will initialize dependencies.login and call dependencies.login.AssertExpectations(t)
func initializeMocks[T any](t *testing.T) *T {
	var dependency T
	ptr := &dependency
	if t == nil {
		log.Fatal().Msgf("received a nil %T", t)
	}
	typeOf := reflect.TypeOf(ptr)
	if typeOf.Kind() != reflect.Pointer {
		log.Fatal().Msg("received a non-pointer value")
	}
	valueOf := reflect.ValueOf(ptr).Elem()
	if valueOf.Kind() != reflect.Struct {
		log.Fatal().Msgf("received a pointer to a non-struct value, %T", valueOf.Interface())
	}
	fieldLen := valueOf.NumField()
	type assertExpectationsTestify interface {
		AssertExpectations(t mock.TestingT) bool
	}
	type assertExpectationsT interface {
		AssertExpectations(t *testing.T) bool
	}

	interfaceTypeTestify := reflect.TypeOf((*assertExpectationsTestify)(nil)).Elem()
	interfaceTypeT := reflect.TypeOf((*assertExpectationsT)(nil)).Elem()

	for i := 0; i < fieldLen; i++ {
		var target any
		fieldValue := valueOf.Field(i)
		fieldType := fieldValue.Type()
		switch fieldValue.Kind() {
		// If the field is a pointer, we only set it if unitialized
		case reflect.Pointer:
			// No need to re-initialize this value if it's already set.
			// No need to continue if the field is not related to mockery.
			if !fieldValue.IsNil() || (!fieldType.Implements(interfaceTypeTestify) && !fieldType.Implements(interfaceTypeT)) {
				continue
			}
			newValue := reflect.New(fieldType.Elem())
			target = setValue(fieldValue, newValue)
		case reflect.Struct:
			ptr := reflect.NewAt(fieldType, unsafe.Pointer(fieldValue.UnsafeAddr()))
			target = ptr.Interface()
		}
		if impl, ok := target.(assertExpectationsTestify); ok {
			t.Cleanup(func() { impl.AssertExpectations(t) })
			continue
		}
		if impl, ok := target.(assertExpectationsT); ok {
			t.Cleanup(func() { impl.AssertExpectations(t) })
		}
	}
	return ptr
}

func setValue(dst reflect.Value, src reflect.Value) any {
	dstType := dst.Type()
	srcType := src.Type()
	if dstType != srcType {
		log.Fatal().Msgf("setProtectedValue: dst and src must have same type. got dst=%s src=%s", dstType, srcType)
	}
	if dst.CanSet() {
		dst.Set(src)
		return src.Interface()
	}
	if dst.Kind() == reflect.String {
		log.Fatal().Msg("cannot set a private string value. strings are immutable")
	}
	dstPtr := reflect.NewAt(dstType, unsafe.Pointer(dst.UnsafeAddr()))
	dstPtr.Elem().Set(src)
	return src.Interface()
}
