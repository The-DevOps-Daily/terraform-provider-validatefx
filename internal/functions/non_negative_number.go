package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewNonNegativeNumberFunction exposes the non_negative_number validator as a Terraform function.
func NewNonNegativeNumberFunction() function.Function {
	return newStringValidationFunction(
		"non_negative_number",
		"Validate that a string represents a non-negative number.",
		"Returns true when the input string is a valid non-negative number (zero or greater).",
		validators.NonNegativeNumber(),
	)
}
