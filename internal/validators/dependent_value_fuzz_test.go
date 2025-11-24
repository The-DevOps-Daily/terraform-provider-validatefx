//go:build gofuzz
// +build gofuzz

package validators

import (
	"testing"
)

func FuzzDependentValue(f *testing.F) {
	// Valid seeds
	f.Add("", "")           // both empty
	f.Add("a", "b")         // both set
	f.Add("", "value")      // condition empty, dependent set
	f.Add("enabled", "cfg") // both set

	// Invalid seeds
	f.Add("value", "")  // condition set, dependent empty
	f.Add("true", "")   // condition true, dependent empty

	f.Fuzz(func(t *testing.T, condition, dependent string) {
		validator := NewDependentValue()
		_ = validator.Validate(condition, dependent)
		// Just ensure it doesn't panic
	})
}
