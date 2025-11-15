package validators

import (
	"fmt"
	"net"
)

// CIDROverlapValidator validates that a set of CIDR blocks do not overlap.
// It returns an error if any two CIDR ranges intersect or if any entry is invalid.
type CIDROverlapValidator struct{}

// NewCIDROverlap returns a new instance of the CIDR overlap validator.
func NewCIDROverlap() *CIDROverlapValidator { return &CIDROverlapValidator{} }

// Validate returns nil if no overlaps are found among the provided CIDR blocks.
func (v *CIDROverlapValidator) Validate(cidrs []string) error {
	if v == nil {
		return fmt.Errorf("validator not initialized")
	}

	// Parse and normalize
	type entry struct {
		ipnet *net.IPNet
		base  net.IP
		raw   string
		ipv6  bool
	}

	parsed := make([]entry, 0, len(cidrs))
	for _, raw := range cidrs {
		if raw == "" { // skip empties; treat as invalid input
			return fmt.Errorf("invalid CIDR: empty string")
		}
		ip, n, err := net.ParseCIDR(raw)
		if err != nil {
			return fmt.Errorf("invalid CIDR %q: %w", raw, err)
		}
		base := ip.Mask(n.Mask)
		parsed = append(parsed, entry{ipnet: n, base: base, raw: raw, ipv6: ip.To4() == nil})
	}

	// Compare for overlap; networks overlap if either base IP is contained in the other.
	for i := 0; i < len(parsed); i++ {
		for j := i + 1; j < len(parsed); j++ {
			a, b := parsed[i], parsed[j]
			// address family mismatch cannot overlap
			if a.ipv6 != b.ipv6 {
				continue
			}
			if a.ipnet.Contains(b.base) || b.ipnet.Contains(a.base) {
				return fmt.Errorf("CIDR overlap detected between %q and %q", a.raw, b.raw)
			}
		}
	}

	return nil
}
