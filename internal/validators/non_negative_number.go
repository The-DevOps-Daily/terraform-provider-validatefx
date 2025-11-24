package validators

import (
	"context"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*nonNegativeNumberValidator)(nil)

// NonNegativeNumber returns a validator that checks if a string represents a non-negative number (zero or greater).
func NonNegativeNumber() frameworkvalidator.String {
	return &nonNegativeNumberValidator{}
}

type nonNegativeNumberValidator struct{}

func (nonNegativeNumberValidator) Description(_ context.Context) string {
	return "Value must be a non-negative number (zero or greater)"
}

func (v nonNegativeNumberValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (nonNegativeNumberValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	num, diag := parseFloat64(value, req.Path)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	if num < 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Negative Number",
			"Value must be zero or greater.",
		)
	}
}
