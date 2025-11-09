package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestListSubset_List(t *testing.T) {
	t.Parallel()

	validator := ListSubset([]string{"admin", "editor", "viewer"})

	t.Run("valid subset", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.ListResponse{}
		validator.ValidateList(context.Background(), frameworkvalidator.ListRequest{
			ConfigValue: mustList(t, "admin", "viewer"),
			Path:        path.Root("roles"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid subset", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.ListResponse{}
		validator.ValidateList(context.Background(), frameworkvalidator.ListRequest{
			ConfigValue: mustList(t, "admin", "operator"),
			Path:        path.Root("roles"),
		}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for invalid element")
		}
	})

	t.Run("null list", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.ListResponse{}
		validator.ValidateList(context.Background(), frameworkvalidator.ListRequest{
			ConfigValue: types.ListNull(types.StringType),
			Path:        path.Root("roles"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for null list, got %v", resp.Diagnostics)
		}
	})

	t.Run("unknown list", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.ListResponse{}
		validator.ValidateList(context.Background(), frameworkvalidator.ListRequest{
			ConfigValue: types.ListUnknown(types.StringType),
			Path:        path.Root("roles"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for unknown list, got %v", resp.Diagnostics)
		}
	})
}

func TestListSubset_Set(t *testing.T) {
	t.Parallel()

	validator := SetSubset([]string{"feature-x", "feature-y"})

	t.Run("valid subset", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.SetResponse{}
		validator.ValidateSet(context.Background(), frameworkvalidator.SetRequest{
			ConfigValue: mustSet(t, "feature-x", "feature-y"),
			Path:        path.Root("features"),
		}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid subset", func(t *testing.T) {
		t.Parallel()

		resp := &frameworkvalidator.SetResponse{}
		validator.ValidateSet(context.Background(), frameworkvalidator.SetRequest{
			ConfigValue: mustSet(t, "feature-x", "feature-z"),
			Path:        path.Root("features"),
		}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for invalid element")
		}
	})
}

func mustList(t *testing.T, values ...string) types.List {
	t.Helper()

	list, diags := types.ListValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build list value: %v", diags)
	}

	return list
}

func mustSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build set value: %v", diags)
	}

	return set
}
