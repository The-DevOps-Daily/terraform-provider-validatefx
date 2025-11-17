package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPortNumberValidator(t *testing.T) {
	t.Parallel()
	v := PortNumber()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	for _, s := range []string{"1", "80", "443", "65535"} {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected valid %q: %v", s, resp.Diagnostics)
		}
	}
	for _, s := range []string{"0", "65536", "-1", "abc", "80.0", ""} {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected invalid %q", s)
		}
	}
}
