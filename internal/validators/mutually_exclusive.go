package validators

import (
	"fmt"
)

// MutuallyExclusiveValidator validates that only one of the provided values is set (non-null/non-empty).
type MutuallyExclusiveValidator struct{}

// NewMutuallyExclusive creates a new MutuallyExclusiveValidator.
func NewMutuallyExclusive() *MutuallyExclusiveValidator {
	return &MutuallyExclusiveValidator{}
}

// Validate returns nil when exactly one of the provided values is set.
// A value is considered "set" if it is non-empty.
func (v *MutuallyExclusiveValidator) Validate(values []string) error {
	if v == nil {
		return fmt.Errorf("validator not initialized")
	}

	setCount := 0
	for _, val := range values {
		if val != "" {
			setCount++
		}
	}

	if setCount == 0 {
		return fmt.Errorf("at least one value must be set")
	}

	if setCount > 1 {
		return fmt.Errorf("only one value must be set, but %d values are set", setCount)
	}

	return nil
}
