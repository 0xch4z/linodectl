package ptr

// To returns a pointer to the given value.
func To[T any](v T) *T {
	return &v
}

// Value dereferences a pointer, returning the zero value if nil.
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
