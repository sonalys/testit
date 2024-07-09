package testit

// NewFactory is a function to create a state factory.
// It allows you to create a valid initial state and modify it with functions.
// It is useful for creating test cases that verify field-by-field validation.
// It's also useful for centralizing the creation of a single valid struct,
// as validation rule changes, you only need to change it in one place.
func NewFactory[T any](defaultState T) func(...func(*T)) T {
	return func(f ...func(*T)) T {
		state := DeepClone(defaultState)
		for _, fn := range f {
			fn(&state)
		}
		return state
	}
}
