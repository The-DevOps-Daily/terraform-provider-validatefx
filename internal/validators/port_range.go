package validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*portRangeValidator)(nil)

var portRangeRegexp = regexp.MustCompile(`^\s*(\d{1,5})\s*-\s*(\d{1,5})\s*$`)

// PortRange returns a validator ensuring a string matches a valid TCP/UDP port range
// in the form "start-end" where 0 <= start <= end <= 65535.
func PortRange() frameworkvalidator.String {
	return &portRangeValidator{}
}

type portRangeValidator struct{}

func (portRangeValidator) Description(_ context.Context) string {
	return "string must be a valid port range (start-end) with values 0..65535 and start <= end"
}

func (v portRangeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (portRangeValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	raw := strings.TrimSpace(value)
	m := portRangeRegexp.FindStringSubmatch(raw)
	if len(m) != 3 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Port Range",
			fmt.Sprintf("Value %q must match the form 'start-end' using ports 0..65535.", value),
		)
		return
	}

	start, _ := strconv.Atoi(m[1])
	end, _ := strconv.Atoi(m[2])

	if start < 0 || start > 65535 || end < 0 || end > 65535 || start > end {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Port Range",
			fmt.Sprintf("Invalid port range %q: start and end must be within 0..65535 and start <= end.", value),
		)
		return
	}
}
