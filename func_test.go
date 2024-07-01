package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_FuncStateful(t *testing.T) {
	method := func(id string) (any, error) {
		return nil, nil
	}

	type dependencies struct {
	}

	type testCase struct {
		id string
	}

	setup := testit.Func(method, func(t *testing.T, d *dependencies, tc *testCase) {
	})

	t.Run("test", setup.Case(
		testCase{id: "1"},
		func(t testit.StatefulTest[dependencies, testCase, func(string) (any, error)]) {

		},
	))
}

func Test_FuncStateless(t *testing.T) {
	method := func(id string) (any, error) {
		return nil, nil
	}

	type dependencies struct {
	}

	setup := testit.FuncStateless(method, func(t *testing.T, d *dependencies) {

	})

	t.Run("test", setup.Expect(
		func(t testit.StatelessTest[dependencies, func(string) (any, error)]) {
			resp, err := t.Run("1")
			t.NoError(err)
			t.Nil(resp)
		},
	))
}
