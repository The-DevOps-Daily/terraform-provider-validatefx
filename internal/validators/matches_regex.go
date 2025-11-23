package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MatchesRegex returns a validation that checks input against the provided regular expression pattern.
func MatchesRegex(pattern string) frameworkvalidator.String {
	return matchesRegexValidator{pattern: pattern}
}

type matchesRegexValidator struct {
	pattern string
}

func (v matchesRegexValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must match regex pattern %q", v.pattern)
}

func (v matchesRegexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v matchesRegexValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	compiled, err := regexp.Compile(v.pattern)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Regex Pattern",
			fmt.Sprintf("Pattern %q is not a valid regular expression: %s", v.pattern, err.Error()),
		)
		return
	}

	if !compiled.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Regex Mismatch",
			fmt.Sprintf("Value %q does not match regex pattern %q", value, v.pattern),
		)
	}
}
