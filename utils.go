package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type panicError struct {
	err error
}

func NoErr[V any](value V, err error) V {
	if err != nil {
		panic(panicError{err: err})
	}
	return value
}

func NotPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		switch r := r.(type) {
		case panicError:
			require.NoError(t, r.err)
		}
	}()
	f()
}
