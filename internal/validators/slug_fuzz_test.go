//go:build gofuzz || go1.18

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzSlug(f *testing.F) {
	// Seed with valid slugs
	f.Add("hello-world")
	f.Add("my-app")
	f.Add("web-server-01")
	f.Add("api-v2")
	f.Add("123")

	// Seed with invalid slugs
	f.Add("Hello-World")
	f.Add("-hello")
	f.Add("hello-")
	f.Add("hello--world")
	f.Add("hello_world")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		v := Slug()
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic
		v.ValidateString(context.Background(), req, resp)
	})
}
