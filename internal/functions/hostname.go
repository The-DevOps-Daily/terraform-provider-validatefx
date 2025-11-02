package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewHostnameFunction exposes the hostname validator as a Terraform function.
func NewHostnameFunction() function.Function {
	return newStringValidationFunction(
		"hostname",
		"Validate that a string is a hostname compliant with RFC 1123.",
		"Returns true when the input string is a valid hostname, including optional trailing dot and punycode labels.",
		validators.Hostname(),
	)
}
