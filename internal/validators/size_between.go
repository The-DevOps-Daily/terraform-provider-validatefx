package validators

import (
	"context"
	"fmt"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*sizeBetweenValidator)(nil)

// SizeBetween returns a validator that checks if a string represents a numeric value within an inclusive range.
func SizeBetween(min, max string) frameworkvalidator.String {
	return &sizeBetweenValidator{
		min: min,
		max: max,
	}
}

type sizeBetweenValidator struct {
	min string
	max string
}

func (v *sizeBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Value must be a number between %s and %s (inclusive)", v.min, v.max)
}

func (v *sizeBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *sizeBetweenValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Reuse the existing Between validation logic
	valid, boundsDiag, valueDiag := EvaluateBetween(value, v.min, v.max)

	if boundsDiag != nil {
		resp.Diagnostics.AddAttributeError(req.Path, boundsDiag.Summary, boundsDiag.Detail)
		return
	}

	if valueDiag != nil {
		resp.Diagnostics.AddAttributeError(req.Path, valueDiag.Summary, valueDiag.Detail)
		return
	}

	if !valid {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value Out of Range",
			fmt.Sprintf("Value must be between %s and %s (inclusive).", v.min, v.max),
		)
	}
}
