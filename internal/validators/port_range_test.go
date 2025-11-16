package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPortRangeValidator(t *testing.T) {
	t.Parallel()

	v := PortRange()

	run := func(val types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: val}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	cases := []struct {
		name    string
		value   types.String
		wantErr bool
	}{
		{"simple", types.StringValue("80-8080"), false},
		{"spaces", types.StringValue(" 0 - 65535 "), false},
		{"equal bounds", types.StringValue("443-443"), false},
		{"out of order", types.StringValue("1000-10"), true},
		{"too high", types.StringValue("1-70000"), true},
		{"negative", types.StringValue("-1-10"), true},
		{"bad format", types.StringValue("80:8080"), true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := run(tc.value)
			if got := resp.Diagnostics.HasError(); got != tc.wantErr {
				t.Fatalf("wantErr=%t got=%t diags=%v", tc.wantErr, got, resp.Diagnostics)
			}
		})
	}

	// null/unknown pass-through
	for name, v := range map[string]types.String{"null": types.StringNull(), "unknown": types.StringUnknown()} {
		t.Run(name, func(t *testing.T) {
			resp := run(v)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}
