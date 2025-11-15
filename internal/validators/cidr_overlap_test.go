package validators

import "testing"

func TestCIDROverlapValidator_NoOverlap(t *testing.T) {
	v := NewCIDROverlap()
	cases := [][]string{
		{"10.0.0.0/24", "10.0.1.0/24"},
		{"192.168.0.0/25", "192.168.0.128/25"},
		{"2001:db8::/32", "2001:db9::/32"},
	}
	for _, c := range cases {
		if err := v.Validate(c); err != nil {
			t.Fatalf("expected no overlap for %v, got %v", c, err)
		}
	}
}

func TestCIDROverlapValidator_Overlap(t *testing.T) {
	v := NewCIDROverlap()
	cases := [][]string{
		{"10.0.0.0/24", "10.0.0.128/25"},
		{"192.168.0.0/24", "192.168.0.0/25"},
		{"2001:db8::/32", "2001:db8:1::/48"},
	}
	for _, c := range cases {
		if err := v.Validate(c); err == nil {
			t.Fatalf("expected overlap for %v", c)
		}
	}
}

func TestCIDROverlapValidator_Invalid(t *testing.T) {
	v := NewCIDROverlap()
	if err := v.Validate([]string{"not-a-cidr", "10.0.0.0/24"}); err == nil {
		t.Fatalf("expected error for invalid CIDR input")
	}
}
