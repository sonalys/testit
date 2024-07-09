package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	Hook[D any] struct {
		*assertions
		T            *testing.T
		Dependencies *D
		After        func()
	}

	setup[D any] struct {
		hooks []func(h *Hook[D])
	}

	DTest[D any] struct {
		*require.Assertions
		T            *testing.T
		Dependencies *D
	}
)

func Setup[D any](s func(t *Hook[D])) setup[D] {
	return setup[D]{
		hooks: []func(h *Hook[D]){s},
	}
}

// Hook is a function to add a hook to the test setup.
// The hook allows you to perform setup and teardown actions before and after the test case.
// You can run teardown actions by assigning a function to the After field of the CHook struct.
// Example:
//
//	setup := testit.SetupWithTestCase(func(t *testit.Hook[dependencies]) {
//		println("setup")
//		t.After = func() {
//			println("teardown")
//		}
//	})
func (th setup[D]) Hook(hook ...func(t *Hook[D])) setup[D] {
	th.hooks = append(th.hooks, hook...)
	return th
}

func (th setup[D]) Expect(steps ...func(t DTest[D])) func(*testing.T) {
	return func(t *testing.T) {
		NotPanics(t, func() {
			assertions := require.New(t)
			dependencies := initializeMocks[D](t)
			after := make([]func(), 0, len(th.hooks))
			for _, pre := range th.hooks {
				hook := &Hook[D]{
					assertions:   assertions,
					Dependencies: dependencies,
					T:            t,
				}
				pre(hook)
				if hook.After != nil {
					after = append(after, hook.After)
				}
			}
			for _, step := range steps {
				step(DTest[D]{
					T:            t,
					Assertions:   assertions,
					Dependencies: dependencies,
				})
			}
			for _, afterFunc := range after {
				afterFunc()
			}
		})
	}
}
