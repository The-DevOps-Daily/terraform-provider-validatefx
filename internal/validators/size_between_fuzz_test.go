//go:build gofuzz || go1.18

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzSizeBetween(f *testing.F) {
	// Seed with valid values
	f.Add("5")
	f.Add("1")
	f.Add("10")
	f.Add("5.5")
	f.Add("0")

	// Seed with invalid values
	f.Add("0")
	f.Add("11")
	f.Add("-1")
	f.Add("abc")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		v := SizeBetween("1", "10")
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic
		v.ValidateString(context.Background(), req, resp)
	})
}
