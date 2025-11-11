package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var integerRe = regexp.MustCompile(`^[+-]?\d+$`)

type integerValidator struct{}

var _ frameworkvalidator.String = (*integerValidator)(nil)

// Integer returns a validator ensuring the string represents a valid integer.
func Integer() frameworkvalidator.String { return &integerValidator{} }

func (v *integerValidator) Description(_ context.Context) string {
	return "value must be a valid integer"
}

func (v *integerValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v *integerValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if value == "" {
		return
	}

	if !integerRe.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Integer",
			fmt.Sprintf("Value %q is not a valid integer", value),
		)
	}
}
