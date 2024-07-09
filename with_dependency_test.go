package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_Stateless(t *testing.T) {
	type dependencies struct{}

	setup := testit.Setup[dependencies](func(t *testit.Hook[dependencies]) {
		println("before")
		t.After = func() {
			println("after")
		}
	})

	withHook := setup.Hook(func(t *testit.Hook[dependencies]) {
		println("before")
		t.After = func() {
			println("after")
		}
	})

	t.Run("test", withHook.Expect(func(t testit.DTest[dependencies]) {
		println("test")
	}))
}
