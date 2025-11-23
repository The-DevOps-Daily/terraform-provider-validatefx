package validators

import (
	"fmt"
)

// NonEmptyListValidator validates that a list is not empty.
type NonEmptyListValidator struct{}

// NewNonEmptyList creates a new NonEmptyListValidator.
func NewNonEmptyList() *NonEmptyListValidator {
	return &NonEmptyListValidator{}
}

// Validate returns nil when the provided list is not empty.
func (v *NonEmptyListValidator) Validate(values []string) error {
	if v == nil {
		return fmt.Errorf("validator not initialized")
	}

	if len(values) == 0 {
		return fmt.Errorf("list must not be empty")
	}

	return nil
}
