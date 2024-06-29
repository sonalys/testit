# TestIt

TestIt is a very simple, but powerful, testing framework.

It allows you to avoid repetition in mock initialization, pre-test cleanup, panic avoidance.

It also helps you build small behavior blocks that can be re-utilized between test cases without sharing states.

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

To start using, you can simply create

```go
package main_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test(t *testing.T) {
	type Dependencies struct{}
	type TestCase struct{}
	type TargetFunc func(id string) (any, error)
	type Test = testit.Test[Dependencies, TestCase, TargetFunc]

	targetFn := func(id string) (any, error) {
		return nil, nil
	}

	setup := testit.New(func(t *testing.T, d *Dependencies, tc *TestCase) TargetFunc {
		println("pre-run setup")
		return targetFn
	})

	additionalStep := func(t Test) {
		println("additional step")
	}

	t.Run("test case", setup.Case(
		TestCase{},
		additionalStep,
		func(t Test) {
			got, err := t.Run("id")
			println("after-run assertions")
			t.NoError(err)
			t.Nil(got)
		},
	))
}
```