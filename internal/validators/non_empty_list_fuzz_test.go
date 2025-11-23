//go:build gofuzz || go1.18

package validators

import (
	"strings"
	"testing"
)

func FuzzNonEmptyList(f *testing.F) {
	// Seed with valid scenarios (non-empty lists)
	f.Add("a")
	f.Add("a,b,c")
	f.Add("one")

	// Seed with invalid scenarios (empty list)
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		var values []string
		if input != "" {
			values = strings.Split(input, ",")
		} else {
			values = []string{}
		}

		v := NewNonEmptyList()

		// Should not panic
		_ = v.Validate(values)
	})
}
