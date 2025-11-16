package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPrivateIPFunction exposes the private IP validator as a Terraform function.
func NewPrivateIPFunction() function.Function {
	return newStringValidationFunction(
		"private_ip",
		"Validate that an IP address is private (RFC1918 / IPv6 ULA).",
		"Returns true when the input IP address is a private IPv4 (RFC1918) or IPv6 ULA (fc00::/7).",
		validators.PrivateIP(),
	)
}
