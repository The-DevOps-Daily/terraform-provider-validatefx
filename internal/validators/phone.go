package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var phoneE164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

var _ frameworkvalidator.String = Phone()

// Phone returns a schema.String validator that ensures values follow the E.164 phone number format.
func Phone() frameworkvalidator.String {
	return phoneValidator{}
}

type phoneValidator struct{}

// Description returns a plain-text description of the validator.
func (phoneValidator) Description(_ context.Context) string {
	return "value must be a valid E.164 phone number (start with '+' followed by 1–15 digits)"
}

// MarkdownDescription returns a markdown-formatted description of the validator.
func (phoneValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid E.164 phone number"
}

// ValidateString performs the actual phone number validation.
func (phoneValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	if !phoneE164Regex.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Phone Number",
			fmt.Sprintf("Value %q is not a valid E.164 phone number. It must start with '+' followed by 1–15 digits.", value),
		)
	}
}
