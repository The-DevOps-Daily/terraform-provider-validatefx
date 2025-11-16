package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPortRangeFunction exposes the port range validator as a Terraform function.
func NewPortRangeFunction() function.Function {
	return newStringValidationFunction(
		"port_range",
		"Validate that a string is a valid port range (start-end).",
		"Returns true when the input string matches a valid port range in the form 'start-end' with ports within 0..65535 and start <= end.",
		validators.PortRange(),
	)
}
