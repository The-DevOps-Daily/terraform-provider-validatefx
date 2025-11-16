package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzURIValidator(f *testing.F) {
	v := URI()
	for _, s := range []string{
		"http://example.com", "https://example.com/path?x=1#frag", "urn:isbn:0451450523",
		"example.com", "http:///path", "://host",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
	})
}
