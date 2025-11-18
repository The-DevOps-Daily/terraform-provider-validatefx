package validators

import (
	"context"
	"fmt"
	"net"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = CIDR()

// CIDR validates IPv4 and IPv6 CIDR blocks.
func CIDR() frameworkvalidator.String {
	return cidrValidator{}
}

type cidrValidator struct{}

func (cidrValidator) Description(_ context.Context) string {
	return "value must be a valid IPv4 or IPv6 CIDR block"
}

func (cidrValidator) MarkdownDescription(_ context.Context) string {
	return cidrValidator{}.Description(nil)
}

func (cidrValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	_, ipNet, err := net.ParseCIDR(value)
	if err == nil {
		if ones, bits := ipNet.Mask.Size(); ones < 0 || bits < 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid CIDR Mask",
				fmt.Sprintf("Value %q has an invalid mask", value),
			)
		}
		return
	}

	summary := "Invalid CIDR"
	detail := fmt.Sprintf("Value %q is not a valid CIDR block: %s", value, err.Error())

	if strings.Contains(value, "/") {
		parts := strings.SplitN(value, "/", 2)
		if len(parts) == 2 && net.ParseIP(parts[0]) != nil {
			summary = "Invalid CIDR Mask"
			detail = fmt.Sprintf("Value %q has an invalid mask: %s", value, err.Error())
		}
	}

	resp.Diagnostics.AddAttributeError(req.Path, summary, detail)
}
