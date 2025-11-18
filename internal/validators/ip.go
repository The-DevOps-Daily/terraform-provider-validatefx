package validators

import (
	"context"
	"fmt"
	"net"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = IP()

// IP returns a schema.String validator that ensures the value is a valid IP address.
func IP() frameworkvalidator.String {
	return ipValidator{}
}

type ipValidator struct{}

func (ipValidator) Description(_ context.Context) string {
	return "value must be a valid IP address"
}

func (ipValidator) MarkdownDescription(_ context.Context) string {
	return ipValidator{}.Description(nil)
}

func (ipValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	if ip := net.ParseIP(value); ip == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid IP Address",
			fmt.Sprintf("Value %q is not a valid IPv4 or IPv6 address", value),
		)
	}
}
