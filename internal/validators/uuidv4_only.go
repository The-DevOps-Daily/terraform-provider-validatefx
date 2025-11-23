package validators

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = UUIDv4Only()

// UUIDv4Only returns a schema.String validator which enforces UUID version 4 only.
func UUIDv4Only() frameworkvalidator.String {
	return uuidv4OnlyValidator{}
}

type uuidv4OnlyValidator struct{}

func (uuidv4OnlyValidator) Description(_ context.Context) string {
	return "value must be a valid UUID version 4"
}

func (v uuidv4OnlyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (uuidv4OnlyValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
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
	if version != 4 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid UUID Version",
			fmt.Sprintf("Value %q is a valid UUID but version %d is not version 4", value, version),
		)
	}
}
