package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewIntegerFunction exposes the integer validator as a Terraform function.
func NewIntegerFunction() function.Function {
	return newStringValidationFunction(
		"integer",
		"Validate that a string represents a valid integer.",
		"Returns true when the input string is a valid integer (optional leading + or -).",
		validators.Integer(),
	)
}
