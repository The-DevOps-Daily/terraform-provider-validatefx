package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringSuffixValidator(t *testing.T) {
	t.Parallel()

	validator := StringSuffix(".log", ".txt")

	t.Run("valid suffix", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("service.log"),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid suffix", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringValue("service.log.bak"),
			Path:        path.Root("value"),
		}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for invalid suffix")
		}
	})

	t.Run("null value", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringNull(),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for null value")
		}
	})

	t.Run("unknown value", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.StringResponse{}
		validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
			ConfigValue: types.StringUnknown(),
			Path:        path.Root("value"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for unknown value")
		}
	})
}

func TestNormalizeSuffixes(t *testing.T) {
	t.Parallel()

	result := normalizeSuffixes([]string{" .log", "", "\t.txt\n"})

	if len(result) != 2 {
		t.Fatalf("expected 2 suffixes, got %d", len(result))
	}

	if result[0] != ".log" || result[1] != ".txt" {
		t.Fatalf("unexpected normalized suffixes: %v", result)
	}
}
