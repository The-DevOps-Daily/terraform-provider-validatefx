package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewMACAddressFunction exposes the MAC address validator as a Terraform function.
func NewMACAddressFunction() function.Function {
	return newStringValidationFunction(
		"mac_address",
		"Validate that a string is a MAC address in colon, dash, or compact format.",
		"Returns true when the input represents a MAC address in colon-separated, dash-separated, or compact hexadecimal form.",
		validators.MACAddress(),
	)
}
