package validators

import (
	"context"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzCIDRValidator cross-checks validity against net.ParseCIDR and ensures no panics.
func FuzzCIDRValidator(f *testing.F) {
	for _, s := range []string{"", "10.0.0.0/8", "192.168.1.0/24", "2001:db8::/32", "1.2.3.4/33", "not/a/cidr"} {
		f.Add(s)
	}
	v := CIDR()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("cidr"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		_, _, err := net.ParseCIDR(s)
		expect := err == nil
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: parse-ok=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
