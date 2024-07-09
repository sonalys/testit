package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Mutation(t *testing.T) {
	stateless := Setup(func(t *Hook[struct{}]) {})
	require.Len(t, stateless.hooks, 0)

	var flag bool

	newStateless := stateless.Hook(func(t *Hook[struct{}]) {
		flag = true
	})

	require.Len(t, stateless.hooks, 0)
	require.Len(t, newStateless.hooks, 1)

	stateless = stateless.Hook(func(t *Hook[struct{}]) {})

	require.Len(t, stateless.hooks, 1)

	stateless.Expect()(t)
	require.False(t, flag)

	newStateless.Expect()(t)
	require.True(t, flag)
}
