package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPasswordStrength(t *testing.T) {
	t.Parallel()
	v := PasswordStrengthValidator()

	run := func(s types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("password"), ConfigValue: s}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	cases := []struct {
		name    string
		val     types.String
		wantErr bool
	}{
		{name: "valid", val: types.StringValue("Abc@1234"), wantErr: false},
		{name: "short", val: types.StringValue("abc"), wantErr: true},
		{name: "missing-kinds", val: types.StringValue("abcdefghi"), wantErr: true},
		{name: "null", val: types.StringNull(), wantErr: false},
		{name: "unknown", val: types.StringUnknown(), wantErr: false},
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
