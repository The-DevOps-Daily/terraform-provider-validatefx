package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

func FuzzSSHPublicKeyValidator(f *testing.F) {
	seeds := []string{
		"",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKJf0N0nH7kz5Zr4xkz0GWWJrPq9uO2m6sR3j0s8v2QG user@example",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7 bogus",
		"not a key",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	v := SSHPublicKeyValidator()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("sshkey"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			// empty must error
			if !resp.Diagnostics.HasError() {
				t.Fatalf("empty must error for SSH key")
			}
			return
		}

		// Oracle via ssh.ParseAuthorizedKey
		_, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(s))
		expect := err == nil && len(rest) == 0
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for SSH key; expect=%v diagErr=%v", expect, resp.Diagnostics.HasError())
		}
	})
}
