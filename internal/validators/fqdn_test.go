package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFQDNValidator(t *testing.T) {
	t.Parallel()
	v := FQDN()

	run := func(s types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: s}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	cases := []struct {
		name    string
		val     types.String
		wantErr bool
	}{
		{"valid simple", types.StringValue("example.com"), false},
		{"valid multi-label", types.StringValue("app.prod.example.com"), false},
		{"invalid no dot", types.StringValue("localhost"), true},
		{"invalid empty label", types.StringValue("example..com"), true},
		{"invalid chars", types.StringValue("exa_mple.com"), true},
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
