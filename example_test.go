package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test_Stateful(t *testing.T) {
	type Dependencies struct{}
	type TestCase struct{}
	type TargetFunc = func(id string) (any, error)
	type Test = testit.StatefulTest[Dependencies, TestCase, TargetFunc]

	targetFn := func(id string) (any, error) {
		return nil, nil
	}

	setup := testit.Stateful(func(t *testing.T, d *Dependencies, tc *TestCase) TargetFunc {
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

func Test_Stateless(t *testing.T) {
	type Dependencies struct{}
	type TargetFunc = func(id string) (any, error)

	setup := testit.Stateless(func(t *testing.T, d *Dependencies) TargetFunc {
		println("pre-run setup")
		return func(id string) (any, error) {
			return nil, nil
		}
	})

	setup = setup.PreCase(func(t *testing.T, d *Dependencies) {
		println("pre-case setup")
	})

	setup = setup.PostCase(func(t *testing.T, d *Dependencies) {
		println("post-case setup")
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
