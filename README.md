# TestIt

TestIt is a very simple, but powerful, testing framework.

It allows you to avoid repetition in mock initialization, pre-test cleanup, panic avoidance.

It also helps you build small behavior blocks that can be re-utilized between test cases without sharing states.

Also check out [github.com/sonalys/fake](https://github.com/sonalys/fake). A type-safe go mock generator.

## Usage

To start using, you can simply create

```go
package main

import (
	"testing"

	"github.com/sonalys/testit"
	"github.com/stretchr/testify/require"
	// ...
)

func Test_Example(test *testing.T) {
	type Dependencies struct {
		m  mocks.SomeMock
		db db.Database
	}
	type Case struct {
		ID string
	}
	type Result struct {
		u   User
		err error
	}

	t := testit.New(test, func(t *testing.T, d *Dependencies, tc *Case) func() (r Result) {
		// Pre-run initialization.
		d.m = mocks.Initialize()
		m.db = database.NewConn(&d.m)
		require.NoError(t, m.db.Cleanup())
		// Test-case run function that will execute the run behavior.
		return func() (r Result) {
			r.u, r.err = m.db.FindUser(tc.ID)
			return
		}
	})

	t.Run("case name", t.Case(
		Case{
			ID: "userID",
		},
   // You can pass multiple functions at each case.
		func(t *testing.T, d *Dependencies, tc *Case, run func() Result) {
			// Logic executed before running the test case.
			d.m.OnIsCached(func(id string) bool {
				return false
			})
			// Executes the common run function from the initial declaration.
			result := run()
			// Assertions after the test is ran.
			require.NoError(t, result.err)
			require.NotNil(t, result.u)
		},
	))
}

```