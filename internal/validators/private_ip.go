package validators

import (
	"context"
	"fmt"
	"net"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*privateIPValidator)(nil)

// PrivateIP returns a validator ensuring a string is an RFC1918 private IP address (IPv4 or IPv6 unique local address).
func PrivateIP() frameworkvalidator.String {
	return &privateIPValidator{}
}

type privateIPValidator struct{}

func (privateIPValidator) Description(_ context.Context) string {
	return "string must be a private IP address (RFC1918 for IPv4 or unique local address for IPv6)"
}

func (v privateIPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (privateIPValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	ip := net.ParseIP(s)
	if ip == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid IP",
			fmt.Sprintf("Value %q is not a valid IP address.", s),
		)
		return
	}

	// RFC1918 private ranges and IPv6 Unique Local Addresses (fc00::/7).
	if isPrivateIP(ip) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Not a Private IP",
		fmt.Sprintf("Value %q is not within private IP ranges.", s),
	)
}

func isPrivateIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return isPrivateIPv4(ip4)
	}
	return isPrivateIPv6(ip)
}

func isPrivateIPv4(ip4 net.IP) bool {
	// 10.0.0.0/8
	if ip4[0] == 10 {
		return true
	}
	// 172.16.0.0/12 => 172.16.0.0 - 172.31.255.255
	if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ip4[0] == 192 && ip4[1] == 168 {
		return true
	}
	return false
}

func isPrivateIPv6(ip net.IP) bool {
	// IPv6 Unique Local Address fc00::/7 (fc00::/8 and fd00::/8)
	// Check first 7 bits are 1111110x; equivalently first byte 0xfc or 0xfd.
	if len(ip) == net.IPv6len {
		b0 := ip[0]
		if b0 == 0xfc || b0 == 0xfd {
			return true
		}
	}
	return false
}

// IsLinkLocalIP reports whether the IP is link-local (IPv4 169.254.0.0/16, IPv6 fe80::/10).
func IsLinkLocalIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 169 && ip4[1] == 254
	}
	// fe80::/10 => first 10 bits 1111111010; check first byte 0xfe and next nibble 0x8..0xb
	if len(ip) == net.IPv6len {
		return ip[0] == 0xfe && (ip[1]&0xc0) == 0x80
	}
	return false
}

// IsReservedIP reports whether the IP is from commonly reserved/non-routable ranges
// (loopback, documentation ranges, multicast, broadcast, CGNAT, etc.). This is not
// exhaustive but covers the most relevant ranges for public routing checks.
func IsReservedIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return isReservedIPv4(ip4)
	}
	return isReservedIPv6(ip)
}

var reservedIPv4Nets []*net.IPNet
var reservedIPv6Nets []*net.IPNet

func init() {
	// Common non-routable/assigned ranges (IPv4)
	v4 := []string{
		"0.0.0.0/8",          // this network
		"127.0.0.0/8",        // loopback
		"100.64.0.0/10",      // CGNAT
		"192.0.0.0/24",       // IETF Protocol Assignments
		"192.0.2.0/24",       // TEST-NET-1
		"198.51.100.0/24",    // TEST-NET-2
		"203.0.113.0/24",     // TEST-NET-3
		"224.0.0.0/4",        // multicast
		"240.0.0.0/4",        // reserved
		"255.255.255.255/32", // broadcast
	}
	for _, c := range v4 {
		if _, n, err := net.ParseCIDR(c); err == nil {
			reservedIPv4Nets = append(reservedIPv4Nets, n)
		}
	}

	// Common non-routable/assigned ranges (IPv6)
	v6 := []string{
		"::/128",        // unspecified
		"::1/128",       // loopback
		"ff00::/8",      // multicast
		"2001:db8::/32", // documentation
	}
	for _, c := range v6 {
		if _, n, err := net.ParseCIDR(c); err == nil {
			reservedIPv6Nets = append(reservedIPv6Nets, n)
		}
	}
}

func isReservedIPv4(ip4 net.IP) bool {
	ip := net.IP(ip4)
	for _, n := range reservedIPv4Nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func isReservedIPv6(ip net.IP) bool {
	if len(ip) != net.IPv6len {
		return false
	}
	for _, n := range reservedIPv6Nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
