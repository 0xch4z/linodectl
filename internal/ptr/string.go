package ptr

// String returns a pointer to this bool value.
func String(s string) *string {
	return &s
}
