package testit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Mutation(t *testing.T) {
	stateless := Stateless(func(t *testing.T, d *struct{}) struct{} {
		return struct{}{}
	})
	require.Len(t, stateless.pre, 0)
	require.Len(t, stateless.post, 0)

	var flag bool

	newStateless := stateless.PreCase(func(t *testing.T, d *struct{}) {
		flag = true
	})
	newStateless = newStateless.PostCase(func(t *testing.T, d *struct{}) {})

	require.Len(t, stateless.pre, 0)
	require.Len(t, stateless.post, 0)
	require.Len(t, newStateless.pre, 1)
	require.Len(t, newStateless.post, 1)

	stateless = stateless.PreCase(func(t *testing.T, d *struct{}) {})

	require.Len(t, stateless.pre, 1)
	require.Len(t, stateless.post, 0)

	stateless.Expect()(t)
	require.False(t, flag)

	newStateless.Expect()(t)
	require.True(t, flag)
}
