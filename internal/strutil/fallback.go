package strutil

// Fallback returns x unless it's and empty string, then y would be returned.
func Fallback(x, y string) string {
	if x == "" {
		return y
	}
	return x
}

// Fallback returns x unless it's nil or empty, then y would be returned.
func SliceFallback(x, y []string) []string {
	if x == nil || len(x) == 0 {
		return y
	}
	return x
}
