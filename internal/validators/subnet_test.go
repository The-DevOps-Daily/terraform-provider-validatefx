package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSubnetValidator(t *testing.T) {
	t.Parallel()
	v := Subnet()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	valids := []string{"192.168.1.0/24", "10.0.0.0/8", "2001:db8::/64"}
	for _, s := range valids {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected valid subnet %q: %v", s, resp.Diagnostics)
		}
	}

	invalids := []string{"192.168.1.1/24", "10.0.0.1/8", "2001:db8::1/64", "not-a-cidr"}
	for _, s := range invalids {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected invalid subnet %q", s)
		}
	}
}
