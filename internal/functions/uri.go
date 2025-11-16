package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewURIFunction exposes the URI validator as a Terraform function.
func NewURIFunction() function.Function {
	return newStringValidationFunction(
		"uri",
		"Validate that a string is a URI.",
		"Returns true when the input string is a valid URI supporting common schemes (http, https, ftp, ssh, postgres, etc.).",
		validators.URI(),
	)
}
