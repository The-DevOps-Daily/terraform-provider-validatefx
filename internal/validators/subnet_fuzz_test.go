package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzSubnetValidator(f *testing.F) {
	v := Subnet()
	for _, s := range []string{"192.168.1.0/24", "192.168.1.1/24", "10.0.0.0/8", "not-a-cidr"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
	})
}
