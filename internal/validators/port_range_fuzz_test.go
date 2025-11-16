package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzPortRangeValidator(f *testing.F) {
	v := PortRange()

	seeds := []string{
		"80-8080",
		"0-65535",
		"443-443",
		"  22 -  2222 ",
		"-1-10",
		"10-70000",
		"bad",
		"80:8080",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		// No assertion: ensures no panic
	})
}
