package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewUUIDv4OnlyFunction exposes the UUIDv4Only validator as a Terraform function.
func NewUUIDv4OnlyFunction() function.Function {
	return newStringValidationFunction(
		"uuidv4_only",
		"Validate that a string is a UUID version 4.",
		"Returns true when the input is specifically a UUID version 4; false otherwise.",
		validators.UUIDv4Only(),
	)
}
