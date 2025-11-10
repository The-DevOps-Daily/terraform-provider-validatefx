package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestHexValidator(t *testing.T) {
	t.Parallel()

	v := Hex()

	run := func(val types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: val}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	cases := []struct {
		name    string
		val     types.String
		wantErr bool
	}{
		{"valid lowercase", types.StringValue("deadbeef"), false},
		{"valid uppercase", types.StringValue("DEADBEEF"), false},
		{"valid mixed", types.StringValue("CafeBABE"), false},
		{"invalid char", types.StringValue("xyz123"), true},
		{"empty string", types.StringValue(""), true},
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
