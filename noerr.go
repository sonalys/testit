package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func NoErr[T any](value T, err error) func(t *testing.T, msgAndArgs ...any) T {
	return func(t *testing.T, msgAndArgs ...any) T {
		require.NoError(t, err, msgAndArgs)
		return value
	}
}
