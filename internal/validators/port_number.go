package validators

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*portNumberValidator)(nil)

// PortNumber returns a validator ensuring a string represents a valid TCP/UDP port number (1..65535).
func PortNumber() frameworkvalidator.String { return &portNumberValidator{} }

type portNumberValidator struct{}

func (portNumberValidator) Description(_ context.Context) string {
	return "string must be a valid port number (1..65535)"
}

func (v portNumberValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (portNumberValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := strings.TrimSpace(req.ConfigValue.ValueString())
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 65535 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Port Number",
			fmt.Sprintf("Value %q must be an integer between 1 and 65535.", req.ConfigValue.ValueString()),
		)
		return
	}
}
