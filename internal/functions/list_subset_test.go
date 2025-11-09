package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestListSubsetFunction(t *testing.T) {
	t.Parallel()

	fn := NewListSubsetFunction()
	ctx := context.Background()

	run := func(values attr.Value, allowed attr.Value) *function.RunResponse {
		resp := &function.RunResponse{}
		args := function.NewArgumentsData([]attr.Value{values, allowed})
		fn.Run(ctx, function.RunRequest{Arguments: args}, resp)
		return resp
	}

	t.Run("valid subset", func(t *testing.T) {
		t.Parallel()

		values := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
			basetypes.NewStringValue("admin"),
			basetypes.NewStringValue("viewer"),
		})
		allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
			basetypes.NewStringValue("admin"),
			basetypes.NewStringValue("viewer"),
			basetypes.NewStringValue("editor"),
		})

		resp := run(values, allowed)
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		result, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !result.ValueBool() {
			t.Fatalf("expected true result")
		}
	})

	t.Run("invalid subset", func(t *testing.T) {
		t.Parallel()

		values := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
			basetypes.NewStringValue("admin"),
			basetypes.NewStringValue("operator"),
		})
		allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
			basetypes.NewStringValue("admin"),
			basetypes.NewStringValue("viewer"),
		})

		resp := run(values, allowed)
		if resp.Error == nil {
			t.Fatalf("expected error for invalid subset")
		}
	})

	t.Run("unknown values", func(t *testing.T) {
		t.Parallel()

		resp := run(basetypes.NewListUnknown(basetypes.StringType{}), basetypes.NewListValueMust(basetypes.StringType{}, nil))
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		result, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !result.IsUnknown() {
			t.Fatalf("expected unknown result when values unknown")
		}
	})

	t.Run("empty allowed list", func(t *testing.T) {
		t.Parallel()

		resp := run(basetypes.NewListValueMust(basetypes.StringType{}, nil), basetypes.NewListValueMust(basetypes.StringType{}, nil))
		if resp.Error == nil {
			t.Fatalf("expected error when allowed list empty")
		}
	})
}
