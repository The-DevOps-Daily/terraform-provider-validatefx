package validators

import "strings"

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
