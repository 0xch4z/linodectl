package strutil

import "strings"

func Mask(s string, ch rune) string {
	return strings.Repeat(string(ch), len(s))
}
