package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPortNumberFunction exposes the port number validator as a Terraform function.
func NewPortNumberFunction() function.Function {
	return newStringValidationFunction(
		"port_number",
		"Validate that a string is a valid TCP/UDP port number (1..65535).",
		"Returns true when the input string represents an integer between 1 and 65535.",
		validators.PortNumber(),
	)
}
