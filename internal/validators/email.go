package validators

import (
	"context"
	"fmt"
	"net/mail"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = Email()

// Email returns a schema.String validator which enforces RFC 5322 compliant emails.
func Email() frameworkvalidator.String {
	return emailValidator{}
}

type emailValidator struct{}

func (emailValidator) Description(_ context.Context) string {
	return "value must be a valid email address"
}

func (emailValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid RFC 5322 email address"
}

func (emailValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if _, err := mail.ParseAddress(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Email Address",
			fmt.Sprintf("Value %q is not a valid email address: %s", value, err.Error()),
		)
	}
}
