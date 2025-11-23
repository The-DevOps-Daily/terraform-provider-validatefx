package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzMIMEType(f *testing.F) {
	// Seed with valid MIME types
	f.Add("application/json")
	f.Add("text/html")
	f.Add("text/plain")
	f.Add("image/png")
	f.Add("image/svg+xml")
	f.Add("application/xml")
	f.Add("video/mp4")
	f.Add("text/html; charset=utf-8")
	f.Add("application/vnd.api+json")

	// Seed with invalid MIME types
	f.Add("")
	f.Add("notamimetype")
	f.Add("application/")
	f.Add("/json")
	f.Add("application / json")

	f.Fuzz(func(t *testing.T, input string) {
		v := MIMEType()
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic on any input
		v.ValidateString(context.Background(), req, resp)
	})
}
