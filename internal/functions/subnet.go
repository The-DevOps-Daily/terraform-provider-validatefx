package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewSubnetFunction exposes the subnet validator as a Terraform function.
func NewSubnetFunction() function.Function {
	return newStringValidationFunction(
		"subnet",
		"Validate that a string is a subnet address (IP equals network) in CIDR notation.",
		"Returns true when the string is a valid IPv4/IPv6 subnet address in CIDR notation where the IP equals the network address.",
		validators.Subnet(),
	)
}
