package ptr

// BoolValue gets the value of a *bool.
// Nil is treated as false.
func BoolValue(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// Bool returns a pointer to this bool value.
func Bool(v bool) *bool {
	return &v
}
