package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPrivateIPValidator(t *testing.T) {
	t.Parallel()

	v := PrivateIP()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	ok := []string{"10.0.0.1", "172.16.0.1", "172.31.255.254", "192.168.1.1", "fd00::1", "fc00:abcd::1"}
	bad := []string{"8.8.8.8", "172.15.0.1", "172.32.0.1", "1.2.3.4", "2001:db8::1", "::1", "not-an-ip"}

	for _, s := range ok {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected %s to be private, got %v", s, resp.Diagnostics)
		}
	}
	for _, s := range bad {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected %s to be rejected", s)
		}
	}

	// Null/unknown pass-through
	req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringNull()}
	resp := &frameworkvalidator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics for null")
	}

	req = frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringUnknown()}
	resp = &frameworkvalidator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics for unknown")
	}
}
