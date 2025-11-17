package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzPublicIPValidator(f *testing.F) {
	v := PublicIP()
	for _, s := range []string{"8.8.8.8", "1.1.1.1", "10.0.0.1", "fd00::1", "2001:4860:4860::8888", "not-an-ip"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
	})
}
