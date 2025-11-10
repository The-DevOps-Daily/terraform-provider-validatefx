package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewSSHPublicKeyFunction exposes the ssh_public_key validator as a Terraform function.
func NewSSHPublicKeyFunction() function.Function {
	return newStringValidationFunction(
		"ssh_public_key",
		"Validate that a string is a valid SSH public key.",
		"Returns true when the input string is a valid SSH public key in OpenSSH authorized_keys format.",
		validators.SSHPublicKeyValidator(),
	)
}
