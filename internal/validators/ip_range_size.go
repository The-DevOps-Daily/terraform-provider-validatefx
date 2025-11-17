package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IPRangeSizeValidator ensures a CIDR prefix length falls within [Min, Max].
type IPRangeSizeValidator struct {
	Min int
	Max int
}

var _ validator.String = (*IPRangeSizeValidator)(nil)

// NewIPRangeSizeValidator constructs a new validator with inclusive bounds.
func NewIPRangeSizeValidator(min, max int) IPRangeSizeValidator {
	// If caller passes inverted range, normalize to avoid surprising errors.
	if min > max {
		min, max = max, min
	}
	return IPRangeSizeValidator{Min: min, Max: max}
}

func (v IPRangeSizeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("CIDR prefix length must be between /%d and /%d", v.Min, v.Max)
}

func (v IPRangeSizeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v IPRangeSizeValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	if s == "" {
		return
	}

	_, ipNet, err := net.ParseCIDR(s)
	if err != nil || ipNet == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			fmt.Sprintf("Value %q is not a valid CIDR: %v", s, err),
		)
		return
	}

	ones, bits := ipNet.Mask.Size()
	if ones < 0 || bits <= 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR Mask",
			fmt.Sprintf("Value %q has an invalid mask", s),
		)
		return
	}

	if ones < v.Min || ones > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Mask Out Of Range",
			fmt.Sprintf("CIDR %q has prefix /%d which is outside allowed range /%d to /%d.", s, ones, v.Min, v.Max),
		)
		return
	}
}
