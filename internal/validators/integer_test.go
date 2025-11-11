package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIntegerValidator_Valid(t *testing.T) {
	t.Parallel()
	v := Integer()
	cases := []string{"0", "42", "-7", "+10"}
	for _, c := range cases {
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			Path:        path.Root("value"),
			ConfigValue: types.StringValue(c),
		}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for %q, got: %v", c, resp.Diagnostics)
		}
	}
}

func TestIntegerValidator_Invalid(t *testing.T) {
	t.Parallel()
	v := Integer()
	cases := []string{"", "  ", "3.14", "1e3", "--1", "++2", "a42"}
	for _, c := range cases {
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			Path:        path.Root("value"),
			ConfigValue: types.StringValue(c),
		}, resp)
		if c == "" || c == "  " {
			// empty should be ignored (treated like optional)
			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for empty-like input %q", c)
			}
			continue
		}
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for %q", c)
		}
	}
}

func TestIntegerValidator_NullUnknown(t *testing.T) {
	t.Parallel()
	v := Integer()
	for name, val := range map[string]types.String{
		"null":    types.StringNull(),
		"unknown": types.StringUnknown(),
	} {
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			Path:        path.Root("value"),
			ConfigValue: val,
		}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("%s: unexpected diagnostics: %v", name, resp.Diagnostics)
		}
	}
}
