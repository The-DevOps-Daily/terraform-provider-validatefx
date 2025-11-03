package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type inListValidator struct {
	allowed    []string
	ignoreCase bool
	lookup     map[string]struct{}
}

var _ frameworkvalidator.String = (*inListValidator)(nil)

// NewInListValidator constructs a validator that allows only the provided values.
func NewInListValidator(values []string, ignoreCase bool) frameworkvalidator.String {
	lookup := make(map[string]struct{}, len(values))
	normalized := make([]string, 0, len(values))

	for _, candidate := range values {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}

		key := trimmed
		if ignoreCase {
			key = strings.ToLower(trimmed)
		}

		if _, exists := lookup[key]; exists {
			continue
		}

		lookup[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return &inListValidator{
		allowed:    normalized,
		ignoreCase: ignoreCase,
		lookup:     lookup,
	}
}

func (v *inListValidator) Description(_ context.Context) string {
	return "value must match one of the allowed strings"
}

func (v *inListValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *inListValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if len(v.allowed) == 0 {
		return
	}

	value := req.ConfigValue.ValueString()
	key := value
	if v.ignoreCase {
		key = strings.ToLower(value)
	}

	if _, ok := v.lookup[key]; ok {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Value Not Allowed",
		fmt.Sprintf("Value %q must be one of: %s", value, strings.Join(v.allowed, ", ")),
	)
}
