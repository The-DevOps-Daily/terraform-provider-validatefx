package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewUsernameFunction exposes the username validator as a Terraform function.
func NewUsernameFunction() function.Function {
	return newStringValidationFunction(
		"username",
		"Validate that a string is a valid username.",
		"Returns true when the input string is a valid username (3-20 characters, letters, digits, or underscores).",
		validators.DefaultUsernameValidator(),
	)
}
