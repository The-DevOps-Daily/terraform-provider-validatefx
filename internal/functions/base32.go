package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewBase32Function exposes the base32 validator as a Terraform function.
func NewBase32Function() function.Function {
	return newStringValidationFunction(
		"base32",
		"Validate that a string is Base32 encoded.",
		"Returns true when the input string can be decoded from Base32.",
		validators.Base32Validator(),
	)
}
