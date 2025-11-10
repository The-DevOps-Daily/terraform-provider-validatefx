package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBase32Validator(t *testing.T) {
	t.Parallel()

	v := Base32Validator()
	run := func(s types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: s}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	cases := []struct {
		name    string
		s       types.String
		wantErr bool
	}{
		{"valid 1", types.StringValue("NBSWY3DP"), false},
		{"valid 2", types.StringValue("OZQWY2LEMF2GK==="), false},
		{"invalid", types.StringValue("not-base32"), true},
		{"empty", types.StringValue(""), false},
		{"null", types.StringNull(), false},
		{"unknown", types.StringUnknown(), false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := run(tc.s)
			if tc.wantErr && !resp.Diagnostics.HasError() {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}
