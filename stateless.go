package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	StatelessSetup[D, R any] struct {
		setup func(t *testing.T, d *D) R
	}

	StatelessTest[D, R any] struct {
		*require.Assertions
		T            *testing.T
		Dependencies *D
		Run          R
	}
)

// New is a function to create a new test setup.
func Stateless[Dependencies, RunFn any](
	setup func(t *testing.T, d *Dependencies) RunFn,
) *StatelessSetup[Dependencies, RunFn] {
	return &StatelessSetup[Dependencies, RunFn]{
		setup: setup,
	}
}

// Expect is a function to create a test case with an empty test case.
func (th *StatelessSetup[D, R]) Expect(steps ...func(t StatelessTest[D, R])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			dependencies := initializeMocks[D](t)
			run := th.setup(t, dependencies)
			for _, step := range steps {
				step(StatelessTest[D, R]{
					T:            t,
					Assertions:   require.New(t),
					Dependencies: dependencies,
					Run:          run,
				})
			}
		})
	}
}
