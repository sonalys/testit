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
			if !fieldValue.IsNil() {
				continue
			}
			// No need to continue if the field is not related to mockery.
			if !fieldType.Implements(interfaceTypeTestify) && !fieldType.Implements(interfaceTypeT) {
				continue
			}
			target = initializeStructPointerField(fieldValue)
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

// initializeStructPointerField initializes a pointer field of a struct.
func initializeStructPointerField(field reflect.Value) any {
	fieldType := field.Type()
	// Create a new pointer to the field type.
	// This is necessary to bypass scope limitations.
	ptr := reflect.NewAt(fieldType, unsafe.Pointer(field.UnsafeAddr()))
	newValuePtr := reflect.New(fieldType.Elem())
	ptr.Elem().Set(newValuePtr)
	return newValuePtr.Interface()
}
