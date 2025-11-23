package validators

import (
	"context"
	"fmt"
	"regexp"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = MIMEType()

// MIME type pattern: type/subtype with optional parameters
// Examples: application/json, text/html; charset=utf-8, image/svg+xml
var mimeTypeRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9+.-]*/[a-zA-Z0-9][a-zA-Z0-9+.-]*(;.+)?$`)

// MIMEType returns a validator that verifies a string is a valid MIME type.
func MIMEType() frameworkvalidator.String {
	return mimeTypeValidator{}
}

type mimeTypeValidator struct{}

func (mimeTypeValidator) Description(_ context.Context) string {
	return "value must be a valid MIME type (e.g. application/json, text/html)"
}

func (v mimeTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (mimeTypeValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if !mimeTypeRe.MatchString(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MIME Type",
			fmt.Sprintf("Value %q is not a valid MIME type. Expected format: type/subtype (e.g. application/json).", value),
		)
	}
}
