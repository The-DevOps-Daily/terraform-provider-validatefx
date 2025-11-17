package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPublicIPFunction exposes the public IP validator as a Terraform function.
func NewPublicIPFunction() function.Function {
	return newStringValidationFunction(
		"public_ip",
		"Validate that an IP address is public (not private).",
		"Returns true when the input IP address is not in private IPv4 ranges or IPv6 ULA.",
		validators.PublicIP(),
	)
}
