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
		value := valueOf.Field(i)
		// bypass memory protection of unexported fields.
		newField := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr()))
		switch value.Kind() {
		// If the field is a pointer, we only set it if unitialized
		case reflect.Pointer:
			// No need to continue if the field is not related to mockery.
			// No need to re-initialize this value if it's already set.
			if !value.IsNil() ||
				!newField.Elem().Type().Implements(interfaceTypeTestify) ||
				!newField.Elem().Type().Implements(interfaceTypeT) {
				continue
			}
			newValue := reflect.New(value.Type().Elem())
			newField.Elem().Set(newValue)
			// If the field is a struct, it's already initialized.
			if impl, ok := newValue.Interface().(assertExpectationsTestify); ok {
				t.Cleanup(func() { impl.AssertExpectations(t) })
			} else if impl, ok := newValue.Interface().(assertExpectationsT); ok {
				t.Cleanup(func() { impl.AssertExpectations(t) })
			}
		case reflect.Struct:
			value := newField.Interface()
			if impl, ok := value.(assertExpectationsTestify); ok {
				t.Cleanup(func() { impl.AssertExpectations(t) })
			} else if impl, ok := value.(assertExpectationsT); ok {
				t.Cleanup(func() { impl.AssertExpectations(t) })
			}
		}
	}
	return ptr
}
