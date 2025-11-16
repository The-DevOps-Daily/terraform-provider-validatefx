package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestInListValidator(t *testing.T) {
	t.Parallel()

	validator := NewInListValidator([]string{"alpha", "beta", "gamma"}, false)

	run := func(value types.String) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("value"),
			ConfigValue: value,
		}
		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), req, resp)
		return resp
	}

	t.Run("allowed value", func(t *testing.T) {
		t.Parallel()

		if resp := run(types.StringValue("beta")); resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics, got %v", resp.Diagnostics)
		}
	})

	t.Run("disallowed value", func(t *testing.T) {
		t.Parallel()

		if resp := run(types.StringValue("delta")); !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for disallowed value")
		}
	})

	t.Run("null and unknown", func(t *testing.T) {
		t.Parallel()

		for name, value := range map[string]types.String{
			"null":    types.StringNull(),
			"unknown": types.StringUnknown(),
		} {
			if resp := run(value); resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for %s", name)
			}
		}
	})
}

func TestInListValidatorIgnoreCase(t *testing.T) {
	t.Parallel()

	validator := NewInListValidator([]string{"ONE", "Two"}, true)

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("two"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected case-insensitive match, got %v", resp.Diagnostics)
	}
}

func TestInListValidatorCustomMessage(t *testing.T) {
	t.Parallel()

	validator := NewInListValidatorWithMessage([]string{"alpha", "beta"}, false, "Custom error message")

	req := frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: types.StringValue("delta"),
	}
	resp := &frameworkvalidator.StringResponse{}
	validator.ValidateString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics for custom message case")
	}
}
