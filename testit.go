package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Setup[D, C, R any] struct {
	setup func(t *testing.T, d *D, tc *C) R
}

type Test[D, C, R any] struct {
	*require.Assertions
	T            *testing.T
	Dependencies *D
	Case         C
	Run          R
}

// New is a function to create a new test setup.
func New[Dependencies, TestCase, RunFn any](
	setup func(t *testing.T, d *Dependencies, tc *TestCase) RunFn,
) *Setup[Dependencies, TestCase, RunFn] {
	return &Setup[Dependencies, TestCase, RunFn]{
		setup: setup,
	}
}

// Case is a function to create a test case with the given test case.
func (th *Setup[D, C, R]) Case(tc C, steps ...func(t Test[D, C, R])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			var dependencies D
			if any(dependencies) != any(nil) {
				initializeMocks(t, &dependencies)
			}
			run := th.setup(t, &dependencies, &tc)
			for _, step := range steps {
				step(Test[D, C, R]{
					T:            t,
					Assertions:   require.New(t),
					Dependencies: &dependencies,
					Case:         tc,
					Run:          run,
				})
			}
		})
	}
}

// Expect is a function to create a test case with an empty test case.
func (th *Setup[D, C, R]) Expect(steps ...func(t Test[D, C, R])) func(*testing.T) {
	var tc C
	return th.Case(tc, steps...)
}
