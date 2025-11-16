package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzPrivateIPValidator(f *testing.F) {
	v := PrivateIP()
	for _, s := range []string{"10.0.0.1", "172.16.0.1", "192.168.1.1", "fd00::1", "8.8.8.8", "2001:db8::1", "not-an-ip"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
	})
}
