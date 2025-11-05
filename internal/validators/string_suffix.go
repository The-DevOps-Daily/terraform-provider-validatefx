package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure interface compliance.
var _ frameworkvalidator.String = (*stringSuffixValidator)(nil)

// StringSuffix returns a validator that ensures a value ends with one of the provided suffixes.
func StringSuffix(suffixes ...string) frameworkvalidator.String {
	return &stringSuffixValidator{suffixes: normalizeSuffixes(suffixes)}
}

type stringSuffixValidator struct {
	suffixes []string
}

func (v *stringSuffixValidator) Description(_ context.Context) string {
	return v.describe()
}

func (v *stringSuffixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *stringSuffixValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if len(v.suffixes) == 0 {
		return
	}

	if hasAllowedSuffix(value, v.suffixes) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Suffix",
		fmt.Sprintf("Value %q must end with one of: %s.", value, strings.Join(v.suffixes, ", ")),
	)
}

func (v *stringSuffixValidator) describe() string {
	if len(v.suffixes) == 0 {
		return "string suffix validation"
	}

	return fmt.Sprintf("string must end with one of: %s", strings.Join(v.suffixes, ", "))
}

func normalizeSuffixes(values []string) []string {
	normalized := make([]string, 0, len(values))

	for _, suffix := range values {
		trimmed := strings.TrimSpace(suffix)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	return normalized
}

func hasAllowedSuffix(value string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(value, suffix) {
			return true
		}
	}

	return false
}
