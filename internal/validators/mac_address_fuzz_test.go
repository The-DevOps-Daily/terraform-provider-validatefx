package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzMACAddressValidator(f *testing.F) {
	seeds := []string{"", "00:1A:2B:3C:4D:5E", "00-1A-2B-3C-4D-5E", "001A2B3C4D5E", "GG:HH:II:JJ:KK:LL"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := MACAddress()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("mac"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// Rely on validator logic; assert no panics.
	})
}
