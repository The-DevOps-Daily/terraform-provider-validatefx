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

func (v *privateIPValidator) Description(_ context.Context) string {
	return "string must be a private IP address (RFC1918 for IPv4 or unique local address for IPv6)"
}

func (v *privateIPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *privateIPValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
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
