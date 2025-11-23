package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var semverPattern = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)

var _ frameworkvalidator.String = SemVer()

// SemVer returns a schema.String validator ensuring semantic versioning per semver.org.
func SemVer() frameworkvalidator.String {
	return semverValidator{}
}

type semverValidator struct{}

func (semverValidator) Description(_ context.Context) string {
	return "value must be a Semantic Version (SemVer 2.0.0)"
}

func (v semverValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (semverValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	if !semverPattern.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Semantic Version",
			fmt.Sprintf("Value %q is not a valid semantic version", value),
		)
	}
}
