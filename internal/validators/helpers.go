package validators

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// normalizeStringList takes a list of strings and returns deduplicated, trimmed values.
// If ignoreCase is true, it normalizes to lowercase for deduplication but preserves original case.
// Returns both display values (original case) and normalized values (lowercase if ignoreCase).
func normalizeStringList(values []string, ignoreCase bool) (display []string, normalized []string) {
	display = make([]string, 0, len(values))
	normalized = make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, raw := range values {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		key := trimmed
		if ignoreCase {
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

// parseFloat64 attempts to parse a string as a float64.
// Returns the parsed value and a diagnostic error if parsing fails.
func parseFloat64(value string, attrPath path.Path) (float64, diag.Diagnostic) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, diag.NewAttributeErrorDiagnostic(
			attrPath,
			"Invalid Number",
			"Value must be a valid number.",
		)
	}
	return num, nil
}

// parseIntInRange attempts to parse a string as an integer within a specified range [min, max].
// Returns the parsed value and a diagnostic error if parsing fails or value is out of range.
func parseIntInRange(value string, min, max int, attrPath path.Path, fieldName string) (int, diag.Diagnostic) {
	num, err := strconv.Atoi(value)
	if err != nil || num < min || num > max {
		return 0, diag.NewAttributeErrorDiagnostic(
			attrPath,
			fmt.Sprintf("Invalid %s", fieldName),
			fmt.Sprintf("Value %q must be an integer between %d and %d.", value, min, max),
		)
	}
	return num, nil
}

// validateStringInMap checks if a string value exists in a predefined map of valid values.
// Returns a diagnostic error if the value is not found in the map.
func validateStringInMap(value string, validValues map[string]bool, attrPath path.Path, errorTitle, fieldType string) diag.Diagnostic {
	if value == "" {
		return nil
	}

	if !validValues[value] {
		return diag.NewAttributeErrorDiagnostic(
			attrPath,
			errorTitle,
			fmt.Sprintf("Value %q is not a valid %s.", value, fieldType),
		)
	}
	return nil
}
