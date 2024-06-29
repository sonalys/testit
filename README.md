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

### A very simple example for Stateful and Stateless tests:

```go
package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_Stateful(t *testing.T) {
	type Dependencies struct{ // Mocks, Infrastructure and other injectables.
		Mock // Mock is automatically initialized when it implements interface { AssertExpectations(*testing.T) }
		DB
	}
	type TestCase struct{ ID int } // Shared state between pre-run, run and post-run.
	type TargetFunc = func(id int) (any, error) // The execution plan for the tested behavior.
	type Test = testit.StatefulTest[Dependencies, TestCase, TargetFunc] // A simple signature alias to make code more readable.

	targetFn := func(id int) (any, error) {
		return nil, nil
	}

	setup := testit.Stateful(func(t *testing.T, d *Dependencies, tc *TestCase) TargetFunc {
		d.DB = newDB()
		d.DB.Clean()
		println("pre-run setup")
		return targetFn
	})

	additionalStep := func(t Test) { // Shared call expectations and other initializations that are used only by certain cases.
		t.Dependencies.Mock.REQUIRE().Func().Return()
		println("additional step")
	}

	t.Run("test case", setup.Case(
		TestCase{
			ID: 1,
		},
		additionalStep,
		func(t Test) {
			got, err := t.Run(t.Case.ID)
			println("after-run assertions")
			t.NoError(err)
			t.Nil(got)
		},
	))
}

func Test_Stateless(t *testing.T) {
	type Dependencies struct{}
	type TargetFunc = func(id string) (any, error)

	setup := testit.Stateless(func(t *testing.T, d *Dependencies) TargetFunc {
		println("pre-run setup")
		return func(id string) (any, error) {
			return nil, nil
		}
	})

	additionalStep := func(t testit.StatelessTest[Dependencies, TargetFunc]) {
		println("additional step")
	}

	t.Run("test case", setup.Expect(
		additionalStep,
		func(t testit.StatelessTest[Dependencies, TargetFunc]) {
			got, err := t.Run("id")
			println("after-run assertions")
			t.NoError(err)
			t.Nil(got)
		},
	))
}
```