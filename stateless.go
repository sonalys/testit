package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	stateless[D, R any] struct {
		pre, post []func(t *testing.T, d *D)
		setup     func(t *testing.T, d *D) R
	}

	StatelessTest[D, R any] struct {
		*require.Assertions
		T            *testing.T
		Dependencies *D
		Run          R
	}
)

// New is a function to create a new test setup.
func Stateless[Dependencies, RunFn any](setup func(t *testing.T, d *Dependencies) RunFn) stateless[Dependencies, RunFn] {
	return stateless[Dependencies, RunFn]{
		setup: setup,
	}
}

// PreCase is a function to add pre-case hooks.
func (th stateless[D, R]) PreCase(pre ...func(t *testing.T, d *D)) stateless[D, R] {
	th.pre = append(th.pre, pre...)
	return th
}

// PostCase is a function to add post-case hooks.
func (th stateless[D, R]) PostCase(post ...func(t *testing.T, d *D)) stateless[D, R] {
	th.post = append(th.post, post...)
	return th
}

// Expect is a function to create a test case with an empty test case.
func (th stateless[D, R]) Expect(steps ...func(t StatelessTest[D, R])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			dependencies := initializeMocks[D](t)
			run := th.setup(t, dependencies)
			for _, pre := range th.pre {
				pre(t, dependencies)
			}
			for _, step := range steps {
				step(StatelessTest[D, R]{
					T:            t,
					Assertions:   require.New(t),
					Dependencies: dependencies,
					Run:          run,
				})
			}
			for _, post := range th.post {
				post(t, dependencies)
			}
		})
	}
}
