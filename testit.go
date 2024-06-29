package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Setup[D, C, R any] struct {
	setup func(t *testing.T, d *D, tc *C) R
}

type Config[D, C, R any] struct {
	*require.Assertions
	Dependencies *D
	Case         C
	Run          R
}

func New[Dependencies, TestCase, RunFn any](
	setup func(t *testing.T, d *Dependencies, tc *TestCase) RunFn,
) *Setup[Dependencies, TestCase, RunFn] {
	return &Setup[Dependencies, TestCase, RunFn]{
		setup: setup,
	}
}

func (th *Setup[D, C, R]) Case(tc C, steps ...func(t Config[D, C, R])) func(*testing.T) {
	return func(t *testing.T) {
		require.NotPanics(t, func() {
			var dependencies D
			if any(dependencies) != any(nil) {
				initializeMocks(t, &dependencies)
			}
			run := th.setup(t, &dependencies, &tc)
			for _, step := range steps {
				step(Config[D, C, R]{
					Assertions:   require.New(t),
					Dependencies: &dependencies,
					Case:         tc,
					Run:          run,
				})
			}
		})
	}
}
