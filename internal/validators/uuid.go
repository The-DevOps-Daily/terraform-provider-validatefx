package validators

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = UUID()

// UUID returns a schema.String validator which enforces RFC 4122 compliant UUIDs (versions 1-5).
func UUID() frameworkvalidator.String {
	return uuidValidator{}
}

type uuidValidator struct{}

func (uuidValidator) Description(_ context.Context) string {
	return "value must be a valid UUID (versions 1-5)"
}

func (uuidValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid UUID (versions 1-5)"
}

func (uuidValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid UUID",
			fmt.Sprintf("Value %q is not a valid UUID: %s", value, err.Error()),
		)
		return
	}

	version := parsed.Version()
	if version < 1 || version > 5 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unsupported UUID Version",
			fmt.Sprintf("Value %q is a UUID but version %d is not supported (expected v1-v5)", value, version),
		)
	}
}
