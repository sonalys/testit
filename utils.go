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
	return fmt.Sprintf("%#v", stack.Trace().TrimBelow(stack.Caller(5)))
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
