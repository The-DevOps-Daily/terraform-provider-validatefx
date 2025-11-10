package validators

import (
	"fmt"
	"sort"
)

// SetEqualsValidator compares two lists of strings for set equality (order-independent, unique elements).
// The validator holds the expected list and validates an input list against it.
type SetEqualsValidator struct {
	expectedSet map[string]struct{}
	expectedOrd []string
}

// NewSetEquals creates a new SetEqualsValidator from the expected values.
func NewSetEquals(expected []string) *SetEqualsValidator {
	set := make(map[string]struct{}, len(expected))
	ord := make([]string, 0, len(expected))
	for _, v := range expected {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		ord = append(ord, v)
	}
	sort.Strings(ord)
	return &SetEqualsValidator{expectedSet: set, expectedOrd: ord}
}

// Validate returns nil when the provided list has the same unique elements as the expected list.
// Duplicates in either slice are ignored for comparison purposes.
func (v *SetEqualsValidator) Validate(values []string) error {
	if v == nil {
		return fmt.Errorf("validator not initialized")
	}

	left := make(map[string]struct{}, len(values))
	ord := make([]string, 0, len(values))
	for _, val := range values {
		if _, seen := left[val]; !seen {
			left[val] = struct{}{}
			ord = append(ord, val)
		}
	}
	sort.Strings(ord)

	if len(left) != len(v.expectedSet) {
		return fmt.Errorf("set mismatch: %v != %v", ord, v.expectedOrd)
	}

	for key := range left {
		if _, ok := v.expectedSet[key]; !ok {
			return fmt.Errorf("set mismatch: %v != %v", ord, v.expectedOrd)
		}
	}

	for key := range v.expectedSet {
		if _, ok := left[key]; !ok {
			return fmt.Errorf("set mismatch: %v != %v", ord, v.expectedOrd)
		}
	}

	return nil
}
