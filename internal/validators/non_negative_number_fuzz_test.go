//go:build gofuzz || go1.18

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzNonNegativeNumber(f *testing.F) {
	// Seed with valid non-negative numbers
	f.Add("0")
	f.Add("1")
	f.Add("42")
	f.Add("3.14")
	f.Add("0.001")
	f.Add("+100")

	// Seed with invalid values
	f.Add("-1")
	f.Add("-3.14")
	f.Add("abc")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		v := NonNegativeNumber()
		req := validator.StringRequest{
			Path:        path.Root("test"),
			ConfigValue: types.StringValue(input),
		}
		resp := &validator.StringResponse{}

		// Should not panic
		v.ValidateString(context.Background(), req, resp)
	})
}
