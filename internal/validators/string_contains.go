package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*stringContainsValidator)(nil)

// StringContains returns a validator that ensures the value contains one of the provided substrings.
func StringContains(substrings []string, ignoreCase bool) frameworkvalidator.String {
	display, normalized := prepareSubstringCandidates(substrings, ignoreCase)

	return &stringContainsValidator{
		substrings: display,
		normalized: normalized,
		ignoreCase: ignoreCase,
	}
}

type stringContainsValidator struct {
	substrings []string
	normalized []string
	ignoreCase bool
}

func (v *stringContainsValidator) Description(_ context.Context) string {
	if len(v.substrings) == 0 {
		return "string substring validation"
	}

	return fmt.Sprintf("string must contain one of: %s", strings.Join(v.substrings, ", "))
}

func (v *stringContainsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *stringContainsValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if len(v.normalized) == 0 {
		return
	}

	value := req.ConfigValue.ValueString()
	candidate := value
	if v.ignoreCase {
		candidate = strings.ToLower(candidate)
	}

	for _, substring := range v.normalized {
		if strings.Contains(candidate, substring) {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Substring Not Found",
		fmt.Sprintf("Value %q must contain one of: %s", value, strings.Join(v.substrings, ", ")),
	)
}

func prepareSubstringCandidates(values []string, lower bool) ([]string, []string) {
	display := make([]string, 0, len(values))
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, raw := range values {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		key := trimmed
		if lower {
			key = strings.ToLower(trimmed)
		}

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}

		display = append(display, trimmed)
		normalized = append(normalized, key)
	}

	return display, normalized
}
