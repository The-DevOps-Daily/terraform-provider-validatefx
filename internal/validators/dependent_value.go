package validators

import (
	"context"
	"fmt"
)

// DependentValueValidator validates that if a condition value is non-empty,
// then a dependent value must also be non-empty.
type DependentValueValidator struct{}

// NewDependentValue creates a new validator that checks dependent value relationships.
func NewDependentValue() *DependentValueValidator {
	return &DependentValueValidator{}
}

// Description returns a plain text description of the validator.
func (DependentValueValidator) Description(_ context.Context) string {
	return "if condition value is set, dependent value must also be set"
}

// MarkdownDescription returns a markdown formatted description of the validator.
func (v DependentValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation on the given values.
// Returns an error if conditionValue is non-empty but dependentValue is empty.
func (DependentValueValidator) Validate(conditionValue, dependentValue string) error {
	// If condition is set (non-empty), dependent must also be set
	if conditionValue != "" && dependentValue == "" {
		return fmt.Errorf("when condition value is set, dependent value must also be provided")
	}

	return nil
}
