package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPositiveNumberFunction exposes the positive_number validator as a Terraform function.
func NewPositiveNumberFunction() function.Function {
	return newStringValidationFunction(
		"positive_number",
		"Validate that a string represents a positive number.",
		"Returns true when the input string is a valid positive number (greater than zero).",
		validators.PositiveNumber(),
	)
}
