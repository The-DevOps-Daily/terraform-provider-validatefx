package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringContainsValidator(t *testing.T) {
	t.Parallel()

	t.Run("valid substring", func(t *testing.T) {
		t.Parallel()

		validator := StringContains([]string{"foo", "bar"}, false)

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("prefix-foo-suffix"),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid substring", func(t *testing.T) {
		t.Parallel()

		validator := StringContains([]string{"foo", "bar"}, false)

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("baz"),
			Path:        path.Root("value"),
		}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for missing substring")
		}
	})

	t.Run("ignore case", func(t *testing.T) {
		t.Parallel()

		validator := StringContains([]string{"Foo"}, true)

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("prefix-foo-suffix"),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected case-insensitive match: %v", resp.Diagnostics)
		}
	})

	t.Run("unknown and null", func(t *testing.T) {
		t.Parallel()

		validator := StringContains([]string{"foo"}, false)

		for name, value := range map[string]types.String{
			"unknown": types.StringUnknown(),
			"null":    types.StringNull(),
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				resp := &frameworkvalidator.StringResponse{}
				validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
					ConfigValue: value,
					Path:        path.Root("value"),
				}, resp)

				if resp.Diagnostics.HasError() {
					t.Fatalf("expected no diagnostics for %s value", name)
				}
			})
		}
	})

	t.Run("empty substrings", func(t *testing.T) {
		t.Parallel()

		validator := StringContains(nil, false)

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("value"),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics when substrings list is empty: %v", resp.Diagnostics)
		}
	})
}

func TestPrepareSubstringCandidates(t *testing.T) {
	t.Parallel()

	display, normalized := normalizeStringList([]string{" Foo ", "bar", "foo"}, true)

	if len(display) != 2 {
		t.Fatalf("expected 2 display values, got %d", len(display))
	}

	if display[0] != "Foo" || display[1] != "bar" {
		t.Fatalf("unexpected display values: %v", display)
	}

	if normalized[0] != "foo" || normalized[1] != "bar" {
		t.Fatalf("unexpected normalized values: %v", normalized)
	}
}
