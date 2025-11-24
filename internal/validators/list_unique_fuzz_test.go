//go:build gofuzz
// +build gofuzz

package validators

import (
	"testing"
)

func FuzzListUnique(f *testing.F) {
	// Valid seeds - unique lists
	f.Add("a,b,c")
	f.Add("1,2,3,4,5")
	f.Add("apple,banana,cherry")
	f.Add("")

	// Invalid seeds - with duplicates
	f.Add("a,b,a")
	f.Add("1,2,1")
	f.Add("x,x,x")

	f.Fuzz(func(t *testing.T, input string) {
		if input == "" {
			// Empty list
			validator := NewListUnique()
			_ = validator.Validate([]string{})
			return
		}

		// Split by comma to create list
		var items []string
		for i, c := range input {
			if c == ',' || i == len(input)-1 {
				continue
			}
			items = append(items, string(c))
		}

		validator := NewListUnique()
		_ = validator.Validate(items)
		// Just ensure it doesn't panic
	})
}
