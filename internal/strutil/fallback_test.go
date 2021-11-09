package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFallback(t *testing.T) {
	t.Run("should return fallback when string is empty", func(t *testing.T) {
		v := Fallback("", "foo")
		assert.Equal(t, v, "foo")
	})

	t.Run("should return first string when not empty", func(t *testing.T) {
		v := Fallback("shouldequal", "shouldnot")
		assert.Equal(t, v, "shouldequal")
	})
}

func TestSliceFallback(t *testing.T) {
	t.Run("should return fallback when slice is empty", func(t *testing.T) {
		v := SliceFallback([]string{}, []string{"hi"})
		assert.Equal(t, v, []string{"hi"})
	})

	t.Run("should return first string when not empty", func(t *testing.T) {
		v := SliceFallback([]string{"me"}, []string{"notme"})
		assert.Equal(t, v, []string{"me"})
	})
}
