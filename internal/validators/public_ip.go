package validators

import (
	"context"
	"fmt"
	"net"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*publicIPValidator)(nil)

// PublicIP returns a validator ensuring a string is a public (non-private) IP address.
func PublicIP() frameworkvalidator.String { return &publicIPValidator{} }

type publicIPValidator struct{}

func (v *publicIPValidator) Description(_ context.Context) string {
	return "string must be a public IP address (not RFC1918/ULA)"
}

func (v *publicIPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *publicIPValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	ip := net.ParseIP(s)
	if ip == nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid IP", fmt.Sprintf("Value %q is not a valid IP address.", s))
		return
	}

	if isPrivateIP(ip) {
		resp.Diagnostics.AddAttributeError(req.Path, "Not a Public IP", fmt.Sprintf("Value %q is a private IP address.", s))
		return
	}
}
