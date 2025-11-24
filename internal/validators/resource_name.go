package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*resourceNameValidator)(nil)

// Terraform resource names must start with a letter or underscore,
// followed by letters, digits, underscores, and hyphens.
var resourceNameRegexp = regexp.MustCompile(`^[a-z_][a-z0-9_-]*$`)

// ResourceName returns a validator ensuring a string matches Terraform resource naming conventions.
// Valid names must:
// - Start with a lowercase letter or underscore
// - Contain only lowercase letters, digits, underscores, and hyphens
// - Not be empty
func ResourceName() frameworkvalidator.String {
	return resourceNameValidator{}
}

type resourceNameValidator struct{}

func (resourceNameValidator) Description(_ context.Context) string {
	return "string must be a valid Terraform resource name (lowercase letters, digits, underscores, and hyphens; must start with letter or underscore)"
}

func (v resourceNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (resourceNameValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Resource Name",
			"Resource name cannot be empty.",
		)
		return
	}

	if !resourceNameRegexp.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Resource Name",
			fmt.Sprintf("Value %q must be a valid Terraform resource name. Names must start with a lowercase letter or underscore and contain only lowercase letters, digits, underscores, and hyphens.", value),
		)
	}
}
