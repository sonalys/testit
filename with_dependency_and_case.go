package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	CHook[D, C any] struct {
		*require.Assertions
		*testing.T
		Dependencies *D
		Case         *C
		After        func()
	}

	withDepAndCase[D, C any] struct {
		hooks []func(t *CHook[D, C])
	}

	DCTest[D, C any] struct {
		*require.Assertions
		T            *testing.T
		Dependencies *D
		Case         *C
	}
)

func SetupWithTestCase[Dependencies, TestCase any](s func(t *CHook[Dependencies, TestCase])) withDepAndCase[Dependencies, TestCase] {
	return withDepAndCase[Dependencies, TestCase]{
		hooks: []func(t *CHook[Dependencies, TestCase]){s},
	}
}

// Hook is a function to add a hook to the test setup.
// The hook allows you to perform setup and teardown actions before and after the test case.
// You can run teardown actions by assigning a function to the After field of the CHook struct.
// Example:
//
//	setup := testit.SetupWithTestCase(func(t *testit.CHook[dependencies, testCase]) {
//		println("setup")
//		t.After = func() {
//			println("teardown")
//		}
//	})
func (th withDepAndCase[D, C]) Hook(hooks ...func(t *CHook[D, C])) withDepAndCase[D, C] {
	th.hooks = append(th.hooks, hooks...)
	return th
}

func (th withDepAndCase[D, C]) Case(tc C, steps ...func(t DCTest[D, C])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			assertions := require.New(t)
			dependencies := initializeMocks[D](t)
			after := make([]func(), 0, len(th.hooks))
			for _, pre := range th.hooks {
				hook := &CHook[D, C]{
					Assertions:   assertions,
					Dependencies: dependencies,
					T:            t,
					Case:         &tc,
				}
				pre(hook)
				if hook.After != nil {
					after = append(after, hook.After)
				}
			}
			for _, step := range steps {
				step(DCTest[D, C]{
					T:            t,
					Assertions:   assertions,
					Dependencies: dependencies,
					Case:         &tc,
				})
			}
			for _, afterFunc := range after {
				afterFunc()
			}
		})
	}
}
