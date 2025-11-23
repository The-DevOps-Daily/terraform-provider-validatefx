package validators

import (
	"fmt"
	"sort"
	"strings"
)

// MapKeysMatchValidator validates that map keys match a set of allowed or required keys.
type MapKeysMatchValidator struct {
	allowedKeys  map[string]struct{}
	requiredKeys map[string]struct{}
}

// NewMapKeysMatch creates a validator that checks map keys against allowed and required sets.
// If allowedKeys is empty, all keys are allowed.
// If requiredKeys is provided, those keys must be present.
func NewMapKeysMatch(allowedKeys, requiredKeys []string) *MapKeysMatchValidator {
	allowed := make(map[string]struct{}, len(allowedKeys))
	for _, k := range allowedKeys {
		allowed[k] = struct{}{}
	}

	required := make(map[string]struct{}, len(requiredKeys))
	for _, k := range requiredKeys {
		required[k] = struct{}{}
	}

	return &MapKeysMatchValidator{
		allowedKeys:  allowed,
		requiredKeys: required,
	}
}

// Validate checks if the provided map keys satisfy the allowed/required constraints.
func (v *MapKeysMatchValidator) Validate(keys []string) error {
	if v == nil {
		return fmt.Errorf("validator not initialized")
	}

	// Check for required keys
	if len(v.requiredKeys) > 0 {
		present := make(map[string]struct{}, len(keys))
		for _, k := range keys {
			present[k] = struct{}{}
		}

		missing := []string{}
		for req := range v.requiredKeys {
			if _, found := present[req]; !found {
				missing = append(missing, req)
			}
		}

		if len(missing) > 0 {
			sort.Strings(missing)
			return fmt.Errorf("missing required keys: %s", strings.Join(missing, ", "))
		}
	}

	// Check for allowed keys (only if allowedKeys is not empty)
	if len(v.allowedKeys) > 0 {
		disallowed := []string{}
		for _, k := range keys {
			if _, allowed := v.allowedKeys[k]; !allowed {
				disallowed = append(disallowed, k)
			}
		}

		if len(disallowed) > 0 {
			sort.Strings(disallowed)
			return fmt.Errorf("disallowed keys: %s", strings.Join(disallowed, ", "))
		}
	}

	return nil
}
