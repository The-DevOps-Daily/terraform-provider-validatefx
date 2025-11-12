package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestJWTValidator(t *testing.T) {
	t.Parallel()
	v := JWT()

	run := func(s types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: s}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	valid := types.StringValue("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

	cases := []struct {
		name    string
		val     types.String
		wantErr bool
	}{
		{"valid", valid, false},
		{"missing parts", types.StringValue("abc.def"), true},
		{"empty segment", types.StringValue("abc..def"), true},
		{"bad base64", types.StringValue("abc.def.!@#"), true},
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
