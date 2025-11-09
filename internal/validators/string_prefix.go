package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*stringPrefixValidator)(nil)

// StringPrefix returns a validator ensuring a string starts with one of the provided prefixes.
func StringPrefix(prefixes []string, ignoreCase bool) frameworkvalidator.String {
	display, normalized := normalizePrefixes(prefixes, ignoreCase)
	return &stringPrefixValidator{
		prefixes:   display,
		normalized: normalized,
		ignoreCase: ignoreCase,
	}
}

type stringPrefixValidator struct {
	prefixes   []string
	normalized []string
	ignoreCase bool
}

func (v *stringPrefixValidator) Description(_ context.Context) string {
	if len(v.prefixes) == 0 {
		return "string prefix validation"
	}

	return fmt.Sprintf("string must start with one of: %s", strings.Join(v.prefixes, ", "))
}

func (v *stringPrefixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *stringPrefixValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
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

	for _, prefix := range v.normalized {
		if strings.HasPrefix(candidate, prefix) {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Prefix",
		fmt.Sprintf("Value %q must start with one of: %s", value, strings.Join(v.prefixes, ", ")),
	)
}

func normalizePrefixes(values []string, lower bool) ([]string, []string) {
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
