package testit

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type fakeMock struct {
	called bool
}

func (f *fakeMock) AssertExpectations(t mock.TestingT) bool {
	f.called = true
	return true
}

func TestInitializeMocks(t *testing.T) {
	t.Run("with pointer", func(t *testing.T) {
		type test struct {
			mock *fakeMock
		}
		testStruct := &test{}
		require.Nil(t, testStruct.mock)
		t.Run("limited t scope", func(t *testing.T) {
			initializeMocks(t, testStruct)
		})
		require.NotNil(t, testStruct.mock)
		require.True(t, testStruct.mock.called)
	})

	t.Run("with struct", func(t *testing.T) {
		type test struct {
			mock fakeMock
		}
		testStruct := &test{}
		t.Run("limited t scope", func(t *testing.T) {
			initializeMocks(t, testStruct)
		})
		require.True(t, testStruct.mock.called)
	})
}
