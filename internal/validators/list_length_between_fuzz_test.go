//go:build gofuzz
// +build gofuzz

package validators

import (
	"testing"
)

func FuzzListLengthBetween(f *testing.F) {
	// Valid seeds
	f.Add(2, 5, 3) // min, max, length
	f.Add(0, 10, 5)
	f.Add(1, 1, 1)
	f.Add(0, 5, 0)

	// Invalid seeds
	f.Add(2, 5, 1)  // too short
	f.Add(2, 5, 10) // too long
	f.Add(5, 10, 3) // below min

	f.Fuzz(func(t *testing.T, min, max, length int) {
		// Skip invalid test inputs
		if min < 0 || max < 0 || min > max || length < 0 {
			return
		}

		// Create a list with the specified length
		list := make([]string, length)
		for i := 0; i < length; i++ {
			list[i] = "item"
		}

		validator := NewListLengthBetween(min, max)
		_ = validator.Validate(list)
		// Just ensure it doesn't panic
	})
}
