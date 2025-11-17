package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzIPRangeSizeValidator exercises the validator across random inputs to catch panics.
func FuzzIPRangeSizeValidator(f *testing.F) {
	v := NewIPRangeSizeValidator(8, 28)
	seeds := []string{
		"10.0.0.0/8",
		"10.0.0.0/16",
		"192.168.1.0/24",
		"2001:db8::/32",
		"2001:db8::/48",
		"not-a-cidr",
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("cidr"),
			ConfigValue: types.StringValue(s),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// No assertions: ensure no panics and diagnostics are produced safely.
	})
}
