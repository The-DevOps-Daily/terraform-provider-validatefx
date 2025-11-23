//go:build gofuzz || go1.18

package validators

import (
	"strings"
	"testing"
)

func FuzzMutuallyExclusive(f *testing.F) {
	// Seed with valid scenarios (exactly one value set)
	f.Add("a,,")
	f.Add(",b,")
	f.Add(",,c")

	// Seed with invalid scenarios
	f.Add(",,")
	f.Add("a,b,")
	f.Add("a,b,c")

	f.Fuzz(func(t *testing.T, input string) {
		values := strings.Split(input, ",")

		v := NewMutuallyExclusive()

		// Should not panic
		_ = v.Validate(values)
	})
}
