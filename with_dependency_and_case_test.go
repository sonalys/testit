package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_Stateful(t *testing.T) {
	type dependencies struct{}
	type testCase struct{}

	setup := testit.SetupWithTestCase(func(t *testit.CHook[dependencies, testCase]) {
		println("setup")
	})

	withHook := setup.Hook(func(t *testit.CHook[dependencies, testCase]) {
		println("hook pre-case")
		t.After = func() {
			println("hook after-case")
		}
	})

	t.Run("test", withHook.Case(testCase{}, func(t testit.DCTest[dependencies, testCase]) {
		println("test")
	}))
}
