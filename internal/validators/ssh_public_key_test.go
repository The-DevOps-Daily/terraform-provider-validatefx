package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSSHPublicKeyValidator(t *testing.T) {
	t.Parallel()

	v := SSHPublicKeyValidator()
	run := func(s types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: s}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	valid := types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKJf0N0nH7kz5Zr4xkz0GWWJrPq9uO2m6sR3j0s8v2QG test@example")

	cases := []struct {
		name    string
		val     types.String
		wantErr bool
	}{
		{"valid ed25519", valid, false},
		{"invalid", types.StringValue("not-a-key"), true},
		{"empty", types.StringValue(""), true},
		{"null", types.StringNull(), false},
		{"unknown", types.StringUnknown(), false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := run(tc.val)
			if tc.wantErr && !resp.Diagnostics.HasError() {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}
