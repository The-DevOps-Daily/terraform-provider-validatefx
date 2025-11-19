package validators

import (
	"context"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPrivateIPValidator(t *testing.T) {
	t.Parallel()

	v := PrivateIP()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	ok := []string{"10.0.0.1", "172.16.0.1", "172.31.255.254", "192.168.1.1", "fd00::1", "fc00:abcd::1"}
	bad := []string{"8.8.8.8", "172.15.0.1", "172.32.0.1", "1.2.3.4", "2001:db8::1", "::1", "not-an-ip"}

	for _, s := range ok {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected %s to be private, got %v", s, resp.Diagnostics)
		}
	}
	for _, s := range bad {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected %s to be rejected", s)
		}
	}

	// Null/unknown pass-through
	req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringNull()}
	resp := &frameworkvalidator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics for null")
	}

	req = frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringUnknown()}
	resp = &frameworkvalidator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics for unknown")
	}
}

// TestIsLinkLocalIP tests the IsLinkLocalIP helper function.
func TestIsLinkLocalIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// IPv4 link-local
		{"IPv4 link-local start", "169.254.0.0", true},
		{"IPv4 link-local middle", "169.254.100.50", true},
		{"IPv4 link-local end", "169.254.255.255", true},
		{"IPv4 not link-local", "169.253.0.0", false},
		{"IPv4 not link-local 2", "169.255.0.0", false},
		{"IPv4 public", "8.8.8.8", false},
		{"IPv4 private", "192.168.1.1", false},
		// IPv6 link-local fe80::/10
		{"IPv6 link-local", "fe80::1", true},
		{"IPv6 link-local full", "fe80:0000:0000:0000:0000:0000:0000:0001", true},
		{"IPv6 link-local upper", "febf:ffff:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"IPv6 not link-local", "fe00::1", false},
		{"IPv6 not link-local 2", "fec0::1", false},
		{"IPv6 public", "2001:4860:4860::8888", false},
		{"IPv6 ULA", "fd00::1", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP %q", tc.ip)
			}
			result := IsLinkLocalIP(ip)
			if result != tc.expected {
				t.Errorf("IsLinkLocalIP(%q) = %v, want %v", tc.ip, result, tc.expected)
			}
		})
	}
}

// TestIsReservedIP tests the IsReservedIP helper function for IPv4 and IPv6.
func TestIsReservedIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// IPv4 reserved ranges
		{"IPv4 loopback", "127.0.0.1", true},
		{"IPv4 this network", "0.0.0.0", true},
		{"IPv4 CGNAT", "100.64.0.1", true},
		{"IPv4 TEST-NET-1", "192.0.2.1", true},
		{"IPv4 TEST-NET-2", "198.51.100.1", true},
		{"IPv4 TEST-NET-3", "203.0.113.1", true},
		{"IPv4 multicast", "224.0.0.1", true},
		{"IPv4 reserved", "240.0.0.1", true},
		{"IPv4 broadcast", "255.255.255.255", true},
		{"IPv4 public", "8.8.8.8", false},
		{"IPv4 private", "192.168.1.1", false},
		// IPv6 reserved ranges
		{"IPv6 unspecified", "::", true},
		{"IPv6 loopback", "::1", true},
		{"IPv6 multicast", "ff02::1", true},
		{"IPv6 documentation", "2001:db8::1", true},
		{"IPv6 public", "2001:4860:4860::8888", false},
		{"IPv6 ULA", "fd00::1", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP %q", tc.ip)
			}
			result := IsReservedIP(ip)
			if result != tc.expected {
				t.Errorf("IsReservedIP(%q) = %v, want %v", tc.ip, result, tc.expected)
			}
		})
	}
}

// TestIsReservedIPv4EdgeCases tests edge cases for IPv4 reserved checking.
func TestIsReservedIPv4EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"0.0.0.0/8 start", "0.0.0.0", true},
		{"0.0.0.0/8 end", "0.255.255.255", true},
		{"127.0.0.0/8 middle", "127.100.50.25", true},
		{"100.64.0.0/10 start", "100.64.0.0", true},
		{"100.64.0.0/10 end", "100.127.255.255", true},
		{"192.0.0.0/24 range", "192.0.0.100", true},
		{"224.0.0.0/4 multicast start", "224.0.0.0", true},
		{"224.0.0.0/4 multicast end", "239.255.255.255", true},
		{"240.0.0.0/4 reserved", "240.0.0.1", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP %q", tc.ip)
			}
			result := IsReservedIP(ip)
			if result != tc.expected {
				t.Errorf("IsReservedIP(%q) = %v, want %v", tc.ip, result, tc.expected)
			}
		})
	}
}

// TestIsReservedIPv6EdgeCases tests edge cases for IPv6 reserved checking.
func TestIsReservedIPv6EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"IPv6 multicast start", "ff00::", true},
		{"IPv6 multicast end", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"IPv6 doc range start", "2001:db8::", true},
		{"IPv6 doc range end", "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"IPv6 not in reserved", "2001:db7:ffff:ffff:ffff:ffff:ffff:ffff", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP %q", tc.ip)
			}
			result := IsReservedIP(ip)
			if result != tc.expected {
				t.Errorf("IsReservedIP(%q) = %v, want %v", tc.ip, result, tc.expected)
			}
		})
	}
}
