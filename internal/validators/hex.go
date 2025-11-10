package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var hexRe = regexp.MustCompile(`^[0-9a-fA-F]+$`)

type hexValidator struct{}

var _ frameworkvalidator.String = (*hexValidator)(nil)

// Hex returns a validator ensuring the string contains only hexadecimal characters (0-9, a-f, A-F).
func Hex() frameworkvalidator.String { return &hexValidator{} }

func (v *hexValidator) Description(_ context.Context) string {
	return "string must contain only hexadecimal characters (0-9, a-f, A-F)"
}

func (v *hexValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v *hexValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	if s == "" || !hexRe.MatchString(s) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Hex String",
			fmt.Sprintf("Value %q must contain only hexadecimal characters (0-9, a-f, A-F).", s),
		)
	}
}
