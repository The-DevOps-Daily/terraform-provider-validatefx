package validators

import (
	"testing"
)

// FuzzCIDROverlapValidator fuzzes the CIDR overlap validator with pairs of strings.
// It ensures robustness (no panics) across arbitrary inputs, and seeds include
// both overlapping and non-overlapping examples.
func FuzzCIDROverlapValidator(f *testing.F) {
	f.Add("10.0.0.0/24", "10.0.1.0/24")     // non-overlap
	f.Add("10.0.0.0/24", "10.0.0.128/25")   // overlap
	f.Add("2001:db8::/32", "2001:db9::/32") // non-overlap IPv6
	f.Add("", "10.0.0.0/24")                // invalid
	f.Add("not-a-cidr", "10.0.0.0/24")      // invalid

	v := NewCIDROverlap()
	f.Fuzz(func(t *testing.T, a, b string) {
		t.Parallel()
		// Robustness check; validator returns error on invalid/overlap
		_ = v.Validate([]string{a, b})
	})
}
