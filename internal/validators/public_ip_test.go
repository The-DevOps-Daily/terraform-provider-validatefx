package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPublicIPValidator(t *testing.T) {
	t.Parallel()
	v := PublicIP()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	for _, s := range []string{"8.8.8.8", "1.1.1.1", "2001:4860:4860::8888"} {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected public IP %q to pass: %v", s, resp.Diagnostics)
		}
	}

	for _, s := range []string{"10.0.0.1", "192.168.1.1", "fd00::1"} {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected private IP %q to be rejected", s)
		}
	}
}
