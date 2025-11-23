package validators

import (
	"context"
	"encoding/json"
	"fmt"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = JSON()

// JSON returns a schema.String validator that ensures the value encodes a JSON object.
func JSON() frameworkvalidator.String {
	return jsonValidator{}
}

type jsonValidator struct{}

func (jsonValidator) Description(_ context.Context) string {
	return "value must be a valid JSON object"
}

func (jsonValidator) MarkdownDescription(_ context.Context) string {
	return "value must be valid JSON"
}

func (jsonValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON",
			fmt.Sprintf("Value %q is not a valid JSON object: %s", value, err.Error()),
		)
		return
	}

	if _, ok := decoded.(map[string]any); !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON Object",
			fmt.Sprintf("Value %q must decode to a JSON object", value),
		)
	}
}
