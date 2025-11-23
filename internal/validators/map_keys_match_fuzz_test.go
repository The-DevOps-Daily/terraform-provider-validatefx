//go:build gofuzz || go1.18

package validators

import (
	"testing"
)

func FuzzMapKeysMatch(f *testing.F) {
	// Seed with valid scenarios
	f.Add("a,b,c", "a,b", "a,b")
	f.Add("a,b", "", "a")
	f.Add("", "a,b", "a,b,c")

	// Seed with invalid scenarios
	f.Add("a,b", "", "a,c")
	f.Add("a,b", "a,b", "a")

	f.Fuzz(func(t *testing.T, allowed, required, input string) {
		allowedList := splitKeys(allowed)
		requiredList := splitKeys(required)
		inputList := splitKeys(input)

		v := NewMapKeysMatch(allowedList, requiredList)

		// Should not panic
		_ = v.Validate(inputList)
	})
}

func splitKeys(s string) []string {
	if s == "" {
		return []string{}
	}
	keys := []string{}
	for _, k := range s {
		keys = append(keys, string(k))
	}
	return keys
}
