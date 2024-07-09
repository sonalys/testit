package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_Example(t *testing.T) {
	type dependency struct {
		// mock *mock.Mock // Mocks are automatically initialized
	}
	type tc struct{}

	setup := testit.SetupWithTestCase(func(t *testit.CHook[dependency, tc]) {
		// setup, pre-cleanup.
		// t.Dependencies.mock.EXPECT().DoSomething()
		t.After = func() {
			// teardown
		}
	})

	withHook := setup.Hook(func(t *testit.CHook[dependency, tc]) {
		// additional hook
		t.After = func() {
			// additional teardown
		}
	})

	additionalStep := func(t testit.DCTest[dependency, tc]) {}

	testfunc := func() (int, error) {
		return 0, nil
	}

	t.Run("example", withHook.Case(tc{},
		// any additional steps can be added manually as well.
		additionalStep,
		func(t testit.DCTest[dependency, tc]) {
			// test execution
			var err error
			t.NoError(err)
			t.FailNow("failing test")

			value := testit.NoErr(testfunc())
			t.Equal(0, value)
		},
	))
}
