package validators

import (
	"context"
	"fmt"
	"net"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*subnetValidator)(nil)

// Subnet validates IPv4/IPv6 CIDR blocks where the IP equals the network address (subnet address).
func Subnet() frameworkvalidator.String { return &subnetValidator{} }

type subnetValidator struct{}

func (subnetValidator) Description(_ context.Context) string {
	return "value must be a subnet address in CIDR notation (IP equals network address)"
}

func (subnetValidator) MarkdownDescription(ctx context.Context) string {
	return subnetValidator{}.Description(ctx)
}

func (subnetValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid CIDR", fmt.Sprintf("Value %q is not a valid CIDR: %v", s, err))
		return
	}
	// Check IP equals network address
	if !ip.Equal(ipNet.IP) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Subnet Address", fmt.Sprintf("Value %q IP must equal network address %s", s, ipNet.IP.String()))
		return
	}
}
