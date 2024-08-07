# TestIt

TestIt is a very simple, but powerful, testing framework.

It allows you to avoid repetition in the following:

* mock initialization
* pre-test cleanup
* panic avoidance
* test execution
* assertions

It also helps you build small blocks that can be re-utilized between test cases without sharing states. Examples:

* Setup a group of expected calls
* Shared assertions for response fields

## Usage

### tools/tools.go

Create the tools/tools.go to declare testit as a development only dependency

```go
//go:build tools

package tools

import (
	_ "github.com/sonalys/testit"
)

```

### A very simple example:

```go
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
```