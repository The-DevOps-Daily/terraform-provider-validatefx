package validators

import (
	"context"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*positiveNumberValidator)(nil)

// PositiveNumber returns a validator that checks if a string represents a positive number (greater than zero).
func PositiveNumber() frameworkvalidator.String {
	return &positiveNumberValidator{}
}

type positiveNumberValidator struct{}

func (positiveNumberValidator) Description(_ context.Context) string {
	return "Value must be a positive number (greater than zero)"
}

func (v positiveNumberValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (positiveNumberValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	num, diag := parseFloat64(value, req.Path)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	if num <= 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Not a Positive Number",
			"Value must be greater than zero.",
		)
	}
}
