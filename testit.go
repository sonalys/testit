package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	TestCase interface {
		run(*testing.T)
	}

	TestHandler[Dependencies, Case, Result any] interface {
		Run(name string, runner TestCase)
		Case(tc Case, setup ...func(t *testing.T, d *Dependencies, tc Case, run func() Result)) TestCase
	}

	testHandler[Dependencies, Case, Result any] struct {
		t     *testing.T
		setup func(t *testing.T, d *Dependencies, tc *Case) func() Result
	}

	testCase[Dependencies, Case, Result any] struct {
		th    *testHandler[Dependencies, Case, Result]
		tc    Case
		setup []func(t *testing.T, d *Dependencies, tc *Case, run func() Result)
	}
)

func (th *testHandler[D, C, R]) Run(name string, runner TestCase) {
	require.NotPanics(th.t, func() {
		th.t.Run(name, runner.run)
	})
}

func (th *testHandler[D, C, R]) Case(tc C, setup ...func(t *testing.T, d *D, tc *C, run func() R)) TestCase {
	return &testCase[D, C, R]{
		th:    th,
		tc:    tc,
		setup: setup,
	}
}

func (ut *testCase[D, C, R]) run(t *testing.T) {
	var dependencies D
	// Ensure that nil interfaces aren't passed through, we need only structs to be initialized.
	if any(dependencies) != any(nil) {
		initializeMocks(t, &dependencies)
	}
	// Setup test handler base logic, retrieving the run func.
	run := ut.th.setup(t, &dependencies, &ut.tc)
	for _, f := range ut.setup {
		// Executes each test case func, parametrizing the run func from test handler.
		f(t, &dependencies, &ut.tc, run)
	}
}

// New initializes a new test handler.
func New[Dependencies, Case, Result any](t *testing.T, setup func(t *testing.T, d *Dependencies, tc *Case) func() (r Result)) *testHandler[Dependencies, Case, Result] {
	return &testHandler[Dependencies, Case, Result]{
		t:     t,
		setup: setup,
	}
}
