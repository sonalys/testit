package testit

func NoErr[D, C, R, V any](t *Test[D, C, R], value V, err error) V {
	t.NoError(err)
	return value
}
