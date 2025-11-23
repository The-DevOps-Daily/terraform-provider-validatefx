package validators

import (
	"context"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*slugValidator)(nil)

// Slug returns a validator that checks if a string is a valid slug.
// A valid slug consists of lowercase letters, digits, and hyphens.
// It must start and end with a letter or digit, and hyphens cannot be consecutive.
func Slug() frameworkvalidator.String {
	return &slugValidator{}
}

type slugValidator struct{}

func (v *slugValidator) Description(_ context.Context) string {
	return "Value must be a valid slug (lowercase letters, digits, and hyphens; no leading/trailing or consecutive hyphens)"
}

func (v *slugValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *slugValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if !slugRegex.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Slug",
			"Value must be a valid slug (lowercase letters, digits, and hyphens; no leading/trailing or consecutive hyphens).",
		)
	}
}
