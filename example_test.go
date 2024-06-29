package testit_test

import (
	"testing"

	"github.com/sonalys/testit"
)

func Test(t *testing.T) {
	type Dependencies struct{}
	type TestCase struct{}
	type TargetFunc = func(id string) (any, error)
	type Config = testit.Config[Dependencies, TestCase, TargetFunc]

	targetFn := func(id string) (any, error) {
		return nil, nil
	}

	setup := testit.New(func(t *testing.T, d *Dependencies, tc *TestCase) TargetFunc {
		println("pre-run setup")
		return targetFn
	})

	additionalStep := func(t Config) {
		println("additional step")
	}

	t.Run("test case", setup.Case(
		TestCase{},
		additionalStep,
		func(t Config) {
			got, err := t.Run("id")
			println("after-run assertions")
			t.NoError(err)
			t.Nil(got)
		},
	))
}
