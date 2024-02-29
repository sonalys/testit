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
func initializeMocks(t *testing.T, m any) {
	if t == nil {
		log.Fatal().Msgf("received a nil %T", t)
	}
	typeOf := reflect.TypeOf(m)
	if typeOf.Kind() != reflect.Pointer {
		log.Fatal().Msg("received a non-pointer value")
	}
	valueOf := reflect.ValueOf(m).Elem()
	if valueOf.Kind() != reflect.Struct {
		log.Fatal().Msgf("received a pointer to a non-struct value, %T", valueOf.Interface())
	}
	fieldLen := valueOf.NumField()
	type mockInterface interface {
		AssertExpectations(t mock.TestingT) bool
	}
	interfaceType := reflect.TypeOf((*mockInterface)(nil)).Elem()
	for i := 0; i < fieldLen; i++ {
		value := valueOf.Field(i)
		// bypass memory protection of unexported fields.
		newField := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr()))
		var impl mockInterface
		var ok bool
		switch value.Kind() {
		// If the field is a pointer, we only set it if unitialized
		case reflect.Pointer:
			// No need to continue if the field is not related to mockery.
			// No need to re-initialize this value if it's already set.
			if !value.IsNil() || !newField.Elem().Type().Implements(interfaceType) {
				continue
			}
			newValue := reflect.New(value.Type().Elem())
			newField.Elem().Set(newValue)
			// If the field is a struct, it's already initialized.
			impl, ok = newValue.Interface().(mockInterface)
		case reflect.Struct:
			// No need to continue if the field is not related to mockery.
			if !newField.Type().Implements(interfaceType) {
				continue
			}
			value := newField.Interface()
			impl, ok = value.(mockInterface)
		}
		if ok {
			t.Cleanup(func() { impl.AssertExpectations(t) })
		}
	}
}
