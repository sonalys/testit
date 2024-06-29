package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	StatefulSetup[D, C, R any] struct {
		setup func(t *testing.T, d *D, tc *C) R
	}

	StatefulTest[D, C, R any] struct {
		*require.Assertions
		T            *testing.T
		Dependencies *D
		Case         *C
		Run          R
	}
)

// New is a function to create a new test setup.
func Stateful[Dependencies, TestCase, RunFn any](setup func(t *testing.T, d *Dependencies, tc *TestCase) RunFn) *StatefulSetup[Dependencies, TestCase, RunFn] {
	return &StatefulSetup[Dependencies, TestCase, RunFn]{
		setup: setup,
	}
}

// Case is a function to create a test case with the given test case.
func (th *StatefulSetup[D, C, R]) Case(tc C, steps ...func(t StatefulTest[D, C, R])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			dependencies := initializeMocks[D](t)
			run := th.setup(t, dependencies, &tc)
			for _, step := range steps {
				step(StatefulTest[D, C, R]{
					T:            t,
					Assertions:   require.New(t),
					Dependencies: dependencies,
					Case:         &tc,
					Run:          run,
				})
			}
		})
	}
}

// Expect is a function to create a test case with an empty test case.
func (th *StatefulSetup[D, C, R]) Expect(steps ...func(t StatefulTest[D, C, R])) func(*testing.T) {
	var tc C
	return th.Case(tc, steps...)
}
