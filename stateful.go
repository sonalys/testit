package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	stateful[D, C, R any] struct {
		pre, post []func(t *testing.T, d *D, tc *C)
		setup     func(t *testing.T, d *D, tc *C) R
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
func Stateful[Dependencies, TestCase, RunFn any](setup func(t *testing.T, d *Dependencies, tc *TestCase) RunFn) stateful[Dependencies, TestCase, RunFn] {
	return stateful[Dependencies, TestCase, RunFn]{
		setup: setup,
	}
}

// PreCase is a function to add pre-case hooks.
func (th stateful[D, C, R]) PreCase(pre ...func(t *testing.T, d *D, tc *C)) stateful[D, C, R] {
	th.pre = append(th.pre, pre...)
	return th
}

// PostCase is a function to add post-case hooks.
func (th stateful[D, C, R]) PostCase(post ...func(t *testing.T, d *D, tc *C)) stateful[D, C, R] {
	th.post = append(th.post, post...)
	return th
}

// Case is a function to create a test case with the given test case.
func (th stateful[D, C, R]) Case(tc C, steps ...func(t StatefulTest[D, C, R])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			dependencies := initializeMocks[D](t)
			run := th.setup(t, dependencies, &tc)
			for _, pre := range th.pre {
				pre(t, dependencies, &tc)
			}
			for _, step := range steps {
				step(StatefulTest[D, C, R]{
					T:            t,
					Assertions:   require.New(t),
					Dependencies: dependencies,
					Case:         &tc,
					Run:          run,
				})
			}
			for _, post := range th.post {
				post(t, dependencies, &tc)
			}
		})
	}
}

// Expect is a function to create a test case with an empty test case.
func (th stateful[D, C, R]) Expect(steps ...func(t StatefulTest[D, C, R])) func(*testing.T) {
	var tc C
	return th.Case(tc, steps...)
}
