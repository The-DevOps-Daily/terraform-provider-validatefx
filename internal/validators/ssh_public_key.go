package validators

import (
	"context"
	"fmt"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"golang.org/x/crypto/ssh"
)

type sshPublicKeyValidator struct{}

var _ frameworkvalidator.String = (*sshPublicKeyValidator)(nil)

// SSHPublicKeyValidator returns a validator that verifies an SSH public key in OpenSSH authorized_keys format.
func SSHPublicKeyValidator() frameworkvalidator.String { return sshPublicKeyValidator{} }

func (sshPublicKeyValidator) Description(_ context.Context) string {
	return "value must be a valid SSH public key in authorized_keys format"
}

func (sshPublicKeyValidator) MarkdownDescription(ctx context.Context) string {
	return sshPublicKeyValidator{}.Description(ctx)
}

func (sshPublicKeyValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	if s == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid SSH Public Key",
			"Value must not be empty and should be in authorized_keys format (e.g., 'ssh-ed25519 AAAA... comment').",
		)
		return
	}

	if _, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(s)); err != nil || len(rest) > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid SSH Public Key",
			fmt.Sprintf("Value %q is not a valid SSH public key: %v", s, err),
		)
	}
}
