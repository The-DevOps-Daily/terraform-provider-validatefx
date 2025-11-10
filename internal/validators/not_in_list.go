package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type notInListValidator struct {
	disallowed []string
	ignoreCase bool
	lookup     map[string]struct{}
}

var _ frameworkvalidator.String = (*notInListValidator)(nil)

// NewNotInListValidator constructs a validator that fails when the value matches any disallowed value.
func NewNotInListValidator(values []string, ignoreCase bool) frameworkvalidator.String {
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

	return &notInListValidator{
		disallowed: normalized,
		ignoreCase: ignoreCase,
		lookup:     lookup,
	}
}

func (v *notInListValidator) Description(_ context.Context) string {
	return "value must not match any disallowed strings"
}

func (v *notInListValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *notInListValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if len(v.disallowed) == 0 {
		return
	}

	value := req.ConfigValue.ValueString()
	key := value
	if v.ignoreCase {
		key = strings.ToLower(value)
	}

	if _, ok := v.lookup[key]; ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value Disallowed",
			fmt.Sprintf("Value %q must not be one of: %s", value, strings.Join(v.disallowed, ", ")),
		)
		return
	}
}
