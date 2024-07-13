package testit

import (
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/go-stack/stack"
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

func getStack() string {
	if !BetterStack {
		return string(debug.Stack())
	}
	// 3 because we ignore this function, the NotPanics and the panic.go internal.
	return fmt.Sprintf("%#v", stack.Trace().TrimAbove(stack.Caller(3)))
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
		default:
			require.Fail(t, "unexpected panic", "%v: %s", r, getStack())
		}
	}()
	f()
}
