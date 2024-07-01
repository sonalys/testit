package testit

import (
	"testing"
)

type (
	runlessStateful[D, C any] func(t *testing.T, d *D, tc *C)
	runlessStateless[D any]   func(t *testing.T, d *D)
)

func Func[D, C, R any](run R, setup runlessStateful[D, C]) stateful[D, C, R] {
	return Stateful(func(t *testing.T, d *D, tc *C) R {
		setup(t, d, tc)
		return run
	})
}

func FuncStateless[D, R any](run R, setup runlessStateless[D]) stateless[D, R] {
	return Stateless(func(t *testing.T, d *D) R {
		setup(t, d)
		return run
	})
}
