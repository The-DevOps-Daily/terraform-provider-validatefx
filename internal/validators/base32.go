package validators

import (
	"context"
	"encoding/base32"
	"fmt"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type base32Validator struct{}

var _ frameworkvalidator.String = (*base32Validator)(nil)

// Base32Validator returns a validator that verifies a string is Base32-encoded.
func Base32Validator() frameworkvalidator.String { return base32Validator{} }

func (base32Validator) Description(_ context.Context) string { return "value must be a Base32 string" }
func (base32Validator) MarkdownDescription(_ context.Context) string {
	return "value must be a Base32 string"
}

func (base32Validator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() == "" {
		return
	}

	if _, err := base32.StdEncoding.DecodeString(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Base32 string",
			fmt.Sprintf("Value %q is not a valid Base32 string", req.ConfigValue.ValueString()),
		)
	}
}
